import type { NextConfig } from 'next'

const nextConfig: NextConfig = {
  turbopack: { root: process.cwd() },
  async rewrites() {
    return [
      {
        source: '/v1/:path*',
        destination: 'http://127.0.0.1:8888/v1/:path*',
      },
      {
        source: '/health/:path*',
        destination: 'http://127.0.0.1:8888/health/:path*',
      },
    ]
  },
}

export default nextConfig
