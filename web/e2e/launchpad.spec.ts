import { expect, test } from '@playwright/test'

const wallet = '0x1111111111111111111111111111111111111111'
const registryTx = `0x${'bb'.repeat(32)}`
const launchTx = `0x${'cc'.repeat(32)}`
const token = '0x8888888888888888888888888888888888888888'

test.beforeEach(async ({ page, request }) => {
  await request.post('http://127.0.0.1:3999/__reset')
  await page.addInitScript(({ wallet, registryTx, launchTx }) => {
    let sends = 0
    const listeners = new Map<string, Set<(...args: unknown[]) => void>>()
    const provider = {
      isMetaMask: true,
      async request({ method }: { method: string; params?: unknown[] }) {
        if (method === 'eth_accounts' || method === 'eth_requestAccounts') return [wallet]
        if (method === 'eth_chainId') return '0x14a34'
        if (method === 'personal_sign') return `0x${'11'.repeat(65)}`
        if (method === 'eth_sendTransaction') return sends++ === 0 ? registryTx : launchTx
        if (method === 'wallet_switchEthereumChain') return null
        throw new Error(`unsupported wallet method ${method}`)
      },
      on(event: string, listener: (...args: unknown[]) => void) {
        const eventListeners = listeners.get(event) || new Set()
        eventListeners.add(listener)
        listeners.set(event, eventListeners)
        return provider
      },
      removeListener(event: string, listener: (...args: unknown[]) => void) {
        listeners.get(event)?.delete(listener)
        return provider
      },
    }
    Object.defineProperty(window, 'ethereum', { configurable: true, value: provider })
  }, { wallet, registryTx, launchTx })
})

test('creates a tenant, launches a token, reconciles the receipt, and lists it under that tenant', async ({ page, request }) => {
  await page.goto('/create-launchpad')
  await page.getByRole('button', { name: '连接钱包' }).click()
  await expect(page.getByRole('button', { name: /0x1111.*断开/ })).toBeVisible()

  await page.getByLabel('Slug').fill('My Pad')
  await page.getByLabel('名称').fill('My Pad')
  await page.getByLabel('简介').fill('Local E2E tenant')
  await page.getByRole('button', { name: '创建 Launchpad' }).click()
  await expect(page).toHaveURL(/\/launchpad\/my-pad$/)
  await expect(page.getByRole('heading', { name: 'My Pad' })).toBeVisible()
  await expect(page.getByText('Local E2E tenant')).toBeVisible()

  await page.getByRole('link', { name: '发行 Token' }).click()
  await page.getByLabel('Token 名称').fill('Demo Token')
  await page.getByLabel('Symbol').fill('DEMO')
  await page.getByRole('button', { name: '生成发行计划' }).click()
  await expect(page.getByText('1.5% 总费率')).toBeVisible()
  await page.getByRole('button', { name: '确认并签名' }).click()
  await expect(page.getByRole('status')).toHaveText('等待索引确认')
  await expect(page.getByText(`Token：${token}`)).toBeVisible()
  await request.post('http://127.0.0.1:3999/__confirm-index')
  await expect(page.getByRole('status')).toHaveText('发行已确认')

  await page.getByRole('link', { name: 'My Pad' }).click()
  await expect(page.getByRole('heading', { name: token })).toBeVisible()
  await expect(page.getByRole('link', { name: '查看与交易' })).toHaveAttribute('href', `/token/${token}`)

  await page.goto('/launchpad/other-pad')
  await expect(page.getByRole('heading', { name: 'Other Pad' })).toBeVisible()
  await expect(page.getByText('暂无 Token')).toBeVisible()
  await expect(page.getByText(token)).not.toBeVisible()
})
