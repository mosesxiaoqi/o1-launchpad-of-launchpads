# EVM 合约保守安全加固设计

日期：2026-08-26

## 目标

按照 `contracts/AGENTS.md` 对当前 Solidity/Foundry 合约进行审计驱动的保守重写：保持外部契约与业务语义，优先补齐安全边界和测试，只保留有基准证明且不降低可审计性的 gas 优化。

本设计不以现有代码风格作为质量上限。现有代码用于识别技术栈、调用关系和风险面；实施结果必须符合 `contracts/AGENTS.md` 的安全、兼容和验证要求。

## 已确认策略

采用“审计驱动的保守重写”：

- 保持现有函数 selector、事件、错误参数、业务费率和 CREATE2 计算规则。
- 先用失败测试证明安全或边界问题，再修改生产代码。
- 以 `LaunchHookV2` 为主要重写对象，其他合约按证据最小修改。
- gas 优化只针对高频路径，并以相同编译配置和测试输入比较。
- 不因为纯结构整理或 gas 优化要求替换当前 Base Sepolia 地址。
- 只有确认存在影响资金、权限或协议不变量的安全问题时，才单独提出重新部署建议。

## 当前基线

技术栈：

- Solidity 0.8.26；
- Foundry；
- Cancun EVM；
- optimizer runs 200；
- via-IR；
- OpenZeppelin；
- Uniswap v4 Core/Periphery；
- Go 后端通过 abigen bindings、`Launched` 事件和部署清单消费合约。

设计前执行以下命令，退出码均为 0：

```bash
forge fmt --check
forge build
forge test
forge test --gas-report
```

基线测试为 8 个 suite、45 项测试通过、0 失败、0 跳过。

需要重点跟踪的基线指标：

- `LaunchHookV2` deployed size：11,495 bytes；
- `MultiTenantLaunchpadFactory.launch` 当前测试最大约 941,545 gas；
- swap、费用预览、Pool 注册、流动性播种及 Escrow claim 的同输入 gas；
- Factory、Hook、Escrow 和 Token 的部署成本与 bytecode 大小。

gas report 会受到测试路径和 instrumentation 影响，因此只比较同一环境、同一提交范围、同一测试输入的前后结果。

## 范围

生产合约：

- `contracts/src/LaunchHookV2.sol`；
- `contracts/src/MultiTenantLaunchpadFactory.sol`；
- `contracts/src/FeeEscrow.sol`；
- `contracts/src/LaunchpadRegistry.sol`；
- `contracts/src/LaunchToken.sol`；
- `contracts/src/base/BaseHook.sol`，仅在必要时读取或做兼容修正。

测试与验证：

- `contracts/test/*.t.sol`；
- 必要的测试 Harness 和 invariant Handler；
- `contracts/script/DeployBaseSepolia.s.sol` 与 wiring 测试；
- ABI 未变化时不重新生成 Go bindings；若编译产物显示 ABI 差异，必须停止并查明原因。

不在范围内：

- 代理或升级框架；
- 新的生产外部 Library；
- 新依赖；
- 后端、前端或索引器功能改造；
- Base Sepolia 广播、迁移或地址替换；
- 与安全、可读性或可测量 gas 无关的重构。

## 架构与职责

保留现有五个生产合约及职责。

### LaunchpadRegistry

继续负责 launchpad ID、owner、treasury、active 状态和创建时间。它不是高频交易路径，默认不做微型 gas 优化。只有测试证明存在安全或一致性问题时才修改生产代码。

### LaunchToken

继续作为构造时一次性铸造、无后续 mint 权限的 ERC-20。默认不修改其 creation code，避免没有安全收益地改变未来 Token CREATE2 init code hash 和预测地址。

### MultiTenantLaunchpadFactory

继续负责：

- 发布参数与配置版本校验；
- caller-bound salt；
- Token CREATE2 预测与部署；
- Token/quote 排序；
- tick 与初始价格计算；
- Pool 初始化；
- Hook Pool 配置注册；
- 永久单边流动性播种；
- Token、Pool 与 launchpad 映射及 `Launched` 事件。

