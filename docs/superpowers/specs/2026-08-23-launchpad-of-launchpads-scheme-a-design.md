# o1「Launchpad of Launchpads」方案 A 架构设计

**状态：** 已于 2026-08-23 确认  
**原始设计稿：** /Users/zhangshuai/Documents/Codex/2026-08-22/fan/outputs/o1-launchpad-of-launchpads-mvp-design.md

## 目标

实现一个可独立运行的 Base Sepolia MVP：用户连接钱包后可以创建带独立品牌的 Launchpad、获得稳定访问路径、在该 Launchpad 下发行 Token、通过链上证据验证 Token 与 Launchpad 的归属，并演示固定的 1% 协议费与独立的 0.5% LaaS fee。

## 范围决策

- Demo 仅运行在 Base Sepolia（链 ID 84532）。
- MVP 使用固定供应量、无管理权限的 ERC-20。o1 官方资料尚不能证明 Base Sepolia 支持与主网一致的 B20 创建能力，因此 B20 不能阻塞演示交付。
- PoolManager、Quoter、StateView、Universal Router 和 Permit2 使用 Base Sepolia 官方 Uniswap v4 部署。
- o1 生产 Launchpad API 不进入测试网执行链路。该 API 只处理已登记的生产合约 suite，无法索引或准备我们自定义测试网 Factory 的交易。
- 项目是一个全新 monorepo：Foundry 合约、Next.js 前端、go-zero REST Gateway、Proto/zRPC 业务服务、Indexer 和 PostgreSQL。
- MVP 使用 /launchpad/{slug}。通配符子域名、自定义域名、Launchpad Owner 分成、Marketplace、多链和独立 Subgraph 均不在首版范围。

## 总体架构

~~~text
浏览器 + 钱包
  ├── 钱包广播 → Base Sepolia 合约
  └── HTTPS/JSON → launchpad-api Gateway
                         └── zRPC/Proto → launchpad-rpc
                                              ├── PostgreSQL
                                              └── Base Sepolia RPC

Base Sepolia → indexer → PostgreSQL

Base Sepolia
  ├── LaunchpadRegistry
  ├── MultiTenantLaunchpadFactory
  ├── LaunchHookV2
  ├── FeeEscrow
  ├── 多个 LaunchToken
  └── 官方 Uniswap v4 合约
~~~

所有交易都由用户钱包签名。API 只保存品牌信息和索引视图，不能代替用户发行 Token 或转移费用。Launchpad 所有权以 Registry 为准；Token 归属以受信 Factory 的事件为准。

## 合约职责

### LaunchpadRegistry

- 通过 keccak256(abi.encode(chainId, normalizedSlug)) 派生 launchpadId。
- 保存 owner、treasury、active 和创建时间。
- 拒绝重复 ID、零地址 owner/treasury 和不规范 slug。
- 只有 owner 可以修改 treasury 或 active。
- Treasury 变更只影响未来 Pool；已创建 Pool 使用冻结配置。

### MultiTenantLaunchpadFactory

- 复用 o1 已验证 B20LaunchpadFactory 的发行流程和 Uniswap v4 集成结构。
- 测试网 MVP 将 B20 创建替换为固定供应量 LaunchToken。
- LaunchParams 增加 bytes32 launchpadId。
- 创建前要求 Registry 中对应 Launchpad 存在且 active。
- Treasury 从 Registry 读取；调用者不能提交费率或 Treasury 地址。
- 负责注册 Hook 配置、初始化 Pool、注入完整 Token 供应量，并发出带 launchpadId 的 Launched 事件。

### LaunchHookV2

