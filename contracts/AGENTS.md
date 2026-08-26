# EVM 智能合约开发规范

本文件约束 `contracts/` 下后续所有 Solidity、Foundry 测试和部署脚本开发。当前代码用于确定适用的技术栈与风险边界，不代表现有实现天然符合本规范。新增代码必须遵守本规范；修改既有代码时，应让被触及区域向本规范收敛，不能以“原来就是这样”为由延续不安全模式。

## 技术栈与范围

- Solidity `0.8.26`，以 `foundry.toml` 为编译配置的唯一来源。
- EVM 目标为 Cancun，启用优化器与 `via_ir`；未经明确需求和完整回归验证，不得调整这些参数。
- 使用 Foundry 编译、测试和部署。
- 使用 OpenZeppelin、Uniswap v4 Core/Periphery 等已锁定依赖；不得复制其实现或直接修改 `lib/`。
- 当前协议包含 Registry、CREATE2 Factory、ERC-20 LaunchToken、Uniswap v4 Hook、FeeEscrow 和部署脚本。
- `src/` 是生产合约，`test/` 是测试，`script/` 是部署与验证脚本，`vendor/` 是外部或已验证材料，不得作为日常业务代码目录。

## 工作原则

- 只做完整解决当前任务所需的最小改动，不重构、重命名或格式化无关代码。
- 先阅读受影响合约、直接调用方、测试、部署脚本及后端 ABI 使用处，再开始修改。
- 优先复用已有库、类型和模式；不要引入没有当前需求的抽象、扩展点或依赖。
- 每个行为变更和缺陷修复都必须有测试；缺陷修复必须包含能复现问题的回归测试。
- 不得删除测试、放宽断言、降低 fuzz 次数或跳过安全检查来制造“通过”。
- 不确定权限、资金流、舍入方向或协议不变量时，必须停止并明确指出假设，不能自行补全业务规则。

## 协议兼容性边界

以下内容视为外部兼容性边界，除非任务明确批准破坏性变更并提供迁移方案，否则不得改变：

- external/public 函数选择器、参数、返回值和 payable 属性；
- 事件名称、参数顺序、类型及 `indexed` 标记；
- 自定义错误选择器和参数；
- CREATE2 salt、init code、部署者与地址推导规则；
- Hook permissions 与 Hook 地址低位标志；
- PoolKey、PoolId、币种顺序和 tick 计算；
- 已部署合约的状态语义、权限和资金归属；
- 后端、索引器、前端及部署记录使用的 ABI 和合约地址。

ABI 或事件发生变化时，必须在同一变更中检查并按需更新：

- `server/common/chain/bindings/` 下的 Go bindings；
- 索引器事件解析和 topic；
- 前端使用的 ABI；
- 部署脚本、部署清单、接口与测试；
- 任何依赖函数选择器、错误选择器或事件结构的集成。

不得手工编辑 Go bindings。先运行 `forge build`，再从 `server/` 运行 `go generate ./common/chain`，最后审查生成差异。

## 核心协议不变量

修改 Registry、Factory、Hook、Token、Escrow、费用或流动性逻辑时，必须明确列出并验证所有受影响不变量。当前协议至少包含：

- Launchpad 标识必须绑定链 ID 与规范化 slug，不得跨链混淆。
- Launchpad 的 owner、treasury、active 状态及所有权不能被事件重放或重复写入篡改。
- 发布必须验证配置版本、截止时间、launchpad 有效性和 quote 资产。
- 用户 salt 必须与调用者绑定；同一派生 salt 不得重复发布。
- CREATE2 预测地址必须与实际部署地址完全一致。
- Token/quote 排序必须由地址确定，PoolKey 与 PoolId 的构造必须保持确定性。
- Hook 地址权限位必须与 `getHookPermissions()` 一致，并在部署时验证。
- 每个 Pool 只能注册一次、播种一次；普通用户不能添加或移除永久流动性。
- 初始 Token 供应量、注入池子的数量及 Token 持有人必须符合发布规则，不得留下未说明的增发权限。
- 正常总费率为 150 bps，其中 protocol fee 为 100 bps、LaaS fee 为 50 bps；除非明确批准协议变更，不得漂移。
- protocol fee 的 Creator/Platform/Referrer 分配为 50%/30%/20%；无效 referrer 的份额归 Platform。
- anti-snipe 窗口、起始费率、精确输入/输出限制和线性衰减必须保持一致。
- 每次收费都必须满足费用守恒：各接收方份额与 surcharge 之和等于实际总收费，不得凭空增加或丢失。
- `FeeEscrow.totalOwed(currency)` 必须等于该币种所有用户债权之和，且不得超过 Escrow 在 PoolManager 中的可用余额。
- claim 必须先清账再进行外部交互；失败时整笔交易回滚，不能出现已扣债权但未付款。
- 任何 `unlockCallback` 只能由指定 PoolManager 调用，回调数据与币种必须按预期验证。

