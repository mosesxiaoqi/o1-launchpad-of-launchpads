import Link from 'next/link'
import { notFound } from 'next/navigation'

import { ConnectWallet } from '@/components/connect-wallet'
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
    <main style={{ borderTop: `6px solid ${launchpad.primary_color}` }}>
      <header><Link href="/">O1 Launchpad</Link><ConnectWallet /></header>
      <h1>{launchpad.name}</h1>
      <p>{launchpad.description}</p>
      <Link href={`/launchpad/${slug}/create-token`}>发行 Token</Link>
      <section aria-label="Tokens">
        {tokenPage.tokens.length ? tokenPage.tokens.map((token) => <TokenCard key={token.token} token={token} />) : <p>暂无 Token</p>}
      </section>
    </main>
  )
}
