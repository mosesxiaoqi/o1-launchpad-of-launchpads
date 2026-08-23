import {
  decodeEventLog,
  isAddressEqual,
  type Address,
  type Hex,
} from 'viem'

import { multiTenantFactoryAbi } from './contracts'

type ReceiptLog = {
  address: Address
  data: Hex
  topics: [] | [Hex, ...Hex[]]
}

type Receipt = {
  status: 'success' | 'reverted'
  logs: readonly ReceiptLog[]
}

export type LaunchedResult = {
  token: Address
  poolId: Hex
  launchpadId: Hex
  creator: Address
  quote: Address
  supply: bigint
  tickSpacing: number
}

export function parseLaunchedReceipt(receipt: Receipt, factory: Address): LaunchedResult {
  if (receipt.status !== 'success') {
    throw new Error('launch transaction reverted')
  }
  const matches: LaunchedResult[] = []
  for (const log of receipt.logs) {
    if (!isAddressEqual(log.address, factory)) continue
    try {
      const decoded = decodeEventLog({
        abi: multiTenantFactoryAbi,
        eventName: 'Launched',
        data: log.data,
        topics: log.topics,
        strict: true,
      })
      matches.push(decoded.args)
    } catch {
      // A trusted contract can emit other events in the same transaction.
    }
  }
  if (matches.length !== 1) {
    throw new Error('expected exactly one trusted Launched event')
  }
  return matches[0]
}
