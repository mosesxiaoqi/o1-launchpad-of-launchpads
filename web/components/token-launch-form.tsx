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
    <section className="token-launcher">
      <div className="form-panel">
        <form className="stack-form launch-form" onSubmit={prepare}>
          <div className="form-heading"><span>TOKEN CONFIG</span><strong>资产标识与元数据</strong></div>
          <div className="form-field"><label htmlFor="token-name">Token 名称</label><input id="token-name" value={name} onChange={(event) => setName(event.target.value)} placeholder="Example Token" /></div>
          <div className="form-field"><label htmlFor="token-symbol">Symbol</label><input id="token-symbol" value={symbol} onChange={(event) => setSymbol(event.target.value)} placeholder="TOKEN" /></div>
          <div className="form-field"><label htmlFor="contract-uri">Contract URI</label><input id="contract-uri" value={contractURI} onChange={(event) => setContractURI(event.target.value)} placeholder="ipfs://…" /></div>
          {!plan && <button className="primary-action" disabled={!ready || phase === 'preparing'}>生成发行计划 <span aria-hidden="true">→</span></button>}
        </form>
        {phase === 'preparing' && <p className="inline-status" role="status">正在生成最新发行计划</p>}
        {phase === 'wallet' && <p className="inline-status" role="status">等待钱包确认交易</p>}
        {phase === 'receipt' && <p className="inline-status" role="status">等待链上 Receipt</p>}
        {phase === 'indexing' && <p className="inline-status" role="status">等待索引确认</p>}
        {phase === 'confirmed' && <p className="inline-status success" role="status">发行已确认</p>}
        {launched && <p className="mono launched-token">Token：{launched.token}</p>}
        {error && <p className="inline-status error" role="alert">{error}</p>}
        {!ready && <p className="form-notice">请连接钱包并切换到 Base Sepolia</p>}
      </div>

      {plan ? (
        <section className="review-card" aria-labelledby="launch-review-title">
          <p className="panel-label">REVIEW / SIGN</p>
          <h2 id="launch-review-title">发行确认</h2>
          <dl className="review-list">
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
          <p className="review-copy">{plan.review}</p>
          <button className="primary-action" disabled={phase !== 'review'} onClick={broadcast}>确认并签名</button>
        </section>
      ) : (
        <aside className="protocol-brief">
          <p className="panel-label">IMMUTABLE PARAMETERS</p>
          <h2>发行前须知</h2>
          <dl className="brief-list">
            <div><dt>供应量</dt><dd>1B / FIXED</dd></div>
            <div><dt>流动性</dt><dd>永久锁定</dd></div>
            <div><dt>正常费率</dt><dd>1.5%</dd></div>
            <div><dt>保护期</dt><dd>16 秒</dd></div>
            <div><dt>Quote</dt><dd>{quoteLabel(config.quote)}</dd></div>
          </dl>
          <p>生成计划不会立即广播交易。你仍需在下一步审阅并通过钱包确认。</p>
        </aside>
      )}
    </section>
  )
}

function message(cause: unknown) {
  return cause instanceof Error ? cause.message : 'Token 发行失败'
}

function quoteLabel(quote: string) {
  return quote.toLowerCase() === '0x0000000000000000000000000000000000000000' ? 'NATIVE ETH' : quote
}