- 复用 o1 LaunchHook 的核心行为：Factory-only 注册、Pool 配置冻结、永久单边流动性、anti-snipe 期间允许 exact-input、以配对资产收取费用和 Referrer 校验。
- 将 protocolFeeBps=100 与 laasFeeBps=50 分开记账。
- 只有 protocol fee 按 Creator/Platform/Referrer 的 50%/30%/20% 分配。
- 没有有效 Referrer 时，其份额归 Platform。
- LaaS fee 全部进入 LaaS Treasury。
- 开盘总费率在 16 秒内从 99% 衰减到 1.5%；高于 1.5% 的 surcharge 归 Platform。
- 禁止外部增加或移除流动性；anti-snipe 期间拒绝 exact-output。

### FeeEscrow

- 保留 o1 通用 owed[recipient][currency] pull-payment 模型。
- 只有新 Hook 可以增加余额。
- claim(recipient,currency) 向原收款人支付；claimTo(currency,to) 允许余额所有者指定接收地址。
- 使用 checks-effects-interactions 与重入保护。

## 费用不变量

正常阶段，对配对资产输入 amount：

~~~text
protocolFee = floor(amount × 100 / 10_000)
laasFee     = floor(amount ×  50 / 10_000)
creator     = floor(protocolFee × 5_000 / 10_000)
referrer    = valid ? floor(protocolFee × 2_000 / 10_000) : 0
platform    = protocolFee - creator - referrer
laas        = laasFee
~~~

整数除法余数保留在 Platform 份额，保证所有 credit 之和始终等于实际收取金额。

## 后端职责

- launchpad-api 是 go-zero REST Gateway：Handler 解析 HTTP 请求，Gateway Logic 调用生成的 zRPC Client；Gateway 负责 Cookie、CORS、限流和 HTTP 错误映射，不直接连接 PostgreSQL 或 Base RPC。
- launchpad-rpc 由 launchpad.proto 和 goctl 生成，业务 Logic 负责身份验证、品牌配置、交易准备和查询；其 ServiceContext 统一构造 PostgreSQL、Base RPC、合约客户端和 Model。
- Indexer 是同一 Go Module 下的独立 go-zero Service，拥有自己的 Config、ServiceContext 和消费循环，与 RPC 共用 common/chain、common/model 和 PostgreSQL Schema。
- Gateway 只调用 launchpad-rpc；Indexer 不调用 Gateway 或 RPC，只写同一 PostgreSQL。MVP 不使用 Etcd、Redis 或消息队列，Gateway 通过 endpoints 直连 RPC。
- 生成一次性钱包签名 challenge，并签发 HttpOnly Session Cookie。
- 只有在验证 Registry 交易成功、事件匹配且 owner 与登录钱包一致后，才保存品牌信息。
- 执行 slug 规范化、唯一性和保留词校验。
- 只索引配置中的 Factory、Hook、Escrow 和 PoolManager。
- 等待 2 个确认后标记为 confirmed。
- 保存区块哈希并保留回看窗口，短重组发生时删除孤块数据并重新播放。
- 提供 Launchpad、Token 列表、Token 详情、费用余额和交易状态查询。

## 配置约定

- Gateway 加载 service/launchpad/api/etc/launchpad-api.yaml；RPC 加载 service/launchpad/rpc/etc/launchpad-rpc.yaml；Indexer 加载 service/indexer/etc/indexer.yaml。
- Gateway YAML 配置 HTTP、zRPC endpoints、Session、CORS 和日志；RPC YAML 配置 PostgreSQL、Base RPC、部署清单与日志；Indexer YAML 配置 PostgreSQL、Base RPC、部署清单、确认数、扫描批次和重组回看窗口。
- 通过 go-zero conf.MustLoad(..., conf.UseEnv()) 加载配置；YAML 中使用 ${DATABASE_URL}、${BASE_SEPOLIA_RPC_URL}、${SESSION_SECRET} 等占位表达式注入敏感值。
- 数据库密码、RPC 密钥和 Session Secret 不写入仓库；非敏感默认值直接写入 YAML。
- deployments/base-sepolia.json 仍是合约地址的唯一来源；YAML 只配置 DeploymentFile 路径，避免地址复制后不一致。

## 对外 REST 路由

MVP 复用 o1 的资源命名，只实现闭环需要的路由：

