import { notFound } from 'next/navigation'
import Link from 'next/link'

import { ConnectWallet } from '@/components/connect-wallet'
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
    <main>
      <header><Link href={`/launchpad/${slug}`}>{launchpad.name}</Link><ConnectWallet /></header>
      <h1>发行 Token</h1>
      <TokenLaunchForm launchpad={launchpad} config={config} />
    </main>
  )
}
