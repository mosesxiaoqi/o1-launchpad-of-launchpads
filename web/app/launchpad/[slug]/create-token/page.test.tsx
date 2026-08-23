import { fireEvent, render, screen } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import { TokenLaunchForm } from '@/components/token-launch-form'

const mocks = vi.hoisted(() => ({
  sendTransactionAsync: vi.fn(),
  waitForTransactionReceipt: vi.fn(),
  prepareLaunch: vi.fn(),
  transaction: vi.fn(),
  parseReceipt: vi.fn(),
}))

vi.mock('wagmi', () => ({
  useAccount: () => ({
    address: '0x1111111111111111111111111111111111111111',
    chainId: 84532,
    isConnected: true,
  }),
  useSendTransaction: () => ({ sendTransactionAsync: mocks.sendTransactionAsync }),
  usePublicClient: () => ({ waitForTransactionReceipt: mocks.waitForTransactionReceipt }),
}))

vi.mock('@/lib/wagmi', () => ({ chain: { id: 84532 } }))
vi.mock('@/lib/api', () => ({
  launchpadApi: { prepareLaunch: mocks.prepareLaunch, transaction: mocks.transaction },
}))
vi.mock('@/lib/receipt', () => ({ parseLaunchedReceipt: mocks.parseReceipt }))

const launchpad = {
  id: `0x${'ab'.repeat(32)}` as `0x${string}`,
  chain_id: 84532,
  slug: 'my-pad',
  name: 'My Pad',
  description: '',
  logo_url: '',
  primary_color: '#5B5CF6',
  owner: '0x1111111111111111111111111111111111111111' as `0x${string}`,
  treasury: '0x2222222222222222222222222222222222222222' as `0x${string}`,
  active: true,
  created_at: '2026-08-24T00:00:00Z',
}

const config = {
  chain_id: 84532,
  registry: '0x3333333333333333333333333333333333333333' as `0x${string}`,
  factory: '0x4444444444444444444444444444444444444444' as `0x${string}`,
  hook: '0x5555555555555555555555555555555555555555' as `0x${string}`,
  fee_escrow: '0x6666666666666666666666666666666666666666' as `0x${string}`,
  quote: '0x7777777777777777777777777777777777777777' as `0x${string}`,
  config_version: 1,
}

function fillToken() {
  fireEvent.change(screen.getByLabelText('Token 名称'), { target: { value: 'Demo Token' } })
  fireEvent.change(screen.getByLabelText('Symbol'), { target: { value: 'DEMO' } })
}

describe('token launch form', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mocks.prepareLaunch.mockResolvedValue({
      chain_id: 84532,
      from: launchpad.owner,
      to: config.factory,
      data: `0x${'12'.repeat(64)}`,
      value: '0',
      deadline: 1_800_000_000,
      review: 'fixed supply, permanent liquidity, 1% protocol fee, 0.5% LaaS fee, 16s anti-snipe',
    })
    mocks.sendTransactionAsync.mockResolvedValue(`0x${'34'.repeat(32)}`)
    mocks.waitForTransactionReceipt.mockResolvedValue({ status: 'success', logs: [] })
    mocks.parseReceipt.mockReturnValue({
      token: '0x8888888888888888888888888888888888888888',
      poolId: `0x${'99'.repeat(32)}`,
      launchpadId: launchpad.id,
      creator: launchpad.owner,
      quote: config.quote,
      supply: 1_000_000_000_000_000_000_000_000_000n,
      tickSpacing: 60,
    })
    mocks.transaction.mockResolvedValue({ status: 'confirmed' })
  })

  it('shows every immutable launch term before wallet signing', async () => {
    render(<TokenLaunchForm launchpad={launchpad} config={config} />)
    fillToken()
    fireEvent.click(screen.getByRole('button', { name: '生成发行计划' }))
    expect(await screen.findByText('1,000,000,000')).toBeInTheDocument()
    for (const text of ['永久流动性', '1% 协议费', '0.5% LaaS', '1.5% 总费率', '16 秒 anti-snipe', config.quote, launchpad.id, launchpad.treasury]) {
      expect(screen.getByText(text)).toBeInTheDocument()
    }
    expect(screen.getByText(/fixed supply/)).toBeInTheDocument()
    expect(mocks.sendTransactionAsync).not.toHaveBeenCalled()
  })

  it('broadcasts the prepared transaction unchanged and shows receipt before indexing', async () => {
    let resolveTransaction!: (value: { status: string }) => void
    mocks.transaction.mockReturnValue(new Promise((done) => (resolveTransaction = done)))
    render(<TokenLaunchForm launchpad={launchpad} config={config} />)
    fillToken()
    fireEvent.click(screen.getByRole('button', { name: '生成发行计划' }))
    fireEvent.click(await screen.findByRole('button', { name: '确认并签名' }))
    expect(await screen.findByText(/0x8888888888888888888888888888888888888888/)).toBeInTheDocument()
    expect(screen.getByText('等待索引确认')).toBeInTheDocument()
    expect(mocks.sendTransactionAsync).toHaveBeenCalledWith({
      to: config.factory,
      data: `0x${'12'.repeat(64)}`,
      value: 0n,
      chainId: 84532,
    })
    resolveTransaction({ status: 'confirmed' })
    expect(await screen.findByText('发行已确认')).toBeInTheDocument()
  })

  it('requires a fresh review when the plan becomes stale', async () => {
    mocks.sendTransactionAsync.mockRejectedValue(new Error('StaleConfig'))
    render(<TokenLaunchForm launchpad={launchpad} config={config} />)
    fillToken()
    fireEvent.click(screen.getByRole('button', { name: '生成发行计划' }))
    fireEvent.click(await screen.findByRole('button', { name: '确认并签名' }))
    expect(await screen.findByText('计划已过期，请重新生成并审阅')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: '生成发行计划' })).toBeInTheDocument()
  })
})
