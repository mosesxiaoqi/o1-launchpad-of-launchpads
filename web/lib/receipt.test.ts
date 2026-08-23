import { encodeAbiParameters, encodeEventTopics, type Address, type Hex } from 'viem'
import { describe, expect, it } from 'vitest'

import { multiTenantFactoryAbi } from './contracts'
import { parseLaunchedReceipt } from './receipt'

const factory = '0x1111111111111111111111111111111111111111' as Address
const token = '0x2222222222222222222222222222222222222222' as Address
const creator = '0x3333333333333333333333333333333333333333' as Address
const quote = '0x4444444444444444444444444444444444444444' as Address
const poolId = `0x${'55'.repeat(32)}` as Hex
const launchpadId = `0x${'66'.repeat(32)}` as Hex

function launchedLog(address: Address) {
  return {
    address,
    topics: encodeEventTopics({
      abi: multiTenantFactoryAbi,
      eventName: 'Launched',
      args: { token, poolId, launchpadId },
    }) as [Hex, ...Hex[]],
    data: encodeAbiParameters(
      [
        { type: 'address' },
        { type: 'address' },
        { type: 'uint256' },
        { type: 'int24' },
      ],
      [creator, quote, 1_000_000n, 60],
    ),
  }
}

describe('parseLaunchedReceipt', () => {
  it('parses only the configured factory event', () => {
    expect(parseLaunchedReceipt({ status: 'success', logs: [launchedLog(factory)] }, factory)).toEqual({
      token,
      poolId,
      launchpadId,
      creator,
      quote,
      supply: 1_000_000n,
      tickSpacing: 60,
    })
  })

  it('rejects reverted receipts and lookalike events from another contract', () => {
    expect(() => parseLaunchedReceipt({ status: 'reverted', logs: [] }, factory)).toThrow()
    expect(() =>
      parseLaunchedReceipt(
        { status: 'success', logs: [launchedLog('0x9999999999999999999999999999999999999999')] },
        factory,
      ),
    ).toThrow()
  })
})
