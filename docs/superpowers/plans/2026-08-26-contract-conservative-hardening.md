# EVM 合约保守安全加固 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 在不改变公开 ABI、CREATE2、Hook 权限、费率和资金语义的前提下，修复已确认的有符号边界问题，补齐协议级测试，并仅保留可测量且可审计的高频路径 gas 优化。

**Architecture:** 保留 Registry、Token、Factory、Hook、Escrow 五合约职责。生产代码只修改 `MultiTenantLaunchpadFactory` 的构造参数边界，以及 `LaunchHookV2` 的内部数值转换、费率计算和 Referrer 校验复用；Registry、Token、Escrow 默认只增加测试。所有变更先由回归测试刻画，再用最小实现修复。

**Tech Stack:** Solidity 0.8.26、Foundry、Cancun EVM、optimizer runs 200、via-IR、OpenZeppelin、Uniswap v4 Core/Periphery。

**Spec:** `docs/superpowers/specs/2026-08-26-contract-conservative-hardening-design.md`

## Global Constraints

- 不新增 external/public 函数，不改变 selector、事件、custom error、参数、返回值或 payable 属性。
- 不改变 `LaunchParams`、`PoolConfig` 字段顺序和 public getter 返回结构。
- 不改变 caller-bound salt、LaunchToken creation code、CREATE2 部署者和地址推导。
- 不改变 Hook permissions、PoolKey/PoolId、币序、tick 对齐规则、费率和分账语义。
- 不引入代理、新依赖、assembly、批量记账接口或新的生产 Library。
- 测试失败时先定位根因；不得为了让测试通过而放宽安全断言。
- 每个生产修改必须有对应测试；gas 优化若无同输入改进则回退。

## File Map

- Modify: `contracts/src/MultiTenantLaunchpadFactory.sol`
- Modify: `contracts/src/LaunchHookV2.sol`
- Modify: `contracts/test/MultiTenantLaunchpadFactory.t.sol`
- Modify: `contracts/test/LaunchHookV2Fees.t.sol`
- Modify: `contracts/test/LaunchHookV2Integration.t.sol`
- Modify: `contracts/test/FeeEscrow.t.sol`
- Verify only: `contracts/src/LaunchpadRegistry.sol`
- Verify only: `contracts/src/LaunchToken.sol`
- Verify only: `contracts/src/FeeEscrow.sol`
- Verify only: `contracts/script/DeployBaseSepolia.s.sol`
- Verify consumers: `server/common/chain/bindings/`, `server/internal/indexer/`, `web/lib/contracts.ts`, `web/lib/receipt.ts`

---

### Task 1: Freeze compatibility and performance baselines

**Files:**
- Inspect: `contracts/out/`
- Inspect: `server/common/chain/bindings/`
- Inspect: `web/lib/contracts.ts`

- [ ] Run the clean baseline:

```bash
cd contracts
forge fmt --check
forge build
forge test
forge test --gas-report
```

Expected: 8 suites、45 tests、0 failures；记录 `LaunchHookV2` deployed size 和 Factory `launch` gas。

- [ ] Save compatibility fingerprints outside the repository before editing:

```bash
cd contracts
forge inspect LaunchHookV2 methodIdentifiers > /tmp/o1-LaunchHookV2-methods.before
forge inspect MultiTenantLaunchpadFactory methodIdentifiers > /tmp/o1-Factory-methods.before
forge inspect FeeEscrow methodIdentifiers > /tmp/o1-FeeEscrow-methods.before
forge inspect LaunchHookV2 abi > /tmp/o1-LaunchHookV2-abi.before
forge inspect MultiTenantLaunchpadFactory abi > /tmp/o1-Factory-abi.before
forge inspect FeeEscrow abi > /tmp/o1-FeeEscrow-abi.before
```

- [ ] Confirm working tree is clean before implementation:

```bash
git status --short
```

### Task 2: Reject unusable Factory tick configuration

