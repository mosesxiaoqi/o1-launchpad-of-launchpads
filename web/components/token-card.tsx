import Link from 'next/link'

import type { Token } from '@/lib/api'

export function TokenCard({ token }: { token: Token }) {
  return (
    <article>
      <h2>{token.token}</h2>
      <dl>
        <dt>Supply</dt><dd>{token.supply}</dd>
        <dt>Status</dt><dd>{token.status}</dd>
        <dt>Creator</dt><dd>{token.creator}</dd>
      </dl>
      <Link href={`/token/${token.token}`}>查看与交易</Link>
    </article>
  )
}
