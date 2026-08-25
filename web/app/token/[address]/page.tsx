import { notFound } from 'next/navigation'

import { SiteHeader } from '@/components/site-header'
import { FeeBalances } from '@/components/fee-balances'
import { SwapForm } from '@/components/swap-form'
import { ApiError, launchpadApi } from '@/lib/api'

export default async function TokenPage({ params }: PageProps<'/token/[address]'>) {
  const { address } = await params
  let token
  let config
  try {
    ;[token, config] = await Promise.all([launchpadApi.getToken(address), launchpadApi.config()])
  } catch (error) {
    if (error instanceof ApiError && error.status === 404) notFound()
    throw error
  }

  let launchpad
  try {
    launchpad = await launchpadApi.getLaunchpad(token.launchpad_slug)
  } catch (error) {
    if (error instanceof ApiError && error.status === 404) notFound()
    throw error
  }

  return (
    <main className="site-shell app-page token-page">
      <SiteHeader backHref={`/launchpad/${token.launchpad_slug}`} backLabel={launchpad.name} />
      <section className="asset-hero">
        <div className="token-avatar">{token.token.slice(2, 4).toUpperCase()}</div>
        <div className="asset-title"><p className="eyebrow"><span>ERC-20</span> FIXED SUPPLY</p><h1>{token.token.slice(0, 8)}…{token.token.slice(-6)}</h1><p className="mono">{token.token}</p></div>
        <span className="asset-status">● {token.status.toUpperCase()}</span>
      </section>
      <section className="asset-facts" aria-label="Token 信息">
        <div><span>SUPPLY</span><strong>{token.supply}</strong></div>
        <div><span>QUOTE</span><strong>{token.quote.toLowerCase() === '0x0000000000000000000000000000000000000000' ? 'NATIVE ETH' : token.quote}</strong></div>
        <div><span>POOL ID</span><strong className="mono truncate">{token.pool_id}</strong></div>
        <div><span>CREATOR</span><strong className="mono truncate">{token.creator}</strong></div>
      </section>
      <div className="trade-layout">
        <div className="ledger-column">
          <div className="section-title"><p className="eyebrow"><span>LEDGER</span> FEE ESCROW</p><h2>费用账本</h2></div>
          <FeeBalances token={token} launchpad={launchpad} config={config} />
        </div>
        <aside className="swap-column"><SwapForm token={token} /></aside>
      </div>
    </main>
  )
}