**Files:**
- Modify: `contracts/test/MultiTenantLaunchpadFactory.t.sol`
- Modify: `contracts/src/MultiTenantLaunchpadFactory.sol`

- [ ] Add a failing constructor-boundary test:

```solidity
function testRejectsMinimumStartTickFrame() public {
    vm.expectRevert(MultiTenantLaunchpadFactory.InvalidConfig.selector);
    new MultiTenantLaunchpadFactory(
        registry, manager, hook, Currency.unwrap(currency1), laasTreasury, type(int24).min, 60
    );
}
```

- [ ] Run the focused test and verify it fails because construction currently succeeds:

```bash
cd contracts
forge test --match-contract MultiTenantLaunchpadFactoryTest --match-test testRejectsMinimumStartTickFrame -vv
```

- [ ] Extend the existing constructor guard with the smallest compatible rejection:

```solidity
|| startTickToken0Frame_ == type(int24).min
```

Use the existing `InvalidConfig` error. Do not alter `launch()` or tick calculations.

- [ ] Run the focused suite:

```bash
cd contracts
forge test --match-contract MultiTenantLaunchpadFactoryTest
```

- [ ] Commit the isolated fix:

```bash
git add contracts/src/MultiTenantLaunchpadFactory.sol contracts/test/MultiTenantLaunchpadFactory.t.sol
git commit -m "fix(contracts): reject overflowing tick frame"
```

### Task 3: Make signed delta magnitude total and explicit

**Files:**
- Modify: `contracts/test/LaunchHookV2Fees.t.sol`
- Modify: `contracts/src/LaunchHookV2.sol`

- [ ] Add a test-only Harness method that calls an internal `int128` magnitude helper:

```solidity
function magnitudeForTest(int128 amount) external pure returns (uint256) {
    return _magnitude(amount);
}
```

Add boundary tests:

```solidity
function testMagnitudeHandlesSignedInt128Boundaries() public view {
    assertEq(hook.magnitudeForTest(type(int128).min), uint256(1) << 127);
    assertEq(hook.magnitudeForTest(type(int128).max), uint256(uint128(type(int128).max)));
    assertEq(hook.magnitudeForTest(-1), 1);
    assertEq(hook.magnitudeForTest(0), 0);
}
```

- [ ] Run the focused test and confirm compilation fails because `_magnitude` does not exist:

```bash
cd contracts
forge test --match-test testMagnitudeHandlesSignedInt128Boundaries
```

- [ ] Add this internal pure helper without assembly:

```solidity
function _magnitude(int128 amount) internal pure returns (uint256) {
    if (amount < 0) return uint256(uint128(-(amount + 1))) + 1;
    return uint256(uint128(amount));
}
```

Use it in both `_unspecified` and `unlockCallback`. Keep `_specifiedMagnitude(int256)` unchanged because it already handles `int256.min` safely.

- [ ] Run the fee, pool and integration suites:

```bash
cd contracts
forge test --match-contract LaunchHookV2FeesTest
forge test --match-contract LaunchHookV2PoolTest
forge test --match-contract LaunchHookV2IntegrationTest
```

- [ ] Commit the isolated fix:

```bash
git add contracts/src/LaunchHookV2.sol contracts/test/LaunchHookV2Fees.t.sol
git commit -m "fix(contracts): handle signed delta boundaries"
```

### Task 4: Unify Hook fee calculation and Referrer validity

**Files:**
- Modify: `contracts/test/LaunchHookV2Fees.t.sol`
- Modify: `contracts/src/LaunchHookV2.sol`

- [ ] Add characterization coverage for all invalid Referrer identities and a valid Referrer at anti-snipe boundaries. Assert preview fee conservation and exact current fee at timestamps `launchTime`, `+8`, `+16`, and `+17`.

- [ ] Run characterization tests before refactoring; they must pass:

```bash
cd contracts
forge test --match-contract LaunchHookV2FeesTest
```

- [ ] Introduce private/internal helpers that consume the already-loaded `PoolConfig memory`:

