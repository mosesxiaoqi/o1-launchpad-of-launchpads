import Link from 'next/link'
import { notFound } from 'next/navigation'

import { ConnectWallet } from '@/components/connect-wallet'
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
    <main>
      <header><Link href={`/launchpad/${token.launchpad_slug}`}>{launchpad.name}</Link><ConnectWallet /></header>
      <h1>Token</h1>
      <dl>
        <dt>Address</dt><dd>{token.token}</dd>
        <dt>Pool ID</dt><dd>{token.pool_id}</dd>
        <dt>Status</dt><dd>{token.status}</dd>
        <dt>Supply</dt><dd>{token.supply}</dd>
      </dl>
      <FeeBalances token={token} launchpad={launchpad} config={config} />
      <SwapForm token={token} />
    </main>
  )
}
