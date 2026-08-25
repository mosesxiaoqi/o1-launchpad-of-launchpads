export type TransactionPhase = 'idle' | 'auth' | 'wallet' | 'confirming' | 'saving' | 'success'

const labels: Record<Exclude<TransactionPhase, 'idle'>, string> = {
  auth: '等待钱包签名',
  wallet: '等待钱包确认交易',
  confirming: '等待 Base Sepolia 确认',
  saving: '正在验证交易并保存品牌',
  success: 'Launchpad 创建成功',
}

export function TransactionStatus({ phase, error }: { phase: TransactionPhase; error?: string }) {
  if (error) {
    return <p className="inline-status error" role="alert" aria-live="assertive">{error}</p>
  }
  if (phase === 'idle') return null
  return <p className="inline-status" role="status" aria-live="polite"><span />{labels[phase]}</p>
}
