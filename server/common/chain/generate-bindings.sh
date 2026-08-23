#!/bin/sh
set -eu

TASK_CHAIN_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
TASK_REPO_DIR=$(CDPATH= cd -- "$TASK_CHAIN_DIR/../../.." && pwd)
TASK_ABIGEN_BIN=${ABIGEN_BIN:-"$(go env GOPATH)/bin/abigen"}

mkdir -p "$TASK_CHAIN_DIR/bindings"
jq -c '.abi' "$TASK_REPO_DIR/contracts/out/MultiTenantLaunchpadFactory.sol/MultiTenantLaunchpadFactory.json" |
  "$TASK_ABIGEN_BIN" --abi - --pkg bindings --type MultiTenantLaunchpadFactory --out "$TASK_CHAIN_DIR/bindings/factory.go"
jq -c '.abi' "$TASK_REPO_DIR/contracts/out/FeeEscrow.sol/FeeEscrow.json" |
  "$TASK_ABIGEN_BIN" --abi - --pkg bindings --type FeeEscrow --out "$TASK_CHAIN_DIR/bindings/feeescrow.go"
