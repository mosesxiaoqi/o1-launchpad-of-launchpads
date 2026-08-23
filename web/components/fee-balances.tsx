'use client'

import { useCallback, useEffect, useState } from 'react'
import { useAccount, usePublicClient, useSendTransaction } from 'wagmi'

import { launchpadApi, type Config, type Launchpad, type Token } from '@/lib/api'
import { chain } from '@/lib/wagmi'

type Role = 'Creator' | 'Protocol' | 'Referrer' | 'LaaS'
type Balances = Record<Role, string>

const emptyBalances: Balances = { Creator: '0', Protocol: '0', Referrer: '0', LaaS: '0' }

export function FeeBalances({ token, launchpad, config }: { token: Token; launchpad: Launchpad; config: Config }) {
  const { address, chainId, isConnected } = useAccount()
  const { sendTransactionAsync } = useSendTransaction()
  const publicClient = usePublicClient()
  const [currency, setCurrency] = useState<string>(token.quote)
  const [balances, setBalances] = useState<Balances>(emptyBalances)
  const [busy, setBusy] = useState<Role | null>(null)
  const [error, setError] = useState('')
  const [status, setStatus] = useState('')

  const recipient = useCallback((role: Role) => {
    if (role === 'Creator') return token.creator
    if (role === 'Protocol') return launchpad.treasury
    if (role === 'LaaS') return config.laas_treasury
    return address || '0x0000000000000000000000000000000000000000'
  }, [address, config.laas_treasury, launchpad.treasury, token.creator])

  const loadBalances = useCallback(async () => {
    const roles: Role[] = ['Creator', 'Protocol', 'Referrer', 'LaaS']
    const results = await Promise.all(roles.map(async (role) => [role, (await launchpadApi.fees(recipient(role), currency)).amount] as const))
    return Object.fromEntries(results) as Balances
  }, [currency, recipient])

  useEffect(() => {
    let cancelled = false
    async function load() {
      try {
        const next = await loadBalances()
        if (!cancelled) setBalances(next)
      } catch (cause) {
        if (!cancelled) setError(message(cause))
      }
    }
    void load()
    return () => { cancelled = true }
  }, [loadBalances])

  async function claim(role: Role) {
    if (!address || !publicClient) return
    setBusy(role)
    setError('')
    setStatus('正在准备领取交易')
    try {
      const plan = await launchpadApi.prepareFeeClaim({
        chain_id: chain.id,
        currency,
        recipient: recipient(role),
        wallet: address,
      })
      const hash = await sendTransactionAsync({
        to: plan.to,
        data: plan.data,
        value: BigInt(plan.value),
        chainId: chain.id,
      })
      setStatus('等待领取交易确认')
      const receipt = await publicClient.waitForTransactionReceipt({ hash })
      if (receipt.status !== 'success') throw new Error('领取交易已回滚')
      setBalances(await loadBalances())
      setStatus('领取成功，余额已刷新')
    } catch (cause) {
      setError(message(cause))
      setStatus('')
    } finally {
      setBusy(null)
    }
  }

  const roles: Role[] = ['Creator', 'Protocol', 'Referrer', 'LaaS']
  return (
    <section aria-labelledby="fees-title">
      <h2 id="fees-title">费用余额</h2>
      <label htmlFor="fee-currency">币种</label>
      <select id="fee-currency" value={currency} onChange={(event) => setCurrency(event.target.value)}>
        <option value={token.quote}>Quote</option>
        <option value={token.token}>Token</option>
      </select>
      {roles.map((role) => {
        const amount = balances[role]
        const ownsRecipient = Boolean(address && recipient(role).toLowerCase() === address.toLowerCase())
        const canClaim = isConnected && chainId === chain.id && ownsRecipient && BigInt(amount || '0') > 0n && !busy
        return (
          <article key={role}>
            <h3>{role}</h3>
            <p>{recipient(role)}</p>
            <p data-testid={`fee-${role}`}>{amount}</p>
            <button type="button" aria-label={`Claim ${role}`} disabled={!canClaim} onClick={() => claim(role)}>
              {busy === role ? '领取中…' : 'Claim'}
            </button>
          </article>
        )
      })}
      {status && <p role="status">{status}</p>}
      {error && <p role="alert">{error}</p>}
    </section>
  )
}

function message(cause: unknown) {
  return cause instanceof Error ? cause.message : '费用操作失败'
}
