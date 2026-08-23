import { encodeAbiParameters, keccak256, type Hex } from 'viem'

const normalizedSlugPattern = /^[a-z0-9]+(?:-[a-z0-9]+)*$/

export function normalizeSlug(value: string): string {
  const slug = value.trim().toLowerCase().replace(/[\s-]+/g, '-')
  if (slug.length < 3 || slug.length > 32 || !normalizedSlugPattern.test(slug)) {
    throw new Error('slug must be 3-32 lowercase letters, numbers, or single internal hyphens')
  }
  return slug
}

export function deriveLaunchpadId(chainId: number, normalizedSlug: string): Hex {
  const slug = normalizeSlug(normalizedSlug)
  return keccak256(
    encodeAbiParameters([{ type: 'uint256' }, { type: 'string' }], [BigInt(chainId), slug]),
  )
}
