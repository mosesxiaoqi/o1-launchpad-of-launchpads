import { defineConfig } from '@playwright/test'

export default defineConfig({
  testDir: './e2e',
  fullyParallel: false,
  workers: 1,
  reporter: 'line',
  use: {
    baseURL: 'http://127.0.0.1:3100',
    trace: 'retain-on-failure',
  },
  webServer: [
    {
      command: 'node e2e/mock-api.mjs',
      url: 'http://127.0.0.1:3999/health/ready',
      reuseExistingServer: !process.env.CI,
    },
    {
      command: 'NEXT_PUBLIC_API_URL=http://127.0.0.1:3999 NEXT_PUBLIC_RPC_URL=http://127.0.0.1:3999/rpc npm run dev -- --hostname 127.0.0.1 --port 3100',
      url: 'http://127.0.0.1:3100',
      reuseExistingServer: !process.env.CI,
    },
  ],
})
