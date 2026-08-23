import { ConnectWallet } from '@/components/connect-wallet'
import { LaunchpadForm } from '@/components/launchpad-form'
import Link from 'next/link'

export default function CreateLaunchpadPage() {
  return (
    <main>
      <header>
        <Link href="/">O1 Launchpad</Link>
        <ConnectWallet />
      </header>
      <h1>创建 Launchpad</h1>
      <p>在 Base Sepolia 注册租户，并保存独立品牌配置。</p>
      <LaunchpadForm />
    </main>
  )
}
