import { notFound } from 'next/navigation'

import { SiteHeader } from '@/components/site-header'
import { TokenLaunchForm } from '@/components/token-launch-form'
import { ApiError, launchpadApi } from '@/lib/api'

export default async function CreateTokenPage({ params }: PageProps<'/launchpad/[slug]/create-token'>) {
  const { slug } = await params
  let launchpad
  let config
  try {
    ;[launchpad, config] = await Promise.all([
      launchpadApi.getLaunchpad(slug),
      launchpadApi.config(),
    ])
  } catch (error) {
    if (error instanceof ApiError && error.status === 404) notFound()
    throw error
  }
  return (
    <main className="site-shell app-page">
      <SiteHeader backHref={`/launchpad/${slug}`} backLabel={launchpad.name} />
      <section className="page-intro compact">
        <p className="eyebrow"><span>LAUNCH / 02</span> FIXED SUPPLY ASSET</p>
        <h1>发行一个<br /><em>链上资产。</em></h1>
        <p>确认名称、Symbol 与元数据后，Factory 将部署 Token 并创建永久流动性。</p>
      </section>
      <TokenLaunchForm launchpad={launchpad} config={config} />
    </main>
  )
}