```solidity
function _currentTotalFeeBps(PoolConfig memory config, uint256 timestamp) private pure returns (uint256)
function _isValidReferrer(PoolConfig memory config, address trader, address referrer) private pure returns (bool)
```

Keep public `currentTotalFeeBps(bytes32,uint256)` as the ABI-compatible storage-loading wrapper. Make `previewFeeSplit`, `_beforeSwap`, `_afterSwap`, and `_distribute` reuse the helpers, eliminating repeated `poolConfig[poolId]` reads and duplicated Referrer conditions.

- [ ] Run focused regression and gas report:

```bash
cd contracts
forge test --match-contract LaunchHookV2FeesTest
forge test --match-contract LaunchHookV2PoolTest
forge test --match-contract LaunchHookV2IntegrationTest
forge test --gas-report
```

Expected: behavior unchanged; retain the refactor only if same-input Hook paths do not regress materially and deployed bytecode remains below the baseline or any increase is justified by the signed-boundary fix.

- [ ] Commit the behavior-preserving refactor:

```bash
git add contracts/src/LaunchHookV2.sol contracts/test/LaunchHookV2Fees.t.sol
git commit -m "refactor(contracts): reuse loaded hook fee config"
```

### Task 5: Cover hookData parsing, exact-output, ordering and partial fills

**Files:**
- Modify: `contracts/test/LaunchHookV2Fees.t.sol`
- Modify: `contracts/test/LaunchHookV2Integration.t.sol`
- Modify: `contracts/test/MultiTenantLaunchpadFactory.t.sol`

- [ ] Expose `_parseHookData` through the test Harness by changing only its production visibility from `private` to `internal`, then add table-driven tests for lengths 0、31、32、63、64、65. Assert:
  - below 32 bytes returns zero Referrer；
  - 32–63 bytes parses Referrer and zero comment；
  - 64+ bytes parses only the first Referrer/comment words；
  - trailing bytes do not change parsed values。

- [ ] Add integration coverage for exact-output after `launchTime + 16`, exact-input buy/sell, and both quote specified/unspecified paths. Assert `totalOwed == PoolManager.balanceOf(Escrow, quote)` after each trade.

- [ ] Add/extend a callback-level partial-fill test that expects `PartialFillUnsupported` when actual specified delta differs from `amountSpecified + specifiedFee`.

- [ ] Add Factory characterization asserting predicted CREATE2 address equals emitted/deployed token and derive salts stay caller-bound. Exercise both token/quote orderings with deterministic quote fixtures when feasible; if the existing `Deployers` currencies cannot deterministically cover both orderings, document the limitation and keep the production ordering code unchanged.

- [ ] Run focused suites:

```bash
cd contracts
forge test --match-contract LaunchHookV2FeesTest
forge test --match-contract LaunchHookV2PoolTest
forge test --match-contract LaunchHookV2IntegrationTest
forge test --match-contract MultiTenantLaunchpadFactoryTest
```

- [ ] Commit tests and the visibility-only parser change:

```bash
git add contracts/src/LaunchHookV2.sol contracts/test/LaunchHookV2Fees.t.sol contracts/test/LaunchHookV2Integration.t.sol contracts/test/MultiTenantLaunchpadFactory.t.sol contracts/test/LaunchHookV2Pool.t.sol
git commit -m "test(contracts): cover swap and hook data boundaries"
```

### Task 6: Strengthen Escrow accounting invariants

**Files:**
- Modify: `contracts/test/FeeEscrow.t.sol`
- Verify only: `contracts/src/FeeEscrow.sol`

- [ ] Add a stateful sequence test or Foundry invariant Handler that alternates funded credits and claims across at least two recipients. After every operation assert:

```solidity
escrow.totalOwed(address(currency))
    == escrow.owed(firstRecipient, address(currency)) + escrow.owed(secondRecipient, address(currency));
assertLe(escrow.totalOwed(address(currency)), manager.balanceOf(address(escrow), uint160(address(currency))));
```

