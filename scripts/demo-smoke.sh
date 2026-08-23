#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
DEPLOYMENT_FILE=${DEPLOYMENT_FILE:-"$ROOT_DIR/deployments/base-sepolia.json"}
API_URL=${SMOKE_API_URL:-}
AUTH_ORIGIN=${SMOKE_AUTH_ORIGIN:-}
RPC_URL=${BASE_SEPOLIA_RPC_URL:-}
KEYSTORE=${SMOKE_KEYSTORE:-}
PASSWORD_FILE=${SMOKE_KEYSTORE_PASSWORD_FILE:-}
ZERO_ADDRESS=0x0000000000000000000000000000000000000000
COOKIE_JAR=$(mktemp)
trap 'rm -f "$COOKIE_JAR"' EXIT

fail() { printf 'smoke: %s\n' "$*" >&2; exit 1; }
need() { command -v "$1" >/dev/null || fail "missing command: $1"; }
lower() { printf '%s' "$1" | tr '[:upper:]' '[:lower:]'; }
same_address() { [[ "$(lower "$1")" == "$(lower "$2")" ]]; }
is_address() { [[ "$1" =~ ^0x[0-9a-fA-F]{40}$ ]]; }
pad_address() { printf '%064s' "${1#0x}" | tr ' ' 0 | tr '[:upper:]' '[:lower:]'; }
require_contract() {
  is_address "$2" || fail "$1 is not a valid address"
  [[ "$(lower "$2")" != "$ZERO_ADDRESS" ]] || fail "$1 is the zero address"
  [[ "$(cast code --rpc-url "$RPC_URL" "$2")" != 0x ]] || fail "$1 has no deployed code"
}
validate_plan() {
  plan=$1
  expected_to=$2
  expected_value=$3
  deadline_required=$4
  expected_selector=$5
  [[ "$(jq -r '.chain_id' <<<"$plan")" == 84532 ]] || fail 'prepared transaction has the wrong chain_id'
  same_address "$(jq -r '.from' <<<"$plan")" "$wallet" || fail 'prepared transaction has the wrong from address'
  same_address "$(jq -r '.to' <<<"$plan")" "$expected_to" || fail 'prepared transaction has an unexpected target'
  [[ "$(jq -r '.value' <<<"$plan")" == "$expected_value" ]] || fail 'prepared transaction has an unexpected value'
  data=$(jq -r '.data' <<<"$plan")
  [[ "$data" =~ ^0x[0-9a-fA-F]+$ ]] || fail 'prepared transaction calldata is invalid'
  [[ "$(lower "${data:0:10}")" == "$(lower "$expected_selector")" ]] || fail 'prepared transaction function selector is unexpected'
  if [[ "$deadline_required" == yes ]]; then
    [[ "$(jq -r '.deadline' <<<"$plan")" -gt "$(date +%s)" ]] || fail 'prepared transaction deadline is expired'
  fi
}
send_plan() {
  plan=$1
  to=$(jq -r '.to' <<<"$plan")
  data=$(jq -r '.data' <<<"$plan")
  value=$(jq -r '.value' <<<"$plan")
  cast call --rpc-url "$RPC_URL" --from "$wallet" --value "$value" "$to" --data "$data" >/dev/null || fail 'prepared transaction simulation failed'
  cast send --json --rpc-url "$RPC_URL" "${wallet_args[@]}" --confirmations 2 --value "$value" "$to" "$data" | jq -r '.transactionHash'
}
need curl
need jq
need cast
http() { curl --connect-timeout 10 --max-time 60 "$@"; }
[[ -n "$RPC_URL" ]] || fail 'BASE_SEPOLIA_RPC_URL is required'
[[ "$API_URL" =~ ^https://[^/]+$ ]] || fail 'SMOKE_API_URL must be an HTTPS origin without a trailing slash'
[[ "$AUTH_ORIGIN" =~ ^https://([A-Za-z0-9.-]+)(:[0-9]+)?$ ]] || fail 'SMOKE_AUTH_ORIGIN must be an HTTPS frontend origin without a path'
auth_domain=${BASH_REMATCH[1]}
[[ -f "$KEYSTORE" ]] || fail 'SMOKE_KEYSTORE must point to a Foundry keystore file'
[[ -f "$PASSWORD_FILE" ]] || fail 'SMOKE_KEYSTORE_PASSWORD_FILE must point to its local password file'
wallet_args=(--keystore "$KEYSTORE" --password-file "$PASSWORD_FILE")
[[ -f "$DEPLOYMENT_FILE" ]] || fail "deployment file not found: $DEPLOYMENT_FILE"

[[ "$(jq -r '.chainId' "$DEPLOYMENT_FILE")" == 84532 ]] || fail 'deployment chainId is not 84532'
for field in poolManager universalRouter registry factory hook feeEscrow quote laasTreasury deploymentBlock; do
  value=$(jq -r ".$field // empty" "$DEPLOYMENT_FILE")
  [[ -n "$value" && "$value" != null ]] || fail "deployment field $field is not populated"
done
[[ "$(jq -r '.deploymentBlock' "$DEPLOYMENT_FILE")" =~ ^[1-9][0-9]*$ ]] || fail 'deploymentBlock must be a positive integer'

chain_id=$(cast chain-id --rpc-url "$RPC_URL")
[[ "$chain_id" == 84532 ]] || fail "RPC chain id is $chain_id, expected 84532"
http -fsS "$API_URL/health/ready" >/dev/null || fail "API is not ready at $API_URL"

api_config=$(http -fsS "$API_URL/v1/config")
[[ "$(jq -r '.chain_id' <<<"$api_config")" == 84532 ]] || fail 'API config chain_id is not 84532'
for pair in 'registry:registry' 'factory:factory' 'hook:hook' 'fee_escrow:feeEscrow' 'quote:quote' 'laas_treasury:laasTreasury'; do
  api_field=${pair%%:*}
  file_field=${pair##*:}
  api_value=$(jq -r ".$api_field" <<<"$api_config")
  file_value=$(jq -r ".$file_field" "$DEPLOYMENT_FILE")
  [[ "$(lower "$api_value")" == "$(lower "$file_value")" ]] || fail "$api_field differs between API and deployment file"
done

quote=$(jq -r '.quote' <<<"$api_config")
[[ "$(lower "$quote")" == "$ZERO_ADDRESS" ]] || fail 'MVP automated Swap requires native ETH Quote; ERC-20 Permit2 approval is not implemented yet'

pool_manager=$(jq -r '.poolManager' "$DEPLOYMENT_FILE")
universal_router=$(jq -r '.universalRouter' "$DEPLOYMENT_FILE")
registry=$(jq -r '.registry' <<<"$api_config")
factory=$(jq -r '.factory' <<<"$api_config")
hook=$(jq -r '.hook' <<<"$api_config")
escrow=$(jq -r '.fee_escrow' <<<"$api_config")
laas=$(jq -r '.laas_treasury' <<<"$api_config")
for contract in "poolManager:$pool_manager" "universalRouter:$universal_router" "registry:$registry" "factory:$factory" "hook:$hook" "feeEscrow:$escrow"; do
  require_contract "${contract%%:*}" "${contract##*:}"
done
same_address "$(cast call "$hook" 'factory()(address)' --rpc-url "$RPC_URL")" "$factory" || fail 'Hook factory wiring mismatch'
same_address "$(cast call "$hook" 'feeEscrow()(address)' --rpc-url "$RPC_URL")" "$escrow" || fail 'Hook escrow wiring mismatch'
same_address "$(cast call "$escrow" 'hook()(address)' --rpc-url "$RPC_URL")" "$hook" || fail 'Escrow hook wiring mismatch'
same_address "$(cast call "$factory" 'registry()(address)' --rpc-url "$RPC_URL")" "$registry" || fail 'Factory registry wiring mismatch'
same_address "$(cast call "$factory" 'poolManager()(address)' --rpc-url "$RPC_URL")" "$pool_manager" || fail 'Factory PoolManager wiring mismatch'
same_address "$(cast call "$factory" 'hook()(address)' --rpc-url "$RPC_URL")" "$hook" || fail 'Factory Hook wiring mismatch'
same_address "$(cast call "$factory" 'quote()(address)' --rpc-url "$RPC_URL")" "$quote" || fail 'Factory Quote wiring mismatch'
same_address "$(cast call "$factory" 'laasTreasury()(address)' --rpc-url "$RPC_URL")" "$laas" || fail 'Factory LaaS wiring mismatch'
same_address "$(cast call "$hook" 'poolManager()(address)' --rpc-url "$RPC_URL")" "$pool_manager" || fail 'Hook PoolManager wiring mismatch'
same_address "$(cast call "$escrow" 'poolManager()(address)' --rpc-url "$RPC_URL")" "$pool_manager" || fail 'Escrow PoolManager wiring mismatch'
suffix=${hook: -4}
(( (16#$suffix & 0x3fff) == 0x2acc )) || fail 'Hook permission bits are invalid'

wallet=$(cast wallet address "${wallet_args[@]}")
treasury=${SMOKE_LAUNCHPAD_TREASURY:-}
referrer=${SMOKE_REFERRER_ADDRESS:-}
[[ "$treasury" =~ ^0x[0-9a-fA-F]{40}$ ]] || fail 'SMOKE_LAUNCHPAD_TREASURY must be a distinct valid address'
[[ "$referrer" =~ ^0x[0-9a-fA-F]{40}$ ]] || fail 'SMOKE_REFERRER_ADDRESS must be a distinct valid address'
[[ "$(lower "$treasury")" != "$ZERO_ADDRESS" && "$(lower "$referrer")" != "$ZERO_ADDRESS" && "$(lower "$laas")" != "$ZERO_ADDRESS" ]] || fail 'treasury, referrer, and LaaS addresses must be non-zero'
wallet_lower=$(lower "$wallet")
treasury_lower=$(lower "$treasury")
referrer_lower=$(lower "$referrer")
laas_lower=$(lower "$laas")
[[ "$wallet_lower" != "$treasury_lower" && "$wallet_lower" != "$referrer_lower" && "$wallet_lower" != "$laas_lower" ]] || fail 'creator, protocol treasury, referrer, and LaaS addresses must be distinct'
[[ "$treasury_lower" != "$referrer_lower" && "$treasury_lower" != "$laas_lower" && "$referrer_lower" != "$laas_lower" ]] || fail 'creator, protocol treasury, referrer, and LaaS addresses must be distinct'

slug=${SMOKE_LAUNCHPAD_SLUG:-"demo-$(date +%s)"}
token_name=${SMOKE_TOKEN_NAME:-O1-Demo-Token}
token_symbol=${SMOKE_TOKEN_SYMBOL:-O1DEMO}
token_uri=${SMOKE_TOKEN_URI:-ipfs://o1-demo}
swap_amount=${SMOKE_SWAP_AMOUNT_WEI:-1000000000000000}
[[ "$swap_amount" =~ ^[1-9][0-9]*$ ]] || fail 'SMOKE_SWAP_AMOUNT_WEI must be a positive integer'
salt=${SMOKE_LAUNCH_SALT:-$(cast keccak "$slug-$wallet")}

challenge_request=$(jq -nc --arg address "$wallet" --arg domain "$auth_domain" --arg uri "$AUTH_ORIGIN" \
  '{chain_id:84532,address:$address,domain:$domain,uri:$uri}')
challenge=$(http -fsS -X POST -H 'content-type: application/json' -d "$challenge_request" "$API_URL/v1/auth/challenge")
challenge_id=$(jq -r '.challenge_id' <<<"$challenge")
message=$(jq -r '.message' <<<"$challenge")
[[ "$message" == "$auth_domain wants you to sign in to o1 Launchpad."* && "$message" == *"URI: $AUTH_ORIGIN"* ]] || fail 'challenge message does not match SMOKE_AUTH_ORIGIN / FRONTEND_ORIGIN'
signature=$(cast wallet sign "${wallet_args[@]}" "$message")
verify_request=$(jq -nc --arg challenge_id "$challenge_id" --arg address "$wallet" --arg signature "$signature" \
  '{challenge_id:$challenge_id,address:$address,signature:$signature}')
http -fsS -c "$COOKIE_JAR" -X POST -H 'content-type: application/json' -d "$verify_request" "$API_URL/v1/auth/verify" >/dev/null

cast call "$registry" 'createLaunchpad(string,address)' "$slug" "$treasury" --from "$wallet" --rpc-url "$RPC_URL" >/dev/null || fail 'Registry create simulation failed'
registry_receipt=$(cast send --json --rpc-url "$RPC_URL" "${wallet_args[@]}" --confirmations 2 \
  "$registry" 'createLaunchpad(string,address)' "$slug" "$treasury")
registry_tx=$(jq -r '.transactionHash' <<<"$registry_receipt")
[[ "$registry_tx" =~ ^0x[0-9a-fA-F]{64}$ ]] || fail 'Registry transaction did not return a hash'

brand_request=$(jq -nc --arg slug "$slug" --arg name "$slug" --arg tx "$registry_tx" \
  '{chain_id:84532,slug:$slug,name:$name,description:"Base Sepolia smoke",logo_url:"",primary_color:"#5B5CF6",registry_tx_hash:$tx}')
launchpad=$(http -fsS -b "$COOKIE_JAR" -X POST -H 'content-type: application/json' -d "$brand_request" "$API_URL/v1/launchpads")
launchpad_id=$(jq -r '.id' <<<"$launchpad")

launch_request=$(jq -nc --arg id "$launchpad_id" --arg name "$token_name" --arg symbol "$token_symbol" \
  --arg uri "$token_uri" --arg salt "$salt" --arg wallet "$wallet" \
  '{chain_id:84532,launchpad_id:$id,name:$name,symbol:$symbol,contract_uri:$uri,salt:$salt,wallet:$wallet}')
launch_plan=$(http -fsS -X POST -H 'content-type: application/json' -d "$launch_request" "$API_URL/v1/launches/prepare")
validate_plan "$launch_plan" "$factory" 0 yes 0x22d8980f
launch_data=$(lower "$(jq -r '.data' <<<"$launch_plan")")
[[ "$launch_data" == *"$(lower "${launchpad_id#0x}")"* && "$launch_data" == *"$(lower "${salt#0x}")"* ]] || fail 'Launch calldata does not bind the requested launchpadId and salt'
launch_tx=$(send_plan "$launch_plan")

token_json=''
for _ in $(seq 1 60); do
  token_page=$(http -fsS "$API_URL/v1/launchpads/$slug/tokens?limit=25" || true)
  token_json=$(jq -c --arg tx "$(lower "$launch_tx")" '.tokens[]? | select((.tx_hash|ascii_downcase)==$tx)' <<<"$token_page" | head -1)
  [[ -n "$token_json" ]] && break
  sleep 2
done
[[ -n "$token_json" ]] || fail 'Indexer did not expose the launched token within 120 seconds'
token=$(jq -r '.token' <<<"$token_json")
pool_id=$(jq -r '.pool_id' <<<"$token_json")
indexed_launchpad_id=$(jq -r '.launchpad_id' <<<"$token_json")
[[ "$(lower "$indexed_launchpad_id")" == "$(lower "$launchpad_id")" ]] || fail 'indexed token belongs to the wrong launchpadId'

sleep 17
quote_request=$(jq -nc --arg token "$token" --arg amount "$swap_amount" --arg wallet "$wallet" --arg referrer "$referrer" \
  '{chain_id:84532,token:$token,amount:$amount,buy:true,wallet:$wallet,referrer:$referrer}')
swap_quote=$(http -fsS -X POST -H 'content-type: application/json' -d "$quote_request" "$API_URL/v1/swaps/quote")
prepare_swap_request=$(jq -nc --arg id "$(jq -r '.quote_id' <<<"$swap_quote")" --arg wallet "$wallet" \
  '{chain_id:84532,quote_id:$id,slippage_bps:100,wallet:$wallet}')
swap_plan=$(http -fsS -X POST -H 'content-type: application/json' -d "$prepare_swap_request" "$API_URL/v1/swaps/prepare")
validate_plan "$swap_plan" "$universal_router" "$swap_amount" yes 0x3593564c
swap_data=$(lower "$(jq -r '.data' <<<"$swap_plan")")
[[ "$swap_data" == *"$(pad_address "$token")"* && "$swap_data" == *"$(pad_address "$hook")"* && "$swap_data" == *"$(pad_address "$referrer")"* ]] || fail 'Swap calldata does not bind the expected Token, Hook, and Referrer'
swap_decoded=$(cast decode-calldata 'execute(bytes,bytes[],uint256)' "$swap_data")
[[ "$(sed -n '1p' <<<"$swap_decoded")" == 0x10 ]] || fail 'Universal Router command is not V4_SWAP'
[[ "$(sed -n '3p' <<<"$swap_decoded" | awk '{print $1}')" == "$(jq -r '.deadline' <<<"$swap_plan")" ]] || fail 'Router calldata deadline differs from prepared metadata'
swap_tx=$(send_plan "$swap_plan")

fee_creator=$(http -fsS "$API_URL/v1/fees/$wallet?currency=$quote")
fee_protocol=$(http -fsS "$API_URL/v1/fees/$treasury?currency=$quote")
fee_referrer=$(http -fsS "$API_URL/v1/fees/$referrer?currency=$quote")
fee_laas=$(http -fsS "$API_URL/v1/fees/$laas?currency=$quote")
for fee in "$fee_creator" "$fee_protocol" "$fee_referrer" "$fee_laas"; do
  [[ $(jq -r '.amount' <<<"$fee") =~ ^[1-9][0-9]*$ ]] || fail 'one of the four fee balances is zero after Swap'
done

claim_request=$(jq -nc --arg currency "$quote" --arg recipient "$wallet" --arg wallet "$wallet" \
  '{chain_id:84532,currency:$currency,recipient:$recipient,wallet:$wallet}')
claim_plan=$(http -fsS -X POST -H 'content-type: application/json' -d "$claim_request" "$API_URL/v1/claims/fees/prepare")
validate_plan "$claim_plan" "$escrow" 0 no 0x21c0b342
claim_decoded=$(cast decode-calldata 'claim(address,address)' "$(jq -r '.data' <<<"$claim_plan")")
same_address "$(sed -n '1p' <<<"$claim_decoded")" "$wallet" || fail 'Claim calldata recipient mismatch'
same_address "$(sed -n '2p' <<<"$claim_decoded")" "$quote" || fail 'Claim calldata currency mismatch'
claim_tx=$(send_plan "$claim_plan")
creator_after=$(http -fsS "$API_URL/v1/fees/$wallet?currency=$quote" | jq -r '.amount')
[[ "$creator_after" == 0 ]] || fail "creator balance is $creator_after after Claim, expected 0"

evidence_file=${SMOKE_EVIDENCE_FILE:-"$ROOT_DIR/artifacts/base-sepolia-smoke-$(date -u +%Y%m%dT%H%M%SZ).json"}
mkdir -p "$(dirname "$evidence_file")"
fee_total=$(jq -n --arg a "$(jq -r '.amount' <<<"$fee_creator")" --arg b "$(jq -r '.amount' <<<"$fee_protocol")" --arg c "$(jq -r '.amount' <<<"$fee_referrer")" --arg d "$(jq -r '.amount' <<<"$fee_laas")" '$a|tonumber + ($b|tonumber) + ($c|tonumber) + ($d|tonumber)')
jq -n --arg registryTx "$registry_tx" --arg launchTx "$launch_tx" --arg poolId "$pool_id" \
  --arg token "$token" --arg launchpadId "$launchpad_id" --arg swapTx "$swap_tx" --arg claimTx "$claim_tx" \
  --arg creatorFee "$(jq -r '.amount' <<<"$fee_creator")" --arg protocolFee "$(jq -r '.amount' <<<"$fee_protocol")" \
  --arg referrerFee "$(jq -r '.amount' <<<"$fee_referrer")" --arg laasFee "$(jq -r '.amount' <<<"$fee_laas")" --arg feeTotal "$fee_total" \
  --arg creatorBlock "$(jq -r '.block_number' <<<"$fee_creator")" --arg protocolBlock "$(jq -r '.block_number' <<<"$fee_protocol")" --arg referrerBlock "$(jq -r '.block_number' <<<"$fee_referrer")" --arg laasBlock "$(jq -r '.block_number' <<<"$fee_laas")" \
  '{chainId:84532,registryTx:$registryTx,launchTx:$launchTx,poolId:$poolId,token:$token,launchpadId:$launchpadId,swapTx:$swapTx,fees:{creator:{amount:$creatorFee,block:$creatorBlock},protocol:{amount:$protocolFee,block:$protocolBlock},referrer:{amount:$referrerFee,block:$referrerBlock},laas:{amount:$laasFee,block:$laasBlock},total:$feeTotal},claimTx:$claimTx,explorer:{registry:("https://sepolia.basescan.org/tx/"+$registryTx),launch:("https://sepolia.basescan.org/tx/"+$launchTx),swap:("https://sepolia.basescan.org/tx/"+$swapTx),claim:("https://sepolia.basescan.org/tx/"+$claimTx)}}' >"$evidence_file"
printf 'Base Sepolia smoke passed. Evidence: %s\n' "$evidence_file"
