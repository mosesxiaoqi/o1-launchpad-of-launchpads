import { afterEach, describe, expect, it, vi } from 'vitest'

import { api } from './api'

afterEach(() => vi.unstubAllGlobals())

describe('api', () => {
  it('always includes the session cookie', async () => {
	const fetchMock = vi.fn().mockResolvedValue(
	  new Response(JSON.stringify({ status: 'ok' }), {
		status: 200,
		headers: { 'content-type': 'application/json' },
	  }),
	)
	vi.stubGlobal('fetch', fetchMock)
	await api('/v1/health')
	expect(fetchMock).toHaveBeenCalledWith(
	  expect.stringContaining('/v1/health'),
	  expect.objectContaining({ credentials: 'include' }),
	)
  })

  it('throws the code and status from a problem response', async () => {
	vi.stubGlobal(
	  'fetch',
	  vi.fn().mockResolvedValue(
		new Response(
		  JSON.stringify({
			type: 'about:blank',
			title: 'Conflict',
			status: 409,
			detail: 'slug exists',
			code: 'conflict',
		  }),
		  { status: 409, headers: { 'content-type': 'application/problem+json' } },
		),
	  ),
	)
	await expect(api('/v1/launchpads')).rejects.toMatchObject({
	  status: 409,
	  code: 'conflict',
	  message: 'slug exists',
	})
  })
})
