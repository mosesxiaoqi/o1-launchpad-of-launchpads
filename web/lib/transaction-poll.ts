type TransactionState = { status: string }

export async function pollTransaction<T extends TransactionState>(
  read: () => Promise<T>,
  options: {
    attempts?: number
    initialDelayMs?: number
    maxDelayMs?: number
    sleep?: (delayMs: number) => Promise<void>
  } = {},
): Promise<T> {
  const attempts = options.attempts ?? 8
  const initialDelay = options.initialDelayMs ?? 250
  const maxDelay = options.maxDelayMs ?? 4_000
  const sleep = options.sleep ?? ((delay) => new Promise((resolve) => setTimeout(resolve, delay)))

  for (let attempt = 0; attempt < attempts; attempt += 1) {
    try {
      const transaction = await read()
      if (transaction.status === 'confirmed' || transaction.status === 'reverted') return transaction
    } catch {
      // Indexing can lag the receipt by a few blocks.
    }
    if (attempt + 1 < attempts) await sleep(Math.min(initialDelay * 2 ** attempt, maxDelay))
  }
  throw new Error('transaction indexing timed out')
}
