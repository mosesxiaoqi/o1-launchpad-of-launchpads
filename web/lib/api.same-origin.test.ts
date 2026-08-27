import { afterEach, expect, it, vi } from 'vitest'

afterEach(() => {
  vi.unstubAllEnvs()
  vi.unstubAllGlobals()
  vi.resetModules()
})

it('uses the current origin when NEXT_PUBLIC_API_URL is empty', async () => {
  vi.stubEnv('NEXT_PUBLIC_API_URL', '')
  const fetchMock = vi.fn().mockResolvedValue(
    new Response(JSON.stringify({ status: 'ok' }), {
      status: 200,
      headers: { 'content-type': 'application/json' },
    }),
  )
  vi.stubGlobal('fetch', fetchMock)

  const { api } = await import('./api')
  await api('/v1/health')

  expect(fetchMock).toHaveBeenCalledWith('/v1/health', expect.any(Object))
})
