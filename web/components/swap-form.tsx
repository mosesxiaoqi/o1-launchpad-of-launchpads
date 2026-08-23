'use client'

import { useState, type FormEvent } from 'react'
import { erc20Abi } from 'viem'
import { useAccount, useBalance, usePublicClient, useReadContract, useSendTransaction } from 'wagmi'

import { launchpadApi, type PreparedTransaction, type Token } from '@/lib/api'
import { chain } from '@/lib/wagmi'

type Quote = Awaited<ReturnType<typeof launchpadApi.quoteSwap>>

export function SwapForm({ token }: { token: Token }) {
  const { address, chainId, isConnected } = useAccount()
  const { sendTransactionAsync } = useSendTransaction()
  const publicClient = usePublicClient()
  const [buy, setBuy] = useState(true)
  const [amount, setAmount] = useState('')
  const [quote, setQuote] = useState<Quote | null>(null)
  const [plan, setPlan] = useState<PreparedTransaction | null>(null)
  const [pending, setPending] = useState(false)
  const [error, setError] = useState('')
  const [status, setStatus] = useState('')
  const inputCurrency = buy ? token.quote : token.token
  const nativeInput = inputCurrency.toLowerCase() === '0x0000000000000000000000000000000000000000'
  const nativeBalance = useBalance({ address, query: { enabled: nativeInput } })
  const tokenBalance = useReadContract({
    abi: erc20Abi,
    address: inputCurrency,
    functionName: 'balanceOf',
    args: address ? [address] : undefined,
    query: { enabled: Boolean(address && !nativeInput) },
  })
  const balance = nativeInput ? nativeBalance.data?.value : tokenBalance.data
  const ready = Boolean(address && isConnected && chainId === chain.id)

  function resetDirection(nextBuy: boolean) {
    setBuy(nextBuy)
    setQuote(null)
    setPlan(null)
    setError('')
  }

  async function getQuote(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setError('')
    setQuote(null)
    setPlan(null)
    if (!/^\d+$/.test(amount) || BigInt(amount) === 0n) {
      setError('金额必须大于 0')
      return
    }
    if (balance !== undefined && BigInt(amount) > balance) {
      setError('余额不足')
      return
    }
    if (!address || !ready) return
    setPending(true)
    try {
      setQuote(await launchpadApi.quoteSwap({ chain_id: chain.id, token: token.token, amount, buy, wallet: address }))
    } catch (cause) {
      setError(message(cause))
    } finally {
      setPending(false)
    }
  }

  async function prepare() {
    if (!quote || !address) return
    setPending(true)
    setError('')
    try {
      setPlan(await launchpadApi.prepareSwap({ chain_id: chain.id, quote_id: quote.quote_id, slippage_bps: 100, wallet: address }))
    } catch (cause) {
      setPlan(null)
      setError(message(cause))
    } finally {
      setPending(false)
    }
  }

  async function broadcast() {
    if (!plan || !publicClient) return
    setPending(true)
    setError('')
    setStatus('等待钱包确认')
    try {
      const hash = await sendTransactionAsync({ to: plan.to, data: plan.data, value: BigInt(plan.value), chainId: chain.id })
      setStatus('等待 Swap 交易确认')
      const receipt = await publicClient.waitForTransactionReceipt({ hash })
      if (receipt.status !== 'success') throw new Error('Swap 交易已回滚')
      setStatus('Swap 成功')
    } catch (cause) {
      setError(message(cause))
      setStatus('')
    } finally {
      setPending(false)
    }
  }

  return (
    <section aria-labelledby="swap-title">
      <h2 id="swap-title">Swap（exact-input）</h2>
      <fieldset>
        <legend>方向</legend>
        <label><input type="radio" name="direction" checked={buy} onChange={() => resetDirection(true)} />买入</label>
        <label><input type="radio" name="direction" checked={!buy} onChange={() => resetDirection(false)} />卖出</label>
      </fieldset>
      <form onSubmit={getQuote}>
        <label htmlFor="swap-amount">输入金额（raw units）</label>
        <input id="swap-amount" inputMode="numeric" value={amount} onChange={(event) => setAmount(event.target.value)} />
        <button disabled={!ready || pending}>获取 Quote</button>
      </form>
      {quote && (
        <section aria-labelledby="quote-title">
          <h3 id="quote-title">报价</h3>
          <p>预计输出：{quote.amount_out}</p>
          <p>Fee：{quote.fee}</p>
          <p>Price Impact：已包含在 Quoter 输出中</p>
          <p>Quote 有效期：{quote.expires_at}</p>
          {!plan && <button type="button" disabled={pending} onClick={prepare}>准备交易</button>}
        </section>
      )}
      {plan && (
        <section aria-labelledby="swap-review-title">
          <h3 id="swap-review-title">签名前确认</h3>
          <p>Slippage：1%</p>
          <p>anti-snipe：16 秒</p>
          <p>Deadline：{plan.deadline}</p>
          <p>{plan.review}</p>
          <button type="button" disabled={pending} onClick={broadcast}>确认并签名</button>
        </section>
      )}
      {!ready && <p>请连接钱包并切换到 Base Sepolia</p>}
      {status && <p role="status">{status}</p>}
      {error && <p role="alert">{error}</p>}
    </section>
  )
}

function message(cause: unknown) {
  return cause instanceof Error ? cause.message : 'Swap 操作失败'
}