若任务有意修改上述不变量，必须在实现前明确标为协议级变更，并同步更新测试、后端、前端、部署配置和迁移说明。

## 访问控制与生命周期

- 每个写状态入口都必须逐一审查调用者、可调用阶段、重复调用和越权路径。
- 不得使用 `tx.origin` 进行认证或授权。
- 构造函数必须验证所有关键地址和配置；零地址仅在业务明确将其作为原生币标识时允许。
- 一次性 wiring、注册、播种和配置冻结必须 fail closed，并有重复调用测试。
- 不要假设 deployer、factory、hook、PoolManager、launchpad owner 和 treasury 是同一信任域。
- 权限变更或关键配置变化必须产生足够的可索引事件；不得仅靠日志文本表达状态变化。
- 删除权限检查、放宽角色、增加任意调用能力或引入管理员后门属于高风险变更，必须得到明确授权。

## 初始化器与可升级合约

当前生产合约以构造函数、immutable 和一次性 wiring 为主，不得为了“以后可能升级”擅自引入代理。

若未来任务明确引入代理或升级机制，则必须额外满足：

- 使用经过审查的标准代理模式，不自创代理协议。
- 实现合约构造函数必须调用 `_disableInitializers()`。
- `initialize`/`reinitialize` 必须有正确的版本和权限保护。
- 初始化或重初始化应幂等；追加数组、递增计数、铸币、转账、外部副作用等非幂等行为必须有明确安全论证和测试。
- 不得重排、删除、改型或复用已有存储槽；继承顺序变化也按存储布局变更处理。
- 必须保存并比较升级前后的 storage layout，测试升级、重复初始化和未授权升级。
- 升级必须说明治理者、延迟、暂停、回滚和迁移策略。

## 外部调用、回调与重入

- 对每个外部调用审查目标可信度、重入、返回值、失败语义、gas 假设和状态一致性。
- 优先遵循 Checks-Effects-Interactions；不能遵循时必须使用可靠重入保护并说明原因。
- 状态必须在外部调用前进入安全状态；外部调用失败必须整体回滚或进入明确可恢复状态。
- 不得假设 ERC-20 一定返回 `bool`、无手续费、无回调、无重基或完全遵循标准。
- 原生币与 ERC-20 路径必须分别测试；`address(0)` 作为 Currency 标识时不得与普通零地址校验混淆。
- `call`、`delegatecall`、任意目标调用和任意 calldata 属于高风险能力；没有明确需求和安全模型不得引入。
- Uniswap v4 `unlock`、`mint`、`burn`、`take`、`settle`、`sync` 和 Hook delta 必须作为一个完整结算系统审查，不能只验证局部余额。
- Hook callback 必须验证 Pool、Currency、方向、exact-input/exact-output、partial fill 和 delta 符号。
- 所有恶意接收者、恶意 token、重入 callback 和拒绝收款场景都应按风险添加测试。

## 数学、费用与精度

- 金额、费率、tick、流动性和 delta 的符号/位宽转换必须显式审查。
- 使用 `SafeCast` 或等价的有界转换；不得以裸强制转换掩盖溢出或符号错误。
- `unchecked` 只允许在边界已被证明且能带来实际收益时使用，并必须有注释与测试。
- 明确每个除法的舍入方向及余数归属。资金分账应采用“最后一方接收余数”或其他守恒方法，不能独立舍入后丢失资金。
- 费用计算必须覆盖零金额、极小金额、最大金额、边界费率、有效/无效 referrer 和 anti-snipe 时间边界。
- exact-output 与 exact-input 的费用基数不同，禁止为复用代码而混淆语义。
- TickMath、LiquidityAmounts、int128/int256/uint128 转换必须测试最小值、最大值、币种顺序翻转和极端 tickSpacing。
- 不得使用浮点数表达任何链上金额或比例。

## 签名、重放与时间