允许在不改变 `launch()` selector、`LaunchParams`、事件和执行语义的前提下，将校验、地址计算和 Pool 几何计算整理为边界清晰的 internal/private 步骤。

### LaunchHookV2

继续作为以下逻辑的唯一中心：

- Uniswap v4 Hook permissions；
- Factory-only wiring、Pool 注册与播种；
- 外部流动性锁定；
- anti-snipe；
- exact-input/exact-output 费用处理；
- quote 资产识别；
- Creator/Platform/Referrer/LaaS/surcharge 分配；
- Trade 事件；
- 向 FeeEscrow 记账。

这是本次主要内部重写对象，但公开 ABI、事件、`PoolConfig` 字段顺序和 getter 返回结构保持不变。

### FeeEscrow

继续使用 PoolManager ERC-6909 claim 作为可赎回资产，并保持：

- `credit(address,address,uint256)`；
- `claim(address,address)`；
- `claimTo(address,address)`；
- `unlockCallback(bytes)`；
- `owed` 与 `totalOwed`；
- `Credited` 与 `Claimed` 事件。

不增加 `creditBatch`。虽然批量接口可能减少 Hook 的外部调用，但会扩大 ABI、状态变更批次和审计面，不符合本次保守策略。

## 兼容性约束

实施不得改变：

- external/public 函数 selector；
- 现有参数、返回值及 payable 属性；
- 事件名称、参数顺序、类型和 `indexed` 标记；
- custom error 的名称、参数和 selector；
- `LaunchParams` 与 `PoolConfig` 的字段顺序；
- Hook permissions；
- `Launched`、`Trade`、`Credited`、`Claimed` 的语义；
- salt 派生：`keccak256(abi.encode(msg.sender, params.salt))`；
- LaunchToken 构造参数及 CREATE2 预测规则；
- PoolKey、PoolId、币种排序和配置；
- 正常费率、anti-snipe 参数和分账规则。

任何意外 ABI 差异都视为设计违例，必须停止实施并重新获得批准。

内部函数可以拆分、合并或调整可见性。为了通过 Harness 测试纯函数边界，可以将 private 调整为 internal，但不得新增生产 external/public 入口。

## 必须保持的协议不变量

- Launchpad ID 绑定 chain ID 与规范化 slug。
- 发布只接受已存在、active、owner/treasury 有效的 launchpad。
- 配置版本与 deadline 必须 fail closed。
- 用户 salt 与 caller 绑定，同一派生 salt 不能重复使用。
- CREATE2 预测地址必须等于实际 Token 地址。
- Token/quote 排序、PoolKey 和 PoolId 必须确定性一致。
- Hook 地址权限位必须与 `getHookPermissions()` 一致。
- Hook 只能 wiring 一次，Pool 只能注册和播种一次。
- 外部账户不能增加或移除永久流动性。
- LaunchToken 固定供应量为 1,000,000,000 ether，全部进入发布流程指定的 Hook/Pool 路径，不留 mint 权限。
- 正常费率为 150 bps：protocol 100 bps，LaaS 50 bps。
- protocol fee 的 Creator/Platform/Referrer 分配为 50%/30%/20%；无效 Referrer 份额归 Platform。
- anti-snipe 起始费率为 9,900 bps，窗口 16 秒，按现行规则衰减；窗口内拒绝 exact-output。
- 每次收费满足费用守恒。
- Escrow 债权总额不超过 PoolManager 中可赎回余额。
- claim 清账先于外部交互，任何失败整体回滚，同一债权不能领取两次。
- 只有指定 PoolManager 可以调用 `unlockCallback`。

## 安全加固设计

### 数值与转换

为所有有符号 magnitude 计算增加统一、安全的转换路径，覆盖：

- `int128.min`；
- `int256.min`；
- swap delta 正负方向；
- liquidity 的 `uint128`/`int256` 转换；
- tick 的 `int24` 取反、除法、对齐和上下界。

