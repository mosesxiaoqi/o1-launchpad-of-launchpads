import Link from 'next/link'

import type { Token } from '@/lib/api'

export function TokenCard({ token }: { token: Token }) {
  return (
    <article className="token-card">
      <div className="token-card-top"><span className="token-symbol">{token.token.slice(2, 4).toUpperCase()}</span><span className="status-ok">● {token.status.toUpperCase()}</span></div>
      <h2 aria-label={token.token}>{token.token.slice(0, 8)}…{token.token.slice(-6)}</h2>
      <p className="mono token-address">{token.token}</p>
      <dl className="token-card-facts">
        <dt>Supply</dt><dd>{token.supply}</dd>
        <dt>Status</dt><dd>{token.status}</dd>
        <dt>Creator</dt><dd>{token.creator}</dd>
      </dl>
      <Link className="card-link" href={`/token/${token.token}`}>查看与交易 <span aria-hidden="true">↗</span></Link>
    </article>
  )
}
