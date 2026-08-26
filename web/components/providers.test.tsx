import { render, screen } from '@testing-library/react'
import type { ReactNode } from 'react'
import { describe, expect, it, vi } from 'vitest'

import { Providers } from './providers'

const mocks = vi.hoisted(() => ({
  reconnectOnMount: undefined as boolean | undefined,
}))

vi.mock('wagmi', () => ({
  WagmiProvider: ({ children, reconnectOnMount }: { children: ReactNode; reconnectOnMount?: boolean }) => {
    mocks.reconnectOnMount = reconnectOnMount
    return children
  },
}))

vi.mock('@/lib/wagmi', () => ({ wagmiConfig: {} }))

describe('Providers', () => {
  it('does not reconnect a wallet when the page opens', () => {
    render(<Providers><span>页面内容</span></Providers>)

    expect(screen.getByText('页面内容')).toBeInTheDocument()
    expect(mocks.reconnectOnMount).toBe(false)
  })
})