先通过 Harness 或真实调用路径编写失败测试。只有测试证明现有代码会在合法或应明确拒绝的输入上违反不变量，才修改生产函数。

Factory 构造配置必须拒绝会导致后续 launch 取反溢出或永久不可用的 tick 参数。使用现有 `InvalidConfig`，不新增错误 selector。

### Factory 与 CREATE2

测试两种 Token/quote 地址顺序，验证：

- 预测地址等于部署地址；
- Pool 币序正确；
- tickLower、tickUpper、initTick 对齐且在 TickMath 边界内；
- liquidity 非零且转换安全；
- 不同 caller 的相同用户 salt 得到不同派生 salt；
- 相同派生 salt 重放失败；
- 错误 quote、stale config、过期 deadline 和无效 launchpad 失败。

### Hook 生命周期与 callback

保持并补强：

- 只有 deployer 可以 wiring；
- wiring 只能成功一次；
- FeeEscrow 必须绑定当前 Hook 和同一 PoolManager；
- 只有 Factory 可以初始化相关 Pool、注册配置和播种；
- 只有 PoolManager 可以进入 callback；
- 未 wiring、未知 Pool、重复注册、重复播种和外部流动性修改均 fail closed；
- callback 的 Pool、Currency、方向、delta 符号和 single-sided 条件符合预期。

不引入新的跨交易 callback 状态锁，除非测试证明当前 PoolManager 信任边界不足。PoolManager 是 immutable 的协议依赖，不能为了防御假设中的恶意 PoolManager 扩大状态和 gas 成本。

### Swap 与费用

必须覆盖以下组合：

- exact-input buy；
- exact-input sell；
- anti-snipe 期间 exact-output；
- anti-snipe 结束后的 exact-output；
- quote 是 specified currency；
- quote 是 unspecified currency；
- Token 为 currency0；
- Token 为 currency1；
- 完整成交；
- partial fill；
- 零费用舍入与极小金额；
- 边界金额；
- 有效/无效 Referrer；
- hookData 长度为 0、31、32、63、64 和大于 64。

所有预览与实际分账使用相同的 Referrer 有效性规则。

费用守恒：

```text
creator + protocol + referrer + laas + surcharge == totalFee
```

实际记账守恒：

```text
FeeEscrow.totalOwed(currency)
    == 所有接收者 owed 之和
    <= PoolManager.balanceOf(FeeEscrow, currency)
```

### Escrow

测试并保持：

- 非 Hook 不能 credit；
- 零 recipient 和零 amount 失败；
- aggregate credit 不能超过实际可赎回余额；
- `claim(recipient,currency)` 只向原 recipient 支付；
- `claimTo(currency,to)` 只移动调用者自己的债权；
- claim 使用 CEI 和 `nonReentrant`；
- 原生币和 ERC-20 均可领取；
- 重入、拒绝支付或 PoolManager 失败时状态完全回滚；
- 重复领取失败。

增加 stateful invariant Handler 时，只建模真实可达操作，不通过测试作弊函数直接写生产存储。

## 内部重写设计

### LaunchHookV2

- 增加接收已加载 `PoolConfig memory` 的内部费率函数，公开 `currentTotalFeeBps` 保持原签名并调用它。
- `_beforeSwap` 和 `_afterSwap` 每条路径只加载一次 PoolConfig，只计算一次 PoolId、费率、方向和 exact-input/output。
- Referrer 有效性收敛为单一 internal/private pure helper，预览和实际分账复用。
- 有符号 delta 转 magnitude 收敛为安全 helper，显式覆盖最小负数。
- 缓存单次调用使用的 Escrow 地址、Currency 地址和配置字段。
- 仅在不改变 Hook delta 和 revert 行为的前提下消除重复分支。
- 循环长度只读取一次；`unchecked` 递增只有在边界证明、测试和 gas 基准同时成立时采用。
- 不使用 assembly，不重排存储，不增加外部函数。

### MultiTenantLaunchpadFactory

