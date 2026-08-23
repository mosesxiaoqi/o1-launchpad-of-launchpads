import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import Page from './page'

const mocks = vi.hoisted(() => ({
  account: {
    address: '0x1111111111111111111111111111111111111111',
    chainId: 84532,
    isConnected: true,
  },
  signMessageAsync: vi.fn(),
  writeContractAsync: vi.fn(),
  waitForTransactionReceipt: vi.fn(),
  push: vi.fn(),
  config: vi.fn(),
  createChallenge: vi.fn(),
  verifySignature: vi.fn(),
  createLaunchpad: vi.fn(),
}))

vi.mock('wagmi', () => ({
  useAccount: () => mocks.account,
  useSignMessage: () => ({ signMessageAsync: mocks.signMessageAsync }),
  useWriteContract: () => ({ writeContractAsync: mocks.writeContractAsync }),
  usePublicClient: () => ({ waitForTransactionReceipt: mocks.waitForTransactionReceipt }),
}))

vi.mock('next/navigation', () => ({ useRouter: () => ({ push: mocks.push }) }))
vi.mock('@/lib/wagmi', () => ({ chain: { id: 84532 } }))
vi.mock('@/components/connect-wallet', () => ({ ConnectWallet: () => <button>连接钱包</button> }))

vi.mock('@/lib/api', () => ({
  launchpadApi: {
    config: mocks.config,
    createChallenge: mocks.createChallenge,
    verifySignature: mocks.verifySignature,
    createLaunchpad: mocks.createLaunchpad,
  },
}))

function fillValidForm() {
  fireEvent.change(screen.getByLabelText('名称'), { target: { value: 'My Pad' } })
  fireEvent.change(screen.getByLabelText('Slug'), { target: { value: 'My Pad' } })
  fireEvent.change(screen.getByLabelText('简介'), { target: { value: 'A test launchpad' } })
}

describe('create launchpad page', () => {
  beforeEach(() => {
    mocks.account.address = '0x1111111111111111111111111111111111111111'
    mocks.account.chainId = 84532
    mocks.account.isConnected = true
    vi.clearAllMocks()
    mocks.config.mockResolvedValue({
      chain_id: 84532,
      registry: '0x2222222222222222222222222222222222222222',
    })
    mocks.createChallenge.mockResolvedValue({ challenge_id: 'challenge', message: 'Sign me' })
    mocks.signMessageAsync.mockResolvedValue(`0x${'11'.repeat(65)}`)
    mocks.verifySignature.mockResolvedValue({ address: mocks.account.address })
    mocks.writeContractAsync.mockResolvedValue(`0x${'22'.repeat(32)}`)
    mocks.waitForTransactionReceipt.mockResolvedValue({ status: 'success' })
    mocks.createLaunchpad.mockResolvedValue({ slug: 'my-pad' })
  })

  it('disables creation without a connected wallet', () => {
    mocks.account.isConnected = false
    render(<Page />)
    expect(screen.getByRole('button', { name: '创建 Launchpad' })).toBeDisabled()
    expect(screen.getByText('请先连接钱包')).toBeInTheDocument()
  })

  it('disables creation on the wrong chain', () => {
    mocks.account.chainId = 1
    render(<Page />)
    expect(screen.getByRole('button', { name: '创建 Launchpad' })).toBeDisabled()
    expect(screen.getByText('请切换到 Base Sepolia')).toBeInTheDocument()
  })

  it('shows validation errors and focuses the first invalid field', async () => {
    render(<Page />)
    fireEvent.click(screen.getByRole('button', { name: '创建 Launchpad' }))
    await waitFor(() => expect(screen.getByLabelText('Slug')).toHaveFocus())
    expect(screen.getByText(/Slug 需要/)).toBeInTheDocument()
    expect(mocks.createChallenge).not.toHaveBeenCalled()
  })

  it('keeps form input when the wallet rejects the signature', async () => {
    mocks.signMessageAsync.mockRejectedValue(new Error('User rejected the request'))
    render(<Page />)
    fillValidForm()
    fireEvent.click(screen.getByRole('button', { name: '创建 Launchpad' }))
    expect(await screen.findByText('User rejected the request')).toBeInTheDocument()
    expect(screen.getByLabelText('名称')).toHaveValue('My Pad')
  })

  it('reports a reverted Registry transaction and keeps input', async () => {
    mocks.waitForTransactionReceipt.mockResolvedValue({ status: 'reverted' })
    render(<Page />)
    fillValidForm()
    fireEvent.click(screen.getByRole('button', { name: '创建 Launchpad' }))
    expect(await screen.findByText('Registry 交易已回滚')).toBeInTheDocument()
    expect(screen.getByLabelText('Slug')).toHaveValue('My Pad')
  })

  it('shows saving branding while API verification is pending', async () => {
    let resolve!: (value: { slug: string }) => void
    mocks.createLaunchpad.mockReturnValue(new Promise((done) => (resolve = done)))
    render(<Page />)
    fillValidForm()
    fireEvent.click(screen.getByRole('button', { name: '创建 Launchpad' }))
    expect(await screen.findByText('正在验证交易并保存品牌')).toBeInTheDocument()
    resolve({ slug: 'my-pad' })
    await waitFor(() => expect(mocks.push).toHaveBeenCalledWith('/launchpad/my-pad'))
  })

  it('signs in, creates on Registry, saves branding, and navigates', async () => {
    render(<Page />)
    fillValidForm()
    fireEvent.click(screen.getByRole('button', { name: '创建 Launchpad' }))
    await waitFor(() => expect(mocks.push).toHaveBeenCalledWith('/launchpad/my-pad'))
    expect(mocks.signMessageAsync).toHaveBeenCalledWith({ message: 'Sign me' })
    expect(mocks.writeContractAsync).toHaveBeenCalledWith(
      expect.objectContaining({ functionName: 'createLaunchpad' }),
    )
    expect(mocks.createLaunchpad).toHaveBeenCalledWith(
      expect.objectContaining({ slug: 'my-pad', registry_tx_hash: `0x${'22'.repeat(32)}` }),
    )
  })
})
