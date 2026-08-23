import { describe, expect, it } from 'vitest'

import { deriveLaunchpadId, normalizeSlug } from './launchpad-id'

describe('normalizeSlug', () => {
  it('normalizes case, spaces, and repeated separators', () => {
    expect(normalizeSlug('  My  Launch--Pad  ')).toBe('my-launch-pad')
  })

  it.each(['ab', '-launch', 'launch-', 'launch_pad', '发射台', 'a'.repeat(33)])(
    'rejects %s',
    (value) => expect(() => normalizeSlug(value)).toThrow(),
  )
})

it('derives the same id for the same chain and normalized slug', () => {
  const id = deriveLaunchpadId(84532, 'my-launch-pad')
  expect(id).toMatch(/^0x[0-9a-f]{64}$/)
  expect(deriveLaunchpadId(84532, 'my-launch-pad')).toBe(id)
  expect(deriveLaunchpadId(1, 'my-launch-pad')).not.toBe(id)
})
