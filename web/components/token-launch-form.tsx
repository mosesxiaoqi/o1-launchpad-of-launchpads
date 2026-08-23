'use client'

import { bytesToHex } from 'viem'
import { useState, type FormEvent } from 'react'
import { useAccount, usePublicClient, useSendTransaction } from 'wagmi'

import { launchpadApi, type Config, type Launchpad, type PreparedTransaction } from '@/lib/api'
import { parseLaunchedReceipt, type LaunchedResult } from '@/lib/receipt'
import { pollTransaction } from '@/lib/transaction-poll'
import { chain } from '@/lib/wagmi'

type Phase = 'idle' | 'preparing' | 'review' | 'wallet' | 'receipt' | 'indexing' | 'confirmed'

export function TokenLaunchForm({ launchpad, config }: { launchpad: Launchpad; config: Config }) {
  const { address, chainId, isConnected } = useAccount()
  const { sendTransactionAsync } = useSendTransaction()
  const publicClient = usePublicClient()
  const [name, setName] = useState('')
  const [symbol, setSymbol] = useState('')
  const [contractURI, setContractURI] = useState('')
  const [plan, setPlan] = useState<PreparedTransaction | null>(null)
  const [phase, setPhase] = useState<Phase>('idle')
  const [error, setError] = useState('')
  const [launched, setLaunched] = useState<LaunchedResult | null>(null)
  const ready = Boolean(address && isConnected && chainId === chain.id && publicClient)

  async function prepare(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    if (!address || !ready) return
    if (!name.trim() || !symbol.trim()) {
      setError('Token 名称和 Symbol 必填')
      return
    }
    setError('')
    setPhase('preparing')
    try {
      const salt = bytesToHex(crypto.getRandomValues(new Uint8Array(32)))
      const prepared = await launchpadApi.prepareLaunch({
        chain_id: chain.id,
        launchpad_id: launchpad.id,
        name: name.trim(),
        symbol: symbol.trim(),
        contract_uri: contractURI,
        salt,
        wallet: address,
      })
      setPlan(prepared)
      setPhase('review')
    } catch (cause) {
      setPhase('idle')
      setError(message(cause))
    }
  }

  async function broadcast() {
    if (!plan || !address || !publicClient) return
    setError('')
    setPhase('wallet')
    try {
      const hash = await sendTransactionAsync({
        to: plan.to,
        data: plan.data,
        value: BigInt(plan.value),
        chainId: chain.id,
      })
      setPhase('receipt')
      const receipt = await publicClient.waitForTransactionReceipt({ hash })
      if (receipt.status !== 'success') throw new Error('Launch 交易已回滚')
      const result = parseLaunchedReceipt(receipt, config.factory)
      setLaunched(result)
      setPhase('indexing')
      const indexed = await pollTransaction(() => launchpadApi.transaction(chain.id, hash))
      if (indexed.status === 'reverted') throw new Error(indexed.failure_reason || 'Launch 已回滚')
      setPhase('confirmed')
    } catch (cause) {
      if (/stale_plan|StaleConfig|LaunchExpired/i.test(message(cause))) {
        invalidatePlan()
      } else {
        setPhase(plan ? 'review' : 'idle')
        setError(message(cause))
      }
    }
  }

  function invalidatePlan() {
    setPlan(null)
    setPhase('idle')
    setError('计划已过期，请重新生成并审阅')
  }

  return (
    <section>
      <form onSubmit={prepare}>
        <label htmlFor="token-name">Token 名称</label>
        <input id="token-name" value={name} onChange={(event) => setName(event.target.value)} />
        <label htmlFor="token-symbol">Symbol</label>
        <input id="token-symbol" value={symbol} onChange={(event) => setSymbol(event.target.value)} />
        <label htmlFor="contract-uri">Contract URI</label>
        <input id="contract-uri" value={contractURI} onChange={(event) => setContractURI(event.target.value)} />
        {!plan && <button disabled={!ready || phase === 'preparing'}>生成发行计划</button>}
      </form>

      {plan && (
        <section aria-labelledby="launch-review-title">
          <h2 id="launch-review-title">发行确认</h2>
          <dl>
            <dt>固定供应量</dt><dd>1,000,000,000</dd>
            <dt>流动性</dt><dd>永久流动性</dd>
            <dt>协议费</dt><dd>1% 协议费</dd>
            <dt>LaaS 费</dt><dd>0.5% LaaS</dd>
            <dt>总费率</dt><dd>1.5% 总费率</dd>
            <dt>保护期</dt><dd>16 秒 anti-snipe</dd>
            <dt>Quote</dt><dd>{config.quote}</dd>
            <dt>Launchpad ID</dt><dd>{launchpad.id}</dd>
            <dt>Treasury</dt><dd>{launchpad.treasury}</dd>
          </dl>
          <p>{plan.review}</p>
          <button disabled={phase !== 'review'} onClick={broadcast}>确认并签名</button>
        </section>
      )}

      {phase === 'preparing' && <p role="status">正在生成最新发行计划</p>}
      {phase === 'wallet' && <p role="status">等待钱包确认交易</p>}
      {phase === 'receipt' && <p role="status">等待链上 Receipt</p>}
      {phase === 'indexing' && <p role="status">等待索引确认</p>}
      {phase === 'confirmed' && <p role="status">发行已确认</p>}
      {launched && <p>Token：{launched.token}</p>}
      {error && <p role="alert">{error}</p>}
      {!ready && <p>请连接钱包并切换到 Base Sepolia</p>}
    </section>
  )
}

function message(cause: unknown) {
  return cause instanceof Error ? cause.message : 'Token 发行失败'
}
