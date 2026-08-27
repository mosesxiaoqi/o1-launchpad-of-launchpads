import { fireEvent, render, screen } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import { ConnectWallet } from './connect-wallet'

const mocks = vi.hoisted(() => ({
  connect: vi.fn(),
  connectError: null as Error | null,
  connector: { id: 'injected', name: 'Injected', type: 'injected' },
  metaMaskConnector: { id: 'io.metamask', name: 'MetaMask', type: 'injected' },
  connectors: [] as Array<{ id: string; name: string; type: string }>,
}))

vi.mock('wagmi', () => ({
  useAccount: () => ({ address: undefined, chainId: undefined, isConnected: false }),
  useConnect: () => ({ connectors: mocks.connectors, mutate: mocks.connect, isPending: false, error: mocks.connectError }),
  useDisconnect: () => ({ mutate: vi.fn() }),
  useSwitchChain: () => ({ mutate: vi.fn(), isPending: false }),
}))

vi.mock('@/lib/wagmi', () => ({ chain: { id: 84532 } }))

describe('ConnectWallet', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mocks.connectError = null
    mocks.connectors = [mocks.connector]
  })

  it('requests Base Sepolia while connecting the selected wallet', () => {
    render(<ConnectWallet />)

    fireEvent.click(screen.getByRole('button', { name: '连接钱包' }))

    expect(mocks.connect).toHaveBeenCalledWith({
      connector: mocks.connector,
      chainId: 84532,
    })
  })

  it('shows the wallet error when the connection request fails', () => {
    mocks.connectError = new Error('Wallet request was rejected')

    render(<ConnectWallet />)

    expect(screen.getByRole('alert')).toHaveTextContent('Wallet request was rejected')
  })

  it('prefers the discovered MetaMask connector over the generic injected connector', () => {
    mocks.connectors = [mocks.connector, mocks.metaMaskConnector]

    render(<ConnectWallet />)
    fireEvent.click(screen.getByRole('button', { name: '连接钱包' }))

    expect(mocks.connect).toHaveBeenCalledWith({
      connector: mocks.metaMaskConnector,
      chainId: 84532,
    })
  })
})
