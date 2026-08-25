import Link from 'next/link'

import { ConnectWallet } from './connect-wallet'

export function SiteHeader({ backHref, backLabel }: { backHref?: string; backLabel?: string }) {
  return (
    <header className="site-header">
      <Link className="brand" href={backHref || '/'} aria-label={backLabel ? `返回 ${backLabel}` : 'O1 Launchpad 首页'}>
        <span className="brand-mark">O1</span>
        <span className="brand-copy">
          <strong>{backLabel || 'LAUNCHPAD'}</strong>
          <small>OF LAUNCHPADS</small>
        </span>
      </Link>
      <div className="header-actions">
        <span className="network-pill"><i /> BASE SEPOLIA <b>84532</b></span>
        <ConnectWallet />
      </div>
    </header>
  )
}
