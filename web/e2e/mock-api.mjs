import { createServer } from 'node:http'
import { encodeAbiParameters, encodeEventTopics, parseAbi } from 'viem'

const wallet = '0x1111111111111111111111111111111111111111'
const registry = '0x2222222222222222222222222222222222222222'
const factory = '0x3333333333333333333333333333333333333333'
const hook = '0x4444444444444444444444444444444444444444'
const escrow = '0x5555555555555555555555555555555555555555'
const quote = '0x6666666666666666666666666666666666666666'
const laasTreasury = '0x7777777777777777777777777777777777777777'
const tokenAddress = '0x8888888888888888888888888888888888888888'
const poolId = `0x${'99'.repeat(32)}`
const launchpadId = `0x${'aa'.repeat(32)}`
const launchTx = `0x${'cc'.repeat(32)}`
const blockHash = `0x${'dd'.repeat(32)}`
const launchedAbi = parseAbi([
  'event Launched(address indexed token, bytes32 indexed poolId, bytes32 indexed launchpadId, address creator, address quote, uint256 supply, int24 tickSpacing)',
])
const launchedTopics = encodeEventTopics({
  abi: launchedAbi,
  eventName: 'Launched',
  args: { token: tokenAddress, poolId, launchpadId },
})
const launchedData = encodeAbiParameters(
  [{ type: 'address' }, { type: 'address' }, { type: 'uint256' }, { type: 'int24' }],
  [wallet, quote, 1_000_000_000n, 60],
)

let launchpad
let indexed = false
let allowIndex = false
let session = false

function corsHeaders() {
  return {
    'access-control-allow-origin': 'http://127.0.0.1:3100',
    'access-control-allow-credentials': 'true',
    'access-control-allow-headers': 'content-type',
    'access-control-allow-methods': 'GET,POST,OPTIONS',
    'content-type': 'application/json',
  }
}

function send(response, status, body, headers = {}) {
  response.writeHead(status, { ...corsHeaders(), ...headers })
  response.end(JSON.stringify(body))
}

async function body(request) {
  const chunks = []
  for await (const chunk of request) chunks.push(chunk)
  return chunks.length ? JSON.parse(Buffer.concat(chunks).toString()) : {}
}

function config() {
  return { chain_id: 84532, registry, factory, hook, fee_escrow: escrow, quote, laas_treasury: laasTreasury, config_version: 1 }
}

function token() {
  return {
    chain_id: 84532, token: tokenAddress, pool_id: poolId, launchpad_id: launchpadId,
    launchpad_slug: 'my-pad', creator: wallet, quote, supply: '1000000000', tx_hash: launchTx,
    block_number: 100, status: 'confirmed', created_at: '2026-08-24T00:00:00Z',
  }
}

function receipt(hash) {
  const logs = hash === launchTx ? [{
    address: factory, topics: launchedTopics, data: launchedData, blockNumber: '0x64',
    transactionHash: launchTx, transactionIndex: '0x0', blockHash, logIndex: '0x0', removed: false,
  }] : []
  return {
    transactionHash: hash, transactionIndex: '0x0', blockHash, blockNumber: '0x64', from: wallet,
    to: hash === launchTx ? factory : registry, cumulativeGasUsed: '0x5208', gasUsed: '0x5208',
    contractAddress: null, logs, logsBloom: `0x${'00'.repeat(256)}`, status: '0x1',
    effectiveGasPrice: '0x1', type: '0x2',
  }
}

async function handleRpc(request, response) {
  const payload = await body(request)
  let result
  if (payload.method === 'eth_chainId') result = '0x14a34'
  else if (payload.method === 'eth_blockNumber') result = '0x66'
  else if (payload.method === 'eth_getTransactionReceipt') result = receipt(payload.params[0])
  else if (payload.method === 'eth_getTransactionByHash') result = null
  else return send(response, 200, { jsonrpc: '2.0', id: payload.id, error: { code: -32601, message: `unsupported ${payload.method}` } })
  send(response, 200, { jsonrpc: '2.0', id: payload.id, result })
}

