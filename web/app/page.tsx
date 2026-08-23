import { ConnectWallet } from '@/components/connect-wallet'
import Link from 'next/link'

export default function Home() {
  return (
    <main>
      <h1>O1 Launchpad of Launchpads</h1>
      <p>Base Sepolia MVP</p>
      <ConnectWallet />
      <Link href="/create-launchpad">创建 Launchpad</Link>
    </main>
  )
}
