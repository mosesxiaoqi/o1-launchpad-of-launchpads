import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import { FeeBalances } from '@/components/fee-balances'
import { SwapForm } from '@/components/swap-form'
import Page from './page'

const mocks = vi.hoisted(() => ({
  wallet: '0x1111111111111111111111111111111111111111' as `0x${string}`,
  account: { address: '0x1111111111111111111111111111111111111111' as `0x${string}`, chainId: 84532, isConnected: true },
  balance: 1_000n,
  fees: vi.fn(),
  prepareFeeClaim: vi.fn(),
  quoteSwap: vi.fn(),
  prepareSwap: vi.fn(),
  sendTransactionAsync: vi.fn(),
  waitForTransactionReceipt: vi.fn(),
  getToken: vi.fn(),
  getLaunchpad: vi.fn(),
  config: vi.fn(),
}))
const wallet = mocks.wallet

vi.mock('wagmi', () => ({
  useAccount: () => mocks.account,
  useBalance: () => ({ data: { value: mocks.balance } }),
  useReadContract: () => ({ data: mocks.balance }),
  useSendTransaction: () => ({ sendTransactionAsync: mocks.sendTransactionAsync }),
  usePublicClient: () => ({ waitForTransactionReceipt: mocks.waitForTransactionReceipt }),
}))
vi.mock('@/components/connect-wallet', () => ({ ConnectWallet: () => <button>连接钱包</button> }))
vi.mock('@/lib/wagmi', () => ({ chain: { id: 84532 } }))
vi.mock('@/lib/api', () => ({
  launchpadApi: {
    getToken: mocks.getToken,
    getLaunchpad: mocks.getLaunchpad,
    config: mocks.config,
    fees: mocks.fees,
    prepareFeeClaim: mocks.prepareFeeClaim,
    quoteSwap: mocks.quoteSwap,
    prepareSwap: mocks.prepareSwap,
  },
}))

const token = {
  chain_id: 84532,
  token: '0x2222222222222222222222222222222222222222' as `0x${string}`,
  pool_id: `0x${'33'.repeat(32)}` as `0x${string}`,
  launchpad_id: `0x${'44'.repeat(32)}` as `0x${string}`,
  launchpad_slug: 'my-pad',
  creator: wallet as `0x${string}`,
  quote: '0x5555555555555555555555555555555555555555' as `0x${string}`,
  supply: '1000000000',
  tx_hash: `0x${'66'.repeat(32)}` as `0x${string}`,
  block_number: 1,
  status: 'confirmed',
  created_at: '2026-08-24T00:00:00Z',
}

const launchpad = {
  id: token.launchpad_id,
  chain_id: 84532,
  slug: 'my-pad',
  name: 'My Pad',
  description: '',
  logo_url: '',
  primary_color: '#123456',
  owner: wallet as `0x${string}`,
  treasury: '0x7777777777777777777777777777777777777777' as `0x${string}`,
  active: true,
  created_at: token.created_at,
}

const config = {
  chain_id: 84532,
  registry: '0x8888888888888888888888888888888888888888' as `0x${string}`,
  factory: '0x9999999999999999999999999999999999999999' as `0x${string}`,
  hook: '0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa' as `0x${string}`,
  fee_escrow: '0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb' as `0x${string}`,
  quote: token.quote,
  laas_treasury: '0xcccccccccccccccccccccccccccccccccccccccc' as `0x${string}`,
  config_version: 1,
}

describe('token page', () => {
  it('shows the actual ERC-20 quote address', async () => {
    mocks.getToken.mockResolvedValue(token)
    mocks.getLaunchpad.mockResolvedValue(launchpad)
    mocks.config.mockResolvedValue(config)
    render(await Page({ params: Promise.resolve({ address: token.token }) } as never))
    expect(screen.getByText(token.quote)).toBeInTheDocument()
  })
})