~~~text
GET  /v1/health
GET  /v1/config
GET  /v1/tokens
GET  /v1/tokens/{chain_id}/{token_address}
GET  /v1/transactions/{chain_id}/{tx_hash}
POST /v1/launches/prepare
POST /v1/swaps/quote
POST /v1/swaps/prepare
POST /v1/claims/fees/prepare
~~~

多租户能力新增：

~~~text
GET  /v1/launchpads
POST /v1/launchpads
GET  /v1/launchpads/{slug}
GET  /v1/launchpads/{slug}/tokens
POST /v1/auth/challenge
POST /v1/auth/verify
~~~

## 数据模型

- launchpads：链 ID、链上 ID、slug、品牌、owner、treasury、Registry 交易和时间戳。
- token_launches：Launchpad 外键、Token、Pool ID、creator、quote、Factory、交易、区块身份和确认状态。
- chain_transactions：交易生命周期和失败原因。
- indexer_checkpoints：每个 chain/factory 的下一个区块与最后 canonical 区块哈希。
- auth_nonces：一次性 challenge、钱包地址、过期和使用时间。

唯一约束为 chain+slug、chain+launchpadId、chain+token、chain+pool 和 chain+txHash。

## 前端流程

### 创建 Launchpad

1. 连接钱包并切换到 Base Sepolia。
2. 使用一次性 challenge 签名登录。
3. 在本地规范化和校验品牌参数。
4. 调用 LaunchpadRegistry.createLaunchpad。
5. 把已确认的交易哈希和品牌配置提交给 API。
6. 跳转至 /launchpad/{slug}。

### 发行 Token

1. 从 API 加载 Launchpad，并与 Registry 读取结果交叉验证。
2. 展示 Token 身份、固定供应量、配对资产、永久流动性、anti-snipe 和 1.5% 正常费率。
3. 调用 POST /v1/launches/prepare，由 RPC 在同一区块读取 Factory 配置、模拟并返回未签名交易。
4. 钱包检查并原样广播返回的 to、data 和 value。
5. 从 Receipt 解析 Launched，立即展示 Token 和 Pool。
6. 轮询本地交易接口，直到 Indexer 确认。

### 浏览与领取

- Launchpad 页面只查询本地 Indexer 认证且绑定到该 launchpadId 的 Token。
- Token 页面组合数据库身份数据与链上/Uniswap 实时读取。
- Claim 通过 POST /v1/claims/fees/prepare 获取未签名交易，再由连接的钱包广播。

## 失败处理

- 钱包链错误：禁止写操作并提供切换 Base Sepolia 的入口。
- Factory 配置过期或 deadline 失效：重新读取配置并要求用户重新审阅、签名。
- slug 重复或命中保留词：交易提交前拒绝。
- Registry 交易与提交品牌不匹配：拒绝写入数据库。
- Receipt 已确认但尚未索引：显示 confirming，持续轮询并允许手动刷新。
- 交易回滚：已知错误展示解码后的 custom error，同时保留表单数据。
- RPC/API 不可用：保留用户输入，不得显示成功状态。

## 验证策略

- Foundry unit/fuzz/invariant 测试覆盖权限、费用计算、anti-snipe 端点、Referrer、Pool 配置冻结、流动性锁定、Claim 和重入。
- Go 测试覆盖签名重放、交易验证、事件解析、幂等、Checkpoint 恢复和重组重放。
- 前端测试覆盖 slug 规范化、费用展示、Receipt 解析和错误网络禁用。
- Base Sepolia smoke 脚本创建 Launchpad、发行 Token、执行 Swap、验证四类费用、领取一笔余额，并确认 Token 归属索引。

## 明确延期

- Base Sepolia B20 创建。
- o1 生产 API 集成。
- 通配符子域名和自定义域名。
- Launchpad Owner LaaS 分成。
- 动态费率、订阅、Marketplace、排名分析、多链、Subgraph、治理、可升级合约和生产审计。
