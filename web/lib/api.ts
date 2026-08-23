export type Problem = {
  type: string
  title: string
  status: number
  detail: string
  code: string
}

export class ApiError extends Error {
  constructor(
    readonly status: number,
    readonly code: string,
    message: string,
  ) {
    super(message)
    this.name = 'ApiError'
  }
}

const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8888'

export async function api<T>(path: string, init: RequestInit = {}): Promise<T> {
  const headers = new Headers(init.headers)
  if (init.body && !headers.has('content-type')) headers.set('content-type', 'application/json')
  const response = await fetch(`${apiUrl}${path}`, {
    ...init,
    headers,
    credentials: 'include',
  })
  if (!response.ok) {
    const problem = (await response.json()) as Partial<Problem>
    throw new ApiError(
      problem.status || response.status,
      problem.code || 'request_failed',
      problem.detail || problem.title || response.statusText,
    )
  }
  if (response.status === 204) return undefined as T
  return response.json() as Promise<T>
}

export type Config = {
  chain_id: number
  registry: `0x${string}`
  factory: `0x${string}`
  hook: `0x${string}`
  fee_escrow: `0x${string}`
  quote: `0x${string}`
  config_version: number
}

export type Launchpad = {
  id: `0x${string}`
  chain_id: number
  slug: string
  name: string
  description: string
  logo_url: string
  primary_color: string
  owner: `0x${string}`
  treasury: `0x${string}`
  active: boolean
  created_at: string
}

export type PreparedTransaction = {
  chain_id: number
  from: `0x${string}`
  to: `0x${string}`
  data: `0x${string}`
  value: string
  deadline: number
  review: string
}

export type Token = {
  chain_id: number
  token: `0x${string}`
  pool_id: `0x${string}`
  launchpad_id: `0x${string}`
  launchpad_slug: string
  creator: `0x${string}`
  quote: `0x${string}`
  supply: string
  tx_hash: `0x${string}`
  block_number: number
  status: string
  created_at: string
}

export type ChainTransaction = {
  chain_id: number
  tx_hash: `0x${string}`
  kind: string
  status: 'pending' | 'confirming' | 'confirmed' | 'reverted'
  block_number: number
  failure_reason: string
  updated_at: string
}

function post<T>(path: string, body: unknown) {
  return api<T>(path, { method: 'POST', body: JSON.stringify(body) })
}

export const launchpadApi = {
  health: () => api<{ status: string }>('/health/ready'),
  config: () => api<Config>('/v1/config'),
  createChallenge: (body: { chain_id: number; address: string; domain: string; uri: string }) =>
    post<{ challenge_id: string; message: string; expires_at: string }>('/v1/auth/challenge', body),
  verifySignature: (body: { challenge_id: string; address: string; signature: string }) =>
    post<{ address: string; expires_at: string }>('/v1/auth/verify', body),
  listLaunchpads: (cursor = '') =>
    api<{ launchpads: Launchpad[]; next_cursor: string }>(
      `/v1/launchpads?limit=25&cursor=${encodeURIComponent(cursor)}`,
    ),
  getLaunchpad: (slug: string) => api<Launchpad>(`/v1/launchpads/${encodeURIComponent(slug)}`),
  listLaunchpadTokens: (slug: string, cursor = '') =>
    api<{ tokens: Token[]; next_cursor: string }>(
      `/v1/launchpads/${encodeURIComponent(slug)}/tokens?limit=25&cursor=${encodeURIComponent(cursor)}`,
    ),
  createLaunchpad: (body: {
    chain_id: number
    slug: string
    name: string
    description?: string
    logo_url?: string
    primary_color?: string
    registry_tx_hash: string
  }) => post<Launchpad>('/v1/launchpads', body),
  prepareLaunch: (body: {
    chain_id: number
    launchpad_id: string
    name: string
    symbol: string
    contract_uri?: string
    salt: string
    wallet: string
  }) =>
    post<PreparedTransaction>('/v1/launches/prepare', body),
  transaction: (chainId: number, hash: string) =>
    api<ChainTransaction>(`/v1/transactions/${chainId}/${encodeURIComponent(hash)}`),
  quoteSwap: (body: Record<string, unknown>) =>
    post<{ quote_id: string; amount_in: string; amount_out: string; fee: string; expires_at: number }>(
      '/v1/swaps/quote',
      body,
    ),
  prepareSwap: (body: Record<string, unknown>) =>
    post<PreparedTransaction>('/v1/swaps/prepare', body),
  fees: (recipient: string, currency: string) =>
    api<{ amount: string; block_number: number }>(
      `/v1/fees/${encodeURIComponent(recipient)}?currency=${encodeURIComponent(currency)}`,
    ),
  prepareFeeClaim: (body: Record<string, unknown>) =>
    post<PreparedTransaction>('/v1/claims/fees/prepare', body),
}
