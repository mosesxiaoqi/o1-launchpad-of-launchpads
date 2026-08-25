import Link from 'next/link'
import { notFound } from 'next/navigation'

import { SiteHeader } from '@/components/site-header'
import { TokenCard } from '@/components/token-card'
import { ApiError, launchpadApi } from '@/lib/api'

export default async function LaunchpadPage({ params }: PageProps<'/launchpad/[slug]'>) {
  const { slug } = await params
  let launchpad
  let tokenPage
  try {
    ;[launchpad, tokenPage] = await Promise.all([
      launchpadApi.getLaunchpad(slug),
      launchpadApi.listLaunchpadTokens(slug),
    ])
  } catch (error) {
    if (error instanceof ApiError && error.status === 404) notFound()
    throw error
  }
  return (
    <main className="site-shell app-page" style={{ '--tenant-accent': launchpad.primary_color } as React.CSSProperties}>
      <SiteHeader />
      <section className="tenant-hero">
        <div className="tenant-identity">
          <div className="tenant-logo">
            {launchpad.name.slice(0, 2).toUpperCase()}
            {launchpad.logo_url && (
              // eslint-disable-next-line @next/next/no-img-element
              <img src={launchpad.logo_url} alt={`${launchpad.name} Logo`} />
            )}
          </div>
          <div><p className="eyebrow"><span>LAUNCHPAD</span> {launchpad.active ? 'ACTIVE' : 'INACTIVE'} TENANT</p><h1>{launchpad.name}</h1><p>{launchpad.description || '独立发行、永久流动性与链上费用分配。'}</p></div>
        </div>
        <Link className="button tenant-button" href={`/launchpad/${slug}/create-token`}>发行 Token <span aria-hidden="true">↗</span></Link>
      </section>
      <section className="tenant-meta" aria-label="Launchpad 信息">
        <div><span>SLUG</span><strong>/{launchpad.slug}</strong></div>
        <div><span>STATUS</span><strong className={launchpad.active ? 'status-ok' : 'status-inactive'}>● {launchpad.active ? 'ACTIVE' : 'INACTIVE'}</strong></div>
        <div><span>CHAIN</span><strong>BASE SEPOLIA</strong></div>
        <div><span>TREASURY</span><strong className="mono truncate">{launchpad.treasury}</strong></div>
      </section>
      <section className="asset-section" aria-label="Tokens">
        <div className="asset-heading"><div><p className="eyebrow"><span>ASSETS</span> DEPLOYED TOKENS</p><h2>发行资产</h2></div><span>{tokenPage.tokens.length.toString().padStart(2, '0')} TOTAL</span></div>
        {tokenPage.tokens.length ? <div className="token-grid">{tokenPage.tokens.map((token) => <TokenCard key={token.token} token={token} />)}</div> : <div className="empty-state"><span>＋</span><h3>暂无 Token</h3><p>创建第一个固定供应量资产并初始化交易池。</p><Link href={`/launchpad/${slug}/create-token`}>开始发行 →</Link></div>}
      </section>
    </main>
  )
}