- 若新增签名授权，必须绑定 chain ID、验证合约、操作类型、参数、nonce 和 deadline。
- 使用 EIP-712 时必须验证 domain separator、类型哈希和合约钱包签名策略。
- nonce 必须一次性消费，失败重试不能产生双花或重复执行。
- 不得将 `block.timestamp` 用作高精度时钟或安全随机数。
- 使用 deadline 时必须测试刚好等于、刚好超过和极端时间值。
- 不得使用 blockhash、timestamp、prevrandao 等可操纵来源实现高价值随机性。

## CREATE2 与部署

- CREATE2 地址由 deployer、salt 和 init code hash 共同决定；修改构造参数、编译配置或 bytecode 都可能改变地址。
- 修改 Hook 构造参数、权限或创建代码时，必须重新验证地址权限位和 salt 挖掘逻辑。
- 部署脚本必须验证网络、chain ID、PoolManager、quote、treasury、tick 参数、预测地址、实际地址和 wiring。
- 私钥只能通过环境或安全签名器提供，不得写入源码、示例、日志或部署产物。
- 未经明确授权，不得广播交易、部署到公共网络、验证合约或更改生产配置。
- 部署后必须验证 bytecode、构造参数、immutable、Hook permissions、factory/escrow wiring 和关键不变量。
- `broadcast/`、部署清单和后端 deployment 配置必须避免泄露敏感信息，并保持地址一致。

## 代码与文档规范

- 保持当前文件布局、命名和 import 分组风格；使用 `forge fmt`，不要手工制造格式差异。
- 公开接口、复杂资金逻辑、非显然安全假设和协议不变量应使用 NatSpec 或精确注释说明“为什么”。
- 注释不得重复代码，也不得用注释替代可执行校验。
- 优先使用 custom errors；不要为了错误文案引入昂贵 revert string。
- 不得提交调试 `console`、被注释掉的旧代码、无归属 TODO 或临时绕过。
- 不得手工编辑 `out/`、`cache/`、第三方 `lib/` 或生成文件。

## 测试要求

开发时先运行最小相关测试，例如：

```bash
cd contracts
forge test --match-path test/LaunchHookV2Fees.t.sol
```

完成前至少运行：

```bash
cd contracts
forge fmt --check
forge build
forge test
```

根据风险补充以下测试：

- 单元测试：成功路径、预期 revert、事件、状态和余额变化；
- 边界测试：零值、最小值、最大值、时间边界、tick 边界和舍入；
- 权限测试：所有未授权主体与重复调用；
- fuzz 测试：费用守恒、地址组合、金额、时间、tick 和配置；
- invariant 测试：供应量、债权/余额、一次性状态、永久流动性、Pool/Token 映射；
- 集成测试：Factory → Token → PoolManager → Hook → FeeEscrow 的完整调用链；
- fork 测试：依赖真实部署地址或链上 PoolManager 行为时；
- 恶意合约测试：重入、拒收、异常 token 和 callback。

高风险合约变更不能只靠 happy-path 单测。若某类测试不适用，交付时必须说明原因。

## Code Review 规则

审查合约差异时，至少回答：

1. 改动影响哪些资金、权限、生命周期和协议不变量？
2. 是否改变 ABI、事件、错误、CREATE2 地址、Hook permissions 或后端解析？
3. 是否新增外部调用、回调、重入或任意执行能力？
4. 所有金额是否守恒，舍入余数归属是否明确？
5. exact-input/output、两种币序和原生币/ERC-20 是否均覆盖？
6. 一次性 wiring、注册、播种、salt 和 claim 能否重复执行？
7. 是否存在未授权调用、零地址、恶意合约或拒绝服务路径？
8. 测试是否能在没有实现修复时失败，并覆盖风险而非只覆盖行数？
9. 部署、Go bindings、索引器和前端是否需要同步？

发现未解决的高风险问题时不得批准。不得把“已有测试通过”当作安全证明。

## 完成条件

- 改动范围与任务一致，无无关重构或依赖升级。
- 相关单测、边界测试和安全测试已添加并通过。
- `forge fmt --check`、`forge build`、`forge test` 已通过。
- ABI 变化已同步生成 bindings，并检查后端/前端/索引器影响。
- 部署与 CREATE2 变化已验证预测地址和 wiring。
- 已说明未执行的 fork、静态分析或外部集成检查及原因。
- CI 必须到达最终状态；pending、跳过或已知 flaky 不能描述为通过。
