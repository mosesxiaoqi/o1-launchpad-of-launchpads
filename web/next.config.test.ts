import { describe, expect, it } from 'vitest'

import nextConfig from './next.config'

describe('next config', () => {
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
