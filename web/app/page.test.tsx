import { render, screen } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'

import Home from './page'

vi.mock('@/components/connect-wallet', () => ({ ConnectWallet: () => <button>连接钱包</button> }))

describe('home page', () => {
  it('presents the launch workflow and immutable protocol facts', () => {
    render(<Home />)

    expect(screen.getByRole('heading', { name: '发行你的 Launchpad。掌控每一次启动。' })).toBeInTheDocument()
    expect(screen.getByRole('link', { name: '创建 Launchpad' })).toHaveAttribute('href', '/create-launchpad')
    expect(screen.getByText('1,000,000,000')).toBeInTheDocument()
    expect(screen.getByText('1.5%')).toBeInTheDocument()
    expect(screen.getByText('16 秒')).toBeInTheDocument()
    expect(screen.getByText('永久流动性')).toBeInTheDocument()
    for (const step of ['建立品牌', '发行资产', '开放交易']) {
      expect(screen.getByRole('heading', { name: step })).toBeInTheDocument()
    }
    expect(screen.getByText('Base Sepolia')).toBeInTheDocument()
  })
})