- 在构造阶段补齐不可用 tick 配置的边界校验。
- 在不改变执行顺序和原子性的前提下，将参数校验、CREATE2 预测、Pool 几何计算整理为小型 internal/private 步骤。
- `usedLaunchSalts` 仍在外部交互前写入，后续失败依赖 EVM 回滚保证原子性。
- Token 部署、Pool 初始化、Hook 注册、播种、映射写入和事件仍在同一交易完成。
- 不拆出新的生产外部 Library。

### FeeEscrow

默认不做结构重写。只有失败测试证明问题时修复；允许使用不改变 ABI、事件或校验顺序的局部缓存。

### Registry 与 Token

默认不修改生产代码。新增测试如果证明安全问题，再回到设计审批，不在实施阶段临时扩大范围。

## Gas 优化原则与验收

优先级固定为：正确性、安全性、兼容性、可审计性、gas。

重点优化高频路径：

- `_beforeSwap`；
- `_afterSwap`；
- fee/referrer 计算；
- Escrow credit/claim；
- 高频 view 费率预览。

Factory launch、Registry 和部署脚本是低频路径，只有明显浪费或安全修复顺带改善时才优化。

测量命令：

```bash
forge test --gas-report
forge snapshot
```

只保留满足以下全部条件的优化：

- 在同一测试输入上可重复观察到改善；
- 不改变公开行为、revert 优先级或事件；
- 不降低检查、SafeCast、重入保护和测试覆盖；
- 不显著增加 bytecode 或认知复杂度；
- 代码审查能够解释节省来源，而非依赖编译器偶然结果。

不采用：

- 为少量 gas 使用 assembly；
- storage 重排或自定义位打包；
- 删除安全检查；
- 增加批量 credit ABI；
- 为单个调用点引入 Library；
- 没有测量证据的微优化。

## 测试实施顺序

1. 保持现有 45 项测试为绿。
2. 为数值边界、反向币序、exact-output、partial fill、callback 与 Escrow invariant 添加测试。
3. 对每个拟修复问题执行红灯验证：未修复实现上测试必须失败。
4. 实施最小生产修复，使对应测试通过。
5. 运行相关 suite，再运行全量 Foundry 测试。
6. 记录修改前后 gas 和 bytecode 指标，只保留有证据的优化。
7. 检查 ABI、method identifiers、Hook permissions 和 CREATE2 规则无意外变化。

## 完成门禁

必须运行：

```bash
cd contracts
forge fmt --check
forge build
forge test
forge test --gas-report
```

还必须检查：

- 新增安全测试、fuzz/invariant 全部通过；
- ABI 与 method identifiers 无意外变化；
- `PoolConfig` getter 与字段顺序未变化；
- Hook permissions 与地址位规则未变化；
- LaunchToken creation code 未变化，除非另行批准；
- 后端 bindings、索引器和前端不需要同步；
- gas 优化有前后数据；
- 无未授权依赖、部署、广播或地址修改。

## 重新部署判定

代码仓库中的未来部署实现可以包含安全加固和已验证 gas 优化，但当前 Base Sepolia 部署默认保持不动。

只有同时满足以下条件时才建议重新部署：

1. 发现的问题影响资金安全、访问控制、债权守恒、流动性锁定、费用正确性或可被外部调用触发的协议不变量；
2. 问题能在当前部署字节码中复现；
3. 无法通过关闭前端入口、后端校验或其他非合约措施充分缓解；
4. 修复收益高于重新部署、迁移地址和更新集成的风险。

重新部署建议必须单独说明：严重性、利用条件、受影响合约、资金影响、缓解措施、新地址、部署顺序、bindings/前端/索引器更新和验证步骤。

纯 gas 优化、代码整理、注释或测试补强不构成重新部署理由。

## 交付物

- 安全边界、回归、fuzz 和 invariant 测试；
- 范围内的最小生产合约修改；
- 修改前后 gas 与 bytecode 对比；
- ABI、Hook permissions 和 CREATE2 兼容性检查结果；
- 是否需要重新部署的明确结论与证据；
- 未执行检查及其原因。