describe('fees', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mocks.account.chainId = 84532
    mocks.balance = 1_000n
    mocks.fees.mockImplementation(async (recipient: string) => ({
      amount: recipient.toLowerCase() === wallet.toLowerCase() ? '10' : '0',
      block_number: 100,
    }))
    mocks.prepareFeeClaim.mockResolvedValue({
      chain_id: 84532, from: wallet, to: config.fee_escrow, data: '0x1234', value: '0', deadline: 0, review: 'claim',
    })
    mocks.sendTransactionAsync.mockResolvedValue(`0x${'dd'.repeat(32)}`)
    mocks.waitForTransactionReceipt.mockResolvedValue({ status: 'success' })
  })

  it('shows Creator, Protocol, Referrer, and LaaS separately and disables zero balances', async () => {
    render(<FeeBalances token={token} launchpad={launchpad} config={config} />)
    for (const role of ['Creator', 'Protocol', 'Referrer', 'LaaS']) expect(screen.getByText(role)).toBeInTheDocument()
    expect(await screen.findByTestId('fee-Creator')).toHaveTextContent('10')
    expect(screen.getByRole('button', { name: 'Claim Protocol' })).toBeDisabled()
  })

  it('refreshes balances after a successful Claim receipt', async () => {
    render(<FeeBalances token={token} launchpad={launchpad} config={config} />)
    fireEvent.click(await screen.findByRole('button', { name: 'Claim Creator' }))
    await waitFor(() => expect(mocks.waitForTransactionReceipt).toHaveBeenCalled())
    await waitFor(() => expect(mocks.fees.mock.calls.length).toBeGreaterThanOrEqual(8))
  })
})

describe('exact-input swap', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mocks.account.chainId = 84532
    mocks.balance = 1_000n
    mocks.quoteSwap.mockResolvedValue({ quote_id: 'quote-1', amount_in: '100', amount_out: '90', fee: '1', expires_at: 1_900_000_000 })
    mocks.prepareSwap.mockResolvedValue({
      chain_id: 84532, from: wallet, to: '0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee',
      data: '0xabcd', value: '100', deadline: 1_900_000_100, review: 'exact-input swap',
    })
    mocks.sendTransactionAsync.mockResolvedValue(`0x${'ff'.repeat(32)}`)
    mocks.waitForTransactionReceipt.mockResolvedValue({ status: 'success' })
  })

  it('rejects zero amount, insufficient balance, and wrong chain', async () => {
    const { rerender } = render(<SwapForm token={token} />)
    fireEvent.click(screen.getByRole('button', { name: '获取 Quote' }))
    expect(screen.getByText('金额必须大于 0')).toBeInTheDocument()
    fireEvent.change(screen.getByLabelText('输入金额（raw units）'), { target: { value: '1001' } })
    fireEvent.click(screen.getByRole('button', { name: '获取 Quote' }))
    expect(screen.getByText('余额不足')).toBeInTheDocument()
    mocks.account.chainId = 1
    rerender(<SwapForm token={token} />)
    expect(screen.getByRole('button', { name: '获取 Quote' })).toBeDisabled()
  })

  it('shows safety details then broadcasts the prepared transaction unchanged', async () => {
    render(<SwapForm token={token} />)
    fireEvent.change(screen.getByLabelText('输入金额（raw units）'), { target: { value: '100' } })
    fireEvent.click(screen.getByRole('button', { name: '获取 Quote' }))
    expect(await screen.findByText(/Fee：1/)).toBeInTheDocument()
    expect(screen.getByText(/Price Impact/)).toBeInTheDocument()
    expect(screen.queryByText(/exact-output/i)).not.toBeInTheDocument()
    fireEvent.click(screen.getByRole('button', { name: '准备交易' }))
    expect(await screen.findByText(/Slippage：1%/)).toBeInTheDocument()
    expect(screen.getByText(/anti-snipe/)).toBeInTheDocument()
    expect(screen.getByText(/Deadline：1900000100/)).toBeInTheDocument()
    fireEvent.click(screen.getByRole('button', { name: '确认并签名' }))
    await waitFor(() => expect(mocks.sendTransactionAsync).toHaveBeenCalledWith({
      to: '0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee', data: '0xabcd', value: 100n, chainId: 84532,
    }))
  })

  it('rejects an expired quote returned by prepare', async () => {
    mocks.prepareSwap.mockRejectedValue(new Error('quote missing or expired'))
    render(<SwapForm token={token} />)
    fireEvent.change(screen.getByLabelText('输入金额（raw units）'), { target: { value: '100' } })
    fireEvent.click(screen.getByRole('button', { name: '获取 Quote' }))
    fireEvent.click(await screen.findByRole('button', { name: '准备交易' }))
    expect(await screen.findByText('quote missing or expired')).toBeInTheDocument()
  })
})
