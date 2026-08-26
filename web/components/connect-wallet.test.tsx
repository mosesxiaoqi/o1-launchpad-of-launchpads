import { fireEvent, render, screen } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import { ConnectWallet } from './connect-wallet'

const mocks = vi.hoisted(() => ({
  connect: vi.fn(),
  connector: { id: 'injected', name: 'Injected', type: 'injected' },
}))

vi.mock('wagmi', () => ({
  useAccount: () => ({ address: undefined, chainId: undefined, isConnected: false }),
  useConnect: () => ({ connectors: [mocks.connector], mutate: mocks.connect, isPending: false }),
  useDisconnect: () => ({ mutate: vi.fn() }),
  useSwitchChain: () => ({ mutate: vi.fn(), isPending: false }),
}))

vi.mock('@/lib/wagmi', () => ({ chain: { id: 84532 } }))

describe('ConnectWallet', () => {
  beforeEach(() => vi.clearAllMocks())

  it('requests Base Sepolia while connecting the selected wallet', () => {
    render(<ConnectWallet />)

    fireEvent.click(screen.getByRole('button', { name: '连接钱包' }))

    expect(mocks.connect).toHaveBeenCalledWith({
      connector: mocks.connector,
      chainId: 84532,
    })
  })
})