createServer(async (request, response) => {
  const url = new URL(request.url, 'http://127.0.0.1:3999')
  if (request.method === 'OPTIONS') return send(response, 204, {})
  if (url.pathname === '/rpc') return handleRpc(request, response)
  if (url.pathname === '/__reset' && request.method === 'POST') {
    launchpad = undefined
    indexed = false
    allowIndex = false
    session = false
    return send(response, 200, { ok: true })
  }
  if (url.pathname === '/__confirm-index' && request.method === 'POST') {
    allowIndex = true
    return send(response, 200, { ok: true })
  }
  if (url.pathname === '/health/ready') return send(response, 200, { status: 'ok' })
  if (url.pathname === '/v1/config') return send(response, 200, config())
  if (url.pathname === '/v1/auth/challenge') {
    const input = await body(request)
    if (input.chain_id !== 84532 || input.address.toLowerCase() !== wallet ||
      input.domain !== '127.0.0.1' || input.uri !== 'http://127.0.0.1:3100') {
      return send(response, 400, { detail: 'invalid challenge request' })
    }
    return send(response, 200, { challenge_id: 'challenge-1', message: 'Sign local E2E challenge', expires_at: '2026-08-25T00:00:00Z' })
  }
  if (url.pathname === '/v1/auth/verify') {
    const input = await body(request)
    if (input.challenge_id !== 'challenge-1' || input.address.toLowerCase() !== wallet || !/^0x[0-9a-f]{130}$/.test(input.signature)) {
      return send(response, 401, { detail: 'invalid signature' })
    }
    session = true
    return send(response, 200, { address: wallet, expires_at: '2026-08-25T00:00:00Z' }, { 'set-cookie': 'o1_session=e2e; Path=/; HttpOnly; SameSite=Lax' })
  }
  if (url.pathname === '/v1/launchpads' && request.method === 'POST') {
    if (!session || !request.headers.cookie?.includes('o1_session=e2e')) return send(response, 401, { detail: 'missing session' })
    const input = await body(request)
    launchpad = {
      id: launchpadId, chain_id: 84532, slug: input.slug, name: input.name, description: input.description,
      logo_url: input.logo_url, primary_color: input.primary_color, owner: wallet, treasury: wallet,
      active: true, created_at: '2026-08-24T00:00:00Z',
    }
    return send(response, 200, launchpad)
  }
  if (url.pathname === '/v1/launchpads/my-pad' && launchpad) return send(response, 200, launchpad)
  if (url.pathname === '/v1/launchpads/other-pad') return send(response, 200, {
    id: `0x${'ee'.repeat(32)}`, chain_id: 84532, slug: 'other-pad', name: 'Other Pad', description: '',
    logo_url: '', primary_color: '#000000', owner: wallet, treasury: wallet, active: true,
    created_at: '2026-08-24T00:00:00Z',
  })
  if (url.pathname === '/v1/launchpads/my-pad/tokens') return send(response, 200, { tokens: indexed ? [token()] : [], next_cursor: '' })
  if (url.pathname === '/v1/launchpads/other-pad/tokens') return send(response, 200, { tokens: [], next_cursor: '' })
  if (url.pathname === '/v1/launches/prepare') return send(response, 200, {
    chain_id: 84532, from: wallet, to: factory, data: '0x1234', value: '0', deadline: 1900000000,
    review: 'Local E2E launch review',
  })
  if (url.pathname === `/v1/transactions/84532/${launchTx}`) {
    indexed = allowIndex
    return send(response, 200, {
      chain_id: 84532, tx_hash: launchTx, kind: 'launch', status: allowIndex ? 'confirmed' : 'confirming', block_number: allowIndex ? 100 : 0,
      failure_reason: '', updated_at: '2026-08-24T00:00:00Z',
    })
  }
  send(response, 404, { type: 'about:blank', title: 'Not found', status: 404, detail: url.pathname, code: 'not_found' })
}).listen(3999, '127.0.0.1')
