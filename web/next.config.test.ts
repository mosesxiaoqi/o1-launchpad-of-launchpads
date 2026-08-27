import { describe, expect, it, vi } from 'vitest'

import nextConfig from './next.config'

describe('next config', () => {
  it('allows the configured public frontend origin to load development assets', async () => {
    vi.stubEnv('FRONTEND_ORIGIN', 'https://demo-tunnel.trycloudflare.com')
    vi.resetModules()

    const { default: tunnelConfig } = await import('./next.config')

    expect(tunnelConfig.allowedDevOrigins).toEqual(['demo-tunnel.trycloudflare.com'])
    vi.unstubAllEnvs()
  })

  it('proxies browser API requests to the loopback gateway', async () => {
    expect(nextConfig.rewrites).toBeTypeOf('function')
    if (!nextConfig.rewrites) return

    await expect(nextConfig.rewrites()).resolves.toEqual([
      {
        source: '/v1/:path*',
        destination: 'http://127.0.0.1:8888/v1/:path*',
      },
      {
        source: '/health/:path*',
        destination: 'http://127.0.0.1:8888/health/:path*',
      },
    ])
  })
})
