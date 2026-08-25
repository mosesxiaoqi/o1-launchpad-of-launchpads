import { render, screen } from '@testing-library/react'
import { beforeEach, expect, it, vi } from 'vitest'

import Page from './page'
import { ApiError } from '@/lib/api'

const mocks = vi.hoisted(() => ({ getLaunchpad: vi.fn(), listLaunchpadTokens: vi.fn() }))

vi.mock('@/components/connect-wallet', () => ({ ConnectWallet: () => <button>连接钱包</button> }))
vi.mock('@/lib/api', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@/lib/api')>()
  return {
    ...actual,
    launchpadApi: {
      ...actual.launchpadApi,
      getLaunchpad: mocks.getLaunchpad,
      listLaunchpadTokens: mocks.listLaunchpadTokens,
    },
  }
})
vi.mock('next/navigation', () => ({ notFound: () => { throw new Error('NEXT_NOT_FOUND') } }))

beforeEach(() => {
  vi.clearAllMocks()
  mocks.getLaunchpad.mockResolvedValue({
    slug: 'my-pad', name: 'My Pad', description: 'My branded launchpad', logo_url: 'https://example.com/logo.png',
    primary_color: '#123456', active: true, treasury: '0x2222222222222222222222222222222222222222',
  })
  mocks.listLaunchpadTokens.mockResolvedValue({
    tokens: [{
      token: '0x1111111111111111111111111111111111111111', supply: '1000', status: 'confirmed',
      creator: '0x2222222222222222222222222222222222222222',
    }],
    next_cursor: '',
  })
})

it('renders tenant branding and trusted local token data', async () => {
  render(await Page({ params: Promise.resolve({ slug: 'my-pad' }) } as never))
  expect(screen.getByRole('heading', { level: 1, name: 'My Pad' })).toBeInTheDocument()
  expect(screen.getByText('My branded launchpad')).toBeInTheDocument()
  expect(screen.getByRole('img', { name: 'My Pad Logo' })).toHaveAttribute('src', 'https://example.com/logo.png')
  expect(screen.getByText('● ACTIVE')).toBeInTheDocument()
  expect(screen.getByText('0x1111111111111111111111111111111111111111')).toBeInTheDocument()
  expect(screen.getByRole('link', { name: '发行 Token' })).toHaveAttribute('href', '/launchpad/my-pad/create-token')
})

it('reports an inactive launchpad from API state', async () => {
  mocks.getLaunchpad.mockResolvedValue({
    slug: 'paused-pad', name: 'Paused Pad', description: '', logo_url: '', primary_color: '#123456',
    active: false, treasury: '0x2222222222222222222222222222222222222222',
  })
  render(await Page({ params: Promise.resolve({ slug: 'paused-pad' }) } as never))
  expect(screen.getByText('● INACTIVE')).toBeInTheDocument()
})

it('renders not found for an unknown slug', async () => {
  mocks.getLaunchpad.mockRejectedValue(new ApiError(404, 'not_found', 'not found'))
  await expect(Page({ params: Promise.resolve({ slug: 'missing' }) } as never)).rejects.toThrow('NEXT_NOT_FOUND')
})
