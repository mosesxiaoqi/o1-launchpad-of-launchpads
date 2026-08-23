import { expect, it, vi } from 'vitest'

import { pollTransaction } from './transaction-poll'

it('polls with bounded exponential backoff until confirmed', async () => {
  const read = vi
    .fn()
    .mockRejectedValueOnce(new Error('not indexed'))
    .mockResolvedValueOnce({ status: 'confirming' })
    .mockResolvedValueOnce({ status: 'confirmed' })
  const sleep = vi.fn().mockResolvedValue(undefined)
  await expect(pollTransaction(read, { sleep, attempts: 4, initialDelayMs: 100 })).resolves.toEqual({
    status: 'confirmed',
  })
  expect(sleep.mock.calls.map(([delay]) => delay)).toEqual([100, 200])
})

it('stops after the configured number of attempts', async () => {
  const read = vi.fn().mockResolvedValue({ status: 'confirming' })
  await expect(
    pollTransaction(read, { sleep: vi.fn().mockResolvedValue(undefined), attempts: 2 }),
  ).rejects.toThrow('transaction indexing timed out')
  expect(read).toHaveBeenCalledTimes(2)
})
