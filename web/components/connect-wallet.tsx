'use client'

import { useAccount, useConnect, useDisconnect, useSwitchChain } from 'wagmi'

import { chain } from '@/lib/wagmi'

export function useBaseSepoliaWriteReady() {
  const { chainId, isConnected } = useAccount()
  return isConnected && chainId === chain.id
}

export function ConnectWallet() {
  const { address, chainId, isConnected } = useAccount()
  const { connectors, mutate: connect, isPending: isConnecting, error: connectError } = useConnect()
  const { mutate: disconnect } = useDisconnect()
  const { mutate: switchChain, isPending: isSwitching } = useSwitchChain()

  if (!isConnected) {
    const connector = connectors.find(({ id, name }) => id === 'io.metamask' || name.toLowerCase().includes('metamask')) ?? connectors[0]
    return (
      <div>
        <button disabled={!connector || isConnecting} onClick={() => connector && connect({ connector, chainId: chain.id })}>
          {isConnecting ? '连接中…' : '连接钱包'}
        </button>
        {connectError && <p role="alert" className="form-error">{connectError.message}</p>}
      </div>
    )
  }

  if (chainId !== chain.id) {
    return (
      <button disabled={isSwitching} onClick={() => switchChain({ chainId: chain.id })}>
        {isSwitching ? '切换中…' : '切换到 Base Sepolia'}
      </button>
    )
  }

  return (
    <button onClick={() => disconnect()} title={address}>
      {address?.slice(0, 6)}…{address?.slice(-4)} · 断开
    </button>
  )
}