- [ ] Include failed over-credit, double-claim, redirected claim, native claim and reentrant token cases in the focused suite. Production `FeeEscrow.sol` remains unchanged unless a new failing test proves a real violation.

- [ ] Run the suite:

```bash
cd contracts
forge test --match-contract FeeEscrowTest -vv
```

- [ ] Commit the tests:

```bash
git add contracts/test/FeeEscrow.t.sol
git commit -m "test(contracts): strengthen escrow accounting invariants"
```

### Task 7: Verify ABI, CREATE2 surface, gas and full repository consumers

**Files:**
- Verify: `contracts/src/*.sol`
- Verify: `contracts/test/*.t.sol`
- Verify: `contracts/script/DeployBaseSepolia.s.sol`
- Verify: `server/common/chain/bindings/`
- Verify: `server/internal/indexer/`
- Verify: `web/lib/contracts.ts`
- Verify: `web/lib/receipt.ts`

- [ ] Format and run full verification from a fresh build:

```bash
cd contracts
forge fmt
forge clean
forge build
forge test
forge test --gas-report
```

- [ ] Compare selectors and full ABI to the frozen baseline:

```bash
cd contracts
forge inspect LaunchHookV2 methodIdentifiers > /tmp/o1-LaunchHookV2-methods.after
forge inspect MultiTenantLaunchpadFactory methodIdentifiers > /tmp/o1-Factory-methods.after
forge inspect FeeEscrow methodIdentifiers > /tmp/o1-FeeEscrow-methods.after
forge inspect LaunchHookV2 abi > /tmp/o1-LaunchHookV2-abi.after
forge inspect MultiTenantLaunchpadFactory abi > /tmp/o1-Factory-abi.after
forge inspect FeeEscrow abi > /tmp/o1-FeeEscrow-abi.after
diff -u /tmp/o1-LaunchHookV2-methods.before /tmp/o1-LaunchHookV2-methods.after
diff -u /tmp/o1-Factory-methods.before /tmp/o1-Factory-methods.after
diff -u /tmp/o1-FeeEscrow-methods.before /tmp/o1-FeeEscrow-methods.after
diff -u /tmp/o1-LaunchHookV2-abi.before /tmp/o1-LaunchHookV2-abi.after
diff -u /tmp/o1-Factory-abi.before /tmp/o1-Factory-abi.after
diff -u /tmp/o1-FeeEscrow-abi.before /tmp/o1-FeeEscrow-abi.after
```

Expected: all diffs empty. If not empty, stop and revert the accidental compatibility change before proceeding.

- [ ] Confirm the integration consumers still reference unchanged events, fields and addresses:

```bash
rg "Launched|Trade|Credited|Claimed|LaunchHookV2|MultiTenantLaunchpadFactory" server web contracts/script
```

Do not regenerate bindings if ABI is identical. Do not change Base Sepolia deployment addresses.

- [ ] Review production diff for forbidden changes and accidental placeholders:

```bash
git diff --check
rg -n "TODO|TBD|FIXME|placeholder" contracts/src contracts/test
git diff -- contracts/src contracts/test
```

- [ ] Commit any final formatting/test-only adjustments, then verify clean status:

```bash
git status --short
```

## Completion Criteria

- Factory explicitly rejects `type(int24).min` with existing `InvalidConfig`。
- Hook signed magnitude handles `int128.min` without arithmetic panic。
- Preview and actual distribution share one Referrer validity rule。
- Swap direction、exact-output、partial fill、hookData and Escrow invariants have tests。
- All Foundry tests pass，format/build/gas report complete。
- Hook、Factory、Escrow selectors and ABI are byte-for-byte unchanged。
- LaunchToken creation code、salt derivation、Hook permissions、PoolKey/PoolId and deployment addresses are unchanged。
- No confirmed live security issue requires redeployment；若验证发现此类问题，单独报告，不自动部署。
