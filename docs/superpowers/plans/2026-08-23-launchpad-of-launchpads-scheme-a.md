# Launchpad of Launchpads 方案 A 实施计划

> **给执行者：** 实施本计划时必须使用 superpowers:subagent-driven-development（推荐）或 superpowers:executing-plans，按任务逐项执行，并使用 - [ ] 复选框记录进度。

**目标：** 交付一个可独立运行的 Base Sepolia MVP：用户可以创建品牌化 Launchpad、在永久 Uniswap v4 Pool 中发行固定供应量 Token、通过 launchpadId 验证归属，并演示固定的 1% 协议费与 0.5% LaaS fee。

**架构：** 全新 monorepo 包含 Foundry 合约、Next.js 前端、go-zero REST Gateway、Proto/zRPC 业务服务和独立 Indexer，数据存储使用 PostgreSQL。Gateway 提供 o1 风格 REST 路由并调用 RPC；RPC 处理业务、数据库和链上读取；Indexer 扫链写库且不调用另外两个服务。钱包广播 API 准备的未签名交易；o1 生产 API 不进入测试网核心链路。

**技术栈：** Solidity 0.8.26、Foundry、OpenZeppelin、Uniswap v4、Go 1.24、go-zero、goctl、go-ethereum、pgx v5、PostgreSQL 16、Next.js 15、TypeScript、wagmi、viem、TanStack Query、Vitest、Playwright。

**设计依据：** docs/superpowers/specs/2026-08-23-launchpad-of-launchpads-scheme-a-design.md；原始方案为 /Users/zhangshuai/Documents/Codex/2026-08-22/fan/outputs/o1-launchpad-of-launchpads-mvp-design.md。

## 全局约束

- 目标网络仅为 Base Sepolia，chainId=84532。
- Token 固定供应量为 1_000_000_000 ether，完整供应量进入永久单边流动性。
- 每个 Pool 冻结 protocolFeeBps=100、laasFeeBps=50，总正常费率 150 bps。
- Creator/Platform/Referrer 只分配 protocol fee，比例为 5000/3000/2000；无有效 Referrer 时其份额归 Platform。
- Anti-snipe 在 16 秒内从 9900 bps 衰减到 150 bps，高于 150 bps 的 surcharge 归 Platform。
- 品牌信息链下保存；owner、treasury、active、Pool 费用策略和 Token 归属以链上为准。
- 浏览器不得包含 o1 API Key 或服务端签名密钥。
- MVP 不包含 o1 生产 API 依赖、Base B20 依赖、通配符/自定义域名、Owner 分成、动态费率、多链、Subgraph、代理升级和管理员提款。
- Base Sepolia Uniswap v4 地址只保存在 deployments/base-sepolia.json，不在多处硬编码。
- 任何涉及资金路径的分支必须先有可运行的失败测试。

## 文件结构

~~~text
contracts/
  src/
    LaunchpadRegistry.sol
    LaunchToken.sol
    FeeEscrow.sol
    LaunchHookV2.sol
    MultiTenantLaunchpadFactory.sol
  script/
    DeployBaseSepolia.s.sol
    SmokeBaseSepolia.s.sol
  test/*.t.sol
  vendor/o1/
server/
  common/
    chain/
    model/
  service/
    launchpad/api/
      etc/launchpad-api.yaml
      launchpad.api
      launchpad.go
      internal/
        config/
        handler/
        logic/
        svc/
        types/
    launchpad/rpc/
      etc/launchpad-rpc.yaml
      launchpad.proto
      launchpad.go
      launchpadclient/
      internal/
        config/
        logic/
        server/
        svc/
    indexer/
      etc/indexer.yaml
      indexer.go
      internal/
        config/
        logic/
        svc/
  migrations/001_init.sql
web/
  app/
  components/
  lib/
deployments/base-sepolia.json
docker-compose.yml
Makefile
docs/demo/base-sepolia-runbook.md
~~~

合约负责协议不变量；Gateway 负责 HTTP 协议与 Cookie；Proto RPC 负责身份验证、品牌配置、交易准备和查询；Indexer 负责链上同步；前端负责展示和钱包交互。Gateway Handler 不含业务逻辑、不连接数据库；数据库/RPC/合约依赖只在 RPC 或 Indexer 的 ServiceContext 构造；只有 RPC 和 Indexer 都使用的链上/Model 代码进入 common。

---

### 任务 1：仓库基线与 o1 源码固化

**文件：**
- 新建：.gitignore、.env.example、Makefile
- 可选新建：docker-compose.yml（仅供没有现成 PostgreSQL 的开发者使用）
- 新建：deployments/base-sepolia.json
- 新建：contracts/vendor/o1/SOURCE.md、contracts/vendor/o1/verified/

**接口与产物：**
- 生成所有层共用的测试网部署清单。
- 固化 Factory、Hook 和 Escrow 的浏览器验证源码及来源信息。

- [x] **步骤 1：初始化 Git 并忽略敏感/生成文件**

执行：git init

忽略 .env、.next、node_modules、contracts/cache、contracts/out、coverage、server/bin 和编辑器文件。

- [x] **步骤 2：获取 o1 已验证生产源码**

从区块浏览器下载：

~~~text
B20LaunchpadFactory  0xa52ad458cE0282a971ecC71C051A32f28946bb9F
LaunchHook           0x985C14BAa2a18316ffdA0aeFB3a632fAdfcA2AcC
FeeEscrow            0xa2cBD9065cec93c443CAFb0837A62800EE7C4A84
~~~

在 SOURCE.md 记录浏览器 URL、编译器版本、优化参数、获取日期和 SHA-256。vendor 文件只读，不直接修改。

- [x] **步骤 3：写入 Base Sepolia 部署清单**

~~~json
{
  "chainId": 84532,
  "poolManager": "0x05E73354cFDd6745C338b50BcFDfA3Aa6fA03408",
  "universalRouter": "0x492e6456d9528771018deb9e87ef7750ef184104",
  "positionManager": "0x4b2c77d209d3405f41a037ec6c77f7f5b8e2ca80",
  "stateView": "0x571291b572ed32ce6751a2cb2486ebee8defb9b4",
  "quoter": "0x4a6513c898fe1b2d0e78d3b0e0a4a151589b1cba",
  "permit2": "0x000000000022D473030F116dDEE9F6B43aC78BA3",
  "registry": null,
  "factory": null,
  "hook": null,
  "feeEscrow": null,
  "deploymentBlock": null
}
~~~

- [x] **步骤 4：准备已有 PostgreSQL 与环境变量模板**

复用本机现有 Docker PostgreSQL，新建独立数据库 o1_launchpad_mvp；Migration 在任务 8 编写后执行。定义 DATABASE_URL、BASE_SEPOLIA_RPC_URL、BASE_SEPOLIA_WS_URL、SESSION_SECRET、NEXT_PUBLIC_API_URL 和 NEXT_PUBLIC_RPC_URL，但不写真实值。docker-compose.yml 不是启动前提。

- [x] **步骤 5：检查密钥并提交**

执行：

~~~bash
git grep -nE '(PRIVATE_KEY=0x|o1_launch_[a-f0-9]{8}_)' -- . ':!docs/superpowers/plans'
~~~

预期：无匹配。

提交信息：chore: establish launchpad MVP baseline

---

### 任务 2：实现 LaunchpadRegistry

**文件：**
- 新建：contracts/foundry.toml、contracts/remappings.txt
- 新建：contracts/src/LaunchpadRegistry.sol
- 测试：contracts/test/LaunchpadRegistry.t.sol

**接口：**
- deriveId(string) returns (bytes32)
- createLaunchpad(string normalizedSlug,address treasury) returns (bytes32)
- getLaunchpad(bytes32) returns (owner,treasury,active,createdAt)
- setTreasury(bytes32,address)
- setActive(bytes32,bool)
- 事件 LaunchpadCreated(bytes32 indexed id,address indexed owner,address indexed treasury,string slug)

- [x] **步骤 1：先写失败测试**

~~~solidity
function testCreateStoresOwnerTreasuryAndActive() public {
    vm.prank(owner);
    bytes32 id = registry.createLaunchpad("ai-pad", treasury);
    (address o,address t,bool active,) = registry.getLaunchpad(id);
    assertEq(o, owner);
    assertEq(t, treasury);
    assertTrue(active);
}

function testOnlyOwnerMutates() public {
    vm.prank(owner);
    bytes32 id = registry.createLaunchpad("ai-pad", treasury);
    vm.expectRevert(LaunchpadRegistry.NotOwner.selector);
    vm.prank(attacker);
    registry.setActive(id, false);
}
~~~

- [x] **步骤 2：确认测试失败**

执行：cd contracts && forge test --match-path test/LaunchpadRegistry.t.sol -vv

预期：LaunchpadRegistry 尚不存在导致编译失败。

- [x] **步骤 3：实现最小 Registry**

ID 使用 keccak256(abi.encode(block.chainid, normalizedSlug))。slug 仅允许 3–32 位小写 ASCII 字母、数字和内部单个连字符；拒绝首尾连字符、连续连字符、重复 ID 和零 Treasury。使用 custom error。

- [x] **步骤 4：补充重复 ID、非法 slug、零 Treasury、事件和 Owner 修改测试**

- [x] **步骤 5：运行并提交**

执行：cd contracts && forge test --match-path test/LaunchpadRegistry.t.sol -vv

提交信息：feat: add onchain launchpad registry

---

### 任务 3：实现 LaunchToken 与 FeeEscrow

**文件：**
- 新建：contracts/src/LaunchToken.sol
- 新建：contracts/src/FeeEscrow.sol
- 测试：contracts/test/LaunchToken.t.sol
- 测试：contracts/test/FeeEscrow.t.sol

**接口：**
- LaunchToken 构造参数为 name、symbol、contractURI、supply、recipient；铸造后不存在管理权限。
- FeeEscrow 提供 owed、credit、claim、claimTo。
- 只有构造时绑定的 Hook 可以 credit。

- [x] **步骤 1：写 LaunchToken 失败测试**

~~~solidity
function testMintsCompleteSupplyOnce() public {
    LaunchToken token = new LaunchToken(
        "Alpha", "ALPHA", "ipfs://meta", 1_000_000_000 ether, recipient
    );
    assertEq(token.totalSupply(), 1_000_000_000 ether);
    assertEq(token.balanceOf(recipient), token.totalSupply());
}
~~~

- [x] **步骤 2：用 OpenZeppelin ERC20 实现一次性构造铸造**

提供 contractURI，但不提供 mint、burn、pause、owner、role 或 metadata 修改接口。

- [x] **步骤 3：写 FeeEscrow 失败测试**

~~~solidity
function testClaimPaysAndClears() public {
    currency.mint(address(escrow), 10 ether);
    vm.prank(hook);
    escrow.credit(recipient, address(currency), 10 ether);
    escrow.claim(recipient, address(currency));
    assertEq(currency.balanceOf(recipient), 10 ether);
    assertEq(escrow.owed(recipient, address(currency)), 0);
}
~~~

- [x] **步骤 4：移植 o1 Escrow 语义**

使用 PoolManager ERC-6909 claim、ReentrancyGuard、CEI 和 custom error；Swap 期间接收 claim，领取时在独立 unlock 中 burn claim 并 take 真实资产。拒绝零地址、零金额和超过实际 claim 余额的记账。直接 SafeERC20 转账与 v4 Router 的 settle 时序不兼容，不采用。

- [x] **步骤 5：增加 Fuzz 与恶意 Token 重入测试**

验证已支付金额加未领取余额永不超过可赎回 claim，Claim 不会重复支付。

- [x] **步骤 6：运行并提交**

~~~bash
cd contracts
forge test --match-path test/LaunchToken.t.sol -vv
forge test --match-path test/FeeEscrow.t.sol -vv
~~~

提交信息：feat: add immutable token and fee escrow

---

### 任务 4：实现 LaunchHookV2 费用引擎

**文件：**
- 新建：contracts/src/LaunchHookV2.sol
- 测试：contracts/test/LaunchHookV2Fees.t.sol

**接口：**
- PoolConfig 冻结 creator、protocolTreasury、laasTreasury、费率、分配比例、launchTime 和 Token 顺序。
- currentTotalFeeBps(poolId,timestamp) 返回 150–9900。
- 内部拆分结果包含 creator、protocol、referrer、laas、surcharge，合计必须等于实际收费。

- [x] **步骤 1：写精确金额失败测试**

1,000 单位配对资产、有 Referrer 时应为 Creator 5、Platform 3、Referrer 2、LaaS 5；无 Referrer 时为 Creator 5、Platform 5、Referrer 0、LaaS 5。

- [x] **步骤 2：写 anti-snipe 失败测试**

launchTime 时为 9900 bps，中间单调不增，launchTime+16 及之后精确等于 150 bps。

- [x] **步骤 3：确认失败**

执行：cd contracts && forge test --match-path test/LaunchHookV2Fees.t.sol -vv

- [x] **步骤 4：移植 o1 Hook 费用计算，PoolManager callback 在任务 5 集成**

本任务冻结费用计算、Referrer 校验与 protocol/LaaS 拆分；任务 5 集成 PoolManager callback、配对资产收费和 Trade 事件。整数余数和开盘 surcharge 进入 Platform。

- [x] **步骤 5：增加 Fuzz 不变量**

对 uint128 金额、时间戳和 Referrer 状态验证：credit 合计等于收费、LaaS 不进入 Creator/Referrer、费率始终在 150–9900。

- [x] **步骤 6：运行并提交**

执行：cd contracts && forge test --match-path test/LaunchHookV2Fees.t.sol --fuzz-runs 1000

提交信息：feat: separate protocol and LaaS fees

---

### 任务 5：实现 Hook Pool 门控和永久流动性

**文件：**
- 修改：contracts/src/LaunchHookV2.sol
- 测试：contracts/test/LaunchHookV2Pool.t.sol

**接口：**
- setFactoryOnce(address) 永久绑定一个 Factory。
- registerPool(PoolKey,PoolConfigInput) 和 seed(PoolKey,uint256) 只能由 Factory 调用一次。
- 外部增加/移除流动性，以及 anti-snipe 期间 exact-output 必须回滚。

- [x] **步骤 1：写非 Factory、重复注册、重复 Seed、非法 PoolManager 和流动性修改失败测试**

- [x] **步骤 2：写 Uniswap v4 集成失败测试**

初始化 Pool、注入完整供应量、执行 exact-input 买卖，并验证配对资产进入 Escrow。

- [x] **步骤 3：移植 o1 的权限位、注册、Seed 和流动性锁定逻辑**

- [x] **步骤 4：验证 Hook 地址权限位并运行测试**

执行：cd contracts && forge test --match-path 'test/LaunchHookV2*.t.sol' -vv

- [x] **步骤 5：提交**

提交信息：feat: enforce frozen launch pools

---

### 任务 6：实现 MultiTenantLaunchpadFactory

**文件：**
- 新建：contracts/src/MultiTenantLaunchpadFactory.sol
- 测试：contracts/test/MultiTenantLaunchpadFactory.t.sol

**接口：**
- LaunchParams 包含 launchpadId、name、symbol、contractURI、salt、quote、expectedConfigVersion、deadline。
- 依赖 Registry、Hook、PoolManager、Quote 配置、供应量、流动性 Bands 和 Treasury Policy。
- 发出 Launched(token,poolId,launchpadId,creator,quote,supply,tickSpacing)。

- [x] **步骤 1：写未知/暂停 Launchpad、过期配置、过期 deadline 和 salt 重放失败测试**

- [x] **步骤 2：写成功发行失败测试**

验证 Token 供应量、Pool ID、Creator、Treasury、100/50 费率和事件中的 launchpadId。

- [x] **步骤 3：移植 o1 Factory 流程**

保留 configVersion、salt 防重放、Quote 注册、PoolKey/Tick/Bands、Pool 初始化、Hook 注册和 Seed；用 new LaunchToken 替换 B20 创建，并增加 Registry 校验。

- [x] **步骤 4：增加不可变性测试**

用户不能提交费率/Treasury；Registry Treasury 修改只影响未来 Pool；每个 Token 只能对应一个 Pool 和 launchpadId。

- [x] **步骤 5：运行完整合约测试**

执行：

~~~bash
cd contracts
forge fmt --check
forge test -vv
forge test --fuzz-runs 1000
~~~

- [x] **步骤 6：提交**

提交信息：feat: launch tokens under registered tenants

---

### 任务 7：确定性部署到 Base Sepolia

**文件：**
- 新建：contracts/script/DeployBaseSepolia.s.sol
- 新建：contracts/script/SmokeBaseSepolia.s.sol
- 测试：contracts/test/DeploymentWiring.t.sol
- 修改：deployments/base-sepolia.json

**接口与不变量：**
- Hook 地址必须满足 Uniswap v4 permission bits。
- Escrow 绑定的 Hook、Hook 绑定的 Factory、Factory 引用的 Registry/Hook/PoolManager 必须一致。

- [x] **步骤 1：写部署 Wiring 失败测试**

- [x] **步骤 2：实现 CREATE2 部署**

寻找合法 Hook salt，预测 Hook 地址；先使用预测地址部署 Escrow，再通过 CREATE2 部署 Hook，然后部署 Registry 和 Factory，最后一次性绑定 Factory。任何预测不一致立即回滚。

- [x] **步骤 3：本地验证**

执行：cd contracts && forge test --match-path test/DeploymentWiring.t.sol -vv

- [ ] **步骤 4：部署并验证源码**

~~~bash
cd contracts
forge script script/DeployBaseSepolia.s.sol:DeployBaseSepolia \
  --rpc-url "$BASE_SEPOLIA_RPC_URL" --broadcast --verify
~~~

- [ ] **步骤 5：执行 Smoke Launch 并更新部署清单**

~~~bash
cd contracts
forge script script/SmokeBaseSepolia.s.sol:SmokeBaseSepolia \
  --rpc-url "$BASE_SEPOLIA_RPC_URL" --broadcast
~~~

预期：产生 Registry 事件、Launch 事件、初始化 Pool 和完整 Supply Seed。

- [ ] **步骤 6：提交**

提交信息：deploy: publish Base Sepolia launch suite

---

### 任务 8：go-zero Gateway、Proto RPC、Indexer 与 PostgreSQL Model 基线

**文件：**
- 新建：server/go.mod
- 新建：server/migrations/001_init.sql
- 新建：server/common/model/launchpadmodel.go、tokenlaunchmodel.go、transactionmodel.go、checkpointmodel.go、authnoncemodel.go
- 新建：server/service/launchpad/api/etc/launchpad-api.yaml
- 新建：server/service/launchpad/api/internal/config/config.go
- 新建：server/service/launchpad/api/internal/svc/servicecontext.go
- 新建：server/service/launchpad/rpc/launchpad.proto
- 新建：server/service/launchpad/rpc/etc/launchpad-rpc.yaml
- 新建：server/service/launchpad/rpc/internal/config/config.go
- 新建：server/service/launchpad/rpc/internal/svc/servicecontext.go
- 新建：server/service/indexer/etc/indexer.yaml
- 新建：server/service/indexer/internal/config/config.go
- 新建：server/service/indexer/internal/svc/servicecontext.go
- 测试：server/common/model/model_test.go、三个服务的 internal/config/config_test.go

**接口：**
- Model 提供 CreateLaunchpad、GetLaunchpadBySlug、ListLaunchpadTokens、UpsertTransaction、ApplyLaunchEvent、GetCheckpoint、SaveCheckpoint。
- Gateway Config 嵌入 rest.RestConf，并包含 LaunchpadRpc 与 Auth；ServiceContext 只构造 zRPC Client 和 HTTP 中间件依赖。
- RPC Config 嵌入 zrpc.RpcServerConf，并包含 Database、Chain、DeploymentFile；ServiceContext 构造 PostgreSQL、Base RPC、合约 Client 和 Model。
- Indexer Config 嵌入 service.ServiceConf，并包含 Database、Chain、Indexer、DeploymentFile；ServiceContext 构造独立连接池和 RPC Client。
- RPC 与 Indexer 连接同一个 PostgreSQL；Gateway 不连接 PostgreSQL。
- 表：launchpads、token_launches、chain_transactions、indexer_checkpoints、auth_nonces

- [x] **步骤 1：先定义 Proto 并生成 RPC**

launchpad.proto 使用 package launchpad、go_package=./pb/launchpad（确保下述 goctl 命令直接生成到预期目录），定义 MVP 所需方法：Health、GetConfig、ListLaunchpads、CreateLaunchpad、GetLaunchpad、ListLaunchpadTokens、ListTokens、GetToken、GetTransaction、PrepareLaunch、QuoteSwap、PrepareSwap、PrepareFeeClaim、CreateAuthChallenge、VerifyAuthSignature。

执行：

~~~bash
cd server/service/launchpad/rpc
goctl rpc protoc launchpad.proto --go_out=. --go-grpc_out=. --zrpc_out=.
~~~

- [x] **步骤 2：定义 REST DSL 并生成 Gateway**

先编写 server/service/launchpad/api/launchpad.api，再执行：

~~~bash
cd server/service/launchpad/api
goctl api go -api launchpad.api -dir .
~~~

Indexer 按 go-zero Service 约定手工建立 indexer.go、etc、internal/config、internal/logic、internal/svc；它不对外提供接口，也不调用 Gateway/RPC，因此不生成 HTTP 或 gRPC Server。

- [x] **步骤 3：编写 Migration**

地址和哈希使用 bytea，原始金额使用 numeric(78,0)，时间使用 timestamptz。建立 chain+slug、chain+launchpadId、chain+token、chain+pool、chain+txHash 唯一约束。

- [x] **步骤 4：编写三个 YAML 配置**

API 配置：

~~~yaml
Name: launchpad-api
Host: 0.0.0.0
Port: 8888
Mode: dev
LaunchpadRpc:
  Endpoints:
    - 127.0.0.1:8080
Auth:
  SessionSecret: ${SESSION_SECRET}
  ChallengeTTLSeconds: 300
~~~

RPC 配置：

~~~yaml
Name: launchpad-rpc
ListenOn: 0.0.0.0:8080
Mode: dev
Database:
  DataSource: ${DATABASE_URL}
  MaxOpenConns: 20
  MaxIdleConns: 5
Chain:
  ChainId: 84532
  HttpRpc: ${BASE_SEPOLIA_RPC_URL}
  WsRpc: ${BASE_SEPOLIA_WS_URL}
  Confirmations: 2
DeploymentFile: ../../../../deployments/base-sepolia.json
~~~

Indexer 配置：

~~~yaml
Name: launchpad-indexer
Mode: dev
Database:
  DataSource: ${DATABASE_URL}
  MaxOpenConns: 10
  MaxIdleConns: 2
Chain:
  ChainId: 84532
  HttpRpc: ${BASE_SEPOLIA_RPC_URL}
  WsRpc: ${BASE_SEPOLIA_WS_URL}
  Confirmations: 2
DeploymentFile: ../../../deployments/base-sepolia.json
Indexer:
  BatchSize: 500
  PollInterval: 2s
  ReorgLookback: 64
~~~

- [x] **步骤 5：写 Config 加载失败测试**

使用 conf.MustLoad(path,&c,conf.UseEnv())；测试 Gateway 缺少 RPC Endpoint/Session Secret、RPC 缺失 DATABASE_URL/错误 chainId、Indexer Confirmations<1/BatchSize>500 时启动失败。密码和 RPC Key 只通过环境变量展开。

- [x] **步骤 6：写 Model 集成失败测试**

覆盖重复归属冲突、事件幂等重放、pending→confirming→confirmed 和孤块数据删除。

- [x] **步骤 7：确认失败**

~~~bash
cd server
go test ./common/model ./service/launchpad/api/internal/config ./service/launchpad/rpc/internal/config ./service/indexer/internal/config -count=1
~~~

测试使用环境变量 DATABASE_URL 连接现有 Docker PostgreSQL 中的 o1_launchpad_mvp 数据库，不负责启动或重建 PostgreSQL 容器。

- [x] **步骤 8：实现 Model 与三个 ServiceContext**

ApplyLaunchEvent 必须在同一数据库事务中插入事件、更新交易并推进 Checkpoint。若唯一键冲突但不可变归属字段不一致，返回错误而不是覆盖。Gateway ServiceContext 只创建 launchpadclient.Launchpad；RPC ServiceContext 创建数据库、链上 Client、Deployment 和 Model；Indexer ServiceContext 创建独立数据库/RPC Client。Handler/Logic 禁止自行连接数据库。

- [x] **步骤 9：运行并提交**

提交信息：feat: add launchpad persistence model

---

### 任务 9：Proto 业务、REST Gateway 与品牌配置

**文件：**
- 修改：server/service/launchpad/api/launchpad.api
- 修改：server/service/launchpad/rpc/launchpad.proto
- 新建：server/service/launchpad/api/internal/handler/auth/*、launchpad/*
- 新建：server/service/launchpad/api/internal/logic/auth/*、launchpad/*（只做 HTTP/zRPC 转换）
- 新建：server/service/launchpad/rpc/internal/logic/createauthchallengelogic.go、verifyauthsignaturelogic.go、createlaunchpadlogic.go、getlaunchpadlogic.go、listlaunchpadslogic.go、listlaunchpadtokenslogic.go
- 新建：server/service/launchpad/api/internal/middleware/sessionmiddleware.go
- 新建：server/common/chain/registry.go
- 测试：RPC Logic、Gateway Logic 与 Handler 测试

**接口：**
- POST /v1/auth/challenge
- POST /v1/auth/verify
- POST /v1/launchpads
- GET /v1/launchpads/{slug}
- GET /v1/launchpads
- GET /v1/launchpads/{slug}/tokens

- [x] **步骤 1：写签名过期、Domain/Chain/Address 错误和重放失败测试**

EIP-191 消息包含 Domain、URI、chainId=84532、Address、随机 Nonce、签发时间和 5 分钟有效期。

- [x] **步骤 2：在 RPC Logic 实现认证**

Gateway Handler 只绑定 Request；Gateway Logic 调用 zRPC。RPC Logic 通过 svcCtx.AuthNonceModel 保存 Nonce Hash并原子消费。验证成功后 Gateway 使用 RPC 返回的 Session Claims 签发 HttpOnly、Secure、SameSite=Lax Cookie。不引入额外认证框架。

- [x] **步骤 3：写 Registry Receipt 验证失败测试**

拒绝回滚交易、错误 Registry、slug 事件不匹配、Owner 不匹配、错误链和不足 2 个确认。

- [x] **步骤 4：在 RPC Logic 实现品牌校验和保存**

Name 1–80 字符，Description 不超过 1000 UTF-8 bytes，Logo 必须 HTTPS，颜色必须 #RRGGBB，slug 规则与合约一致。仅在一个匹配的 LaunchpadCreated 事件属于登录钱包时保存。

- [x] **步骤 5：运行并提交**

执行：cd server && go test ./service/launchpad/api/... ./service/launchpad/rpc/... ./common/chain/... -count=1

提交信息：feat: verify wallet-owned launchpad creation

---

### 任务 10：go-zero Indexer 与查询 API

**文件：**
- 新建：server/common/chain/abi.go、fees.go 和生成的 Binding
- 新建：server/service/indexer/internal/logic/indexerlogic.go、indexerlogic_test.go
- 修改：server/service/indexer/indexer.go
- 新建：server/service/launchpad/api/internal/handler/{token,transaction,fee}/*
- 新建：server/service/launchpad/api/internal/logic/{token,transaction,fee,launch,swap}/*（只调用 zRPC）
- 新建：server/service/launchpad/rpc/internal/logic/{config,token,transaction,fee,launch,swap}/*

**接口：**
- 只从 deploymentBlock 起索引配置中的 Factory/Hook/Escrow/PoolManager。
- GET /v1/launchpads/{slug}/tokens
- GET /v1/tokens/{address}
- GET /v1/transactions/{txHash}
- GET /v1/fees/{recipient}?currency=
- GET /health/live
- GET /health/ready
- POST /v1/launches/prepare
- POST /v1/swaps/quote
- POST /v1/swaps/prepare
- POST /v1/claims/fees/prepare

- [x] **步骤 1：从 Foundry Artifact 生成 ABI Binding**

~~~bash
cd contracts && forge build
cd ../server && go generate ./common/chain
~~~

生成文件提交到仓库，生产构建不依赖 Foundry。

- [x] **步骤 2：写事件解析与幂等失败测试**

使用真实编码的 Launched Log，精确验证 Token、Pool ID、launchpadId、Creator、Quote、Supply、Block、Tx、LogIndex 和 Factory；重复播放只产生一条记录。

- [x] **步骤 3：写重组失败测试**

替换已保存区块的 Hash 和 Log，验证先删除孤块记录再播放 canonical Log。

- [x] **步骤 4：实现有界轮询**

索引到 latest-2，单批最多 500 个区块。Hash 不匹配时最多回退 64 个区块；找不到共同祖先则停止并报错。

- [x] **步骤 5：在 RPC Logic 实现查询与交易准备**

Cursor 使用 created_at+id，默认 25、最大 100；大整数使用十进制字符串；Fee 在同一个 RPC Block 直接读取 Escrow。PrepareLaunch、PrepareSwap 和 PrepareFeeClaim 只返回 chainId、from、to、data、value、deadline 和 review，不签名、不广播。

- [x] **步骤 6：实现 Gateway 路由转换**

Gateway 将 RPC 错误映射为 problem JSON，并保持 o1 兼容的 HTTP 路径和字段命名。Gateway 不读取 PostgreSQL、Base RPC 或部署清单。

- [x] **步骤 7：实现 Indexer 启动与依赖装配**

Indexer 启动时使用 conf.MustLoad(*configFile,&c,conf.UseEnv())，调用 c.MustSetUp()，从 ServiceContext 获取 Model/RPC/Deployment，不直接读取环境变量或硬编码地址。

- [x] **步骤 8：运行并提交**

执行：cd server && go test ./... -count=1

提交信息：feat: index and expose trusted launches

---

### 任务 11：前端基础与类型化客户端

**文件：**
- 新建：web Next.js App Router 工程
- 新建：web/lib/api.ts、contracts.ts、launchpad-id.ts、receipt.ts
- 新建：web/components/providers.tsx、connect-wallet.tsx
- 测试：web/lib/launchpad-id.test.ts、receipt.test.ts

**接口：**
- normalizeSlug、deriveLaunchpadId、parseLaunchedReceipt
- 自动携带 Cookie 的本地 API Client
- wagmi/viem 只配置 Base Sepolia

- [x] **步骤 1：只安装必要依赖**

Next、wagmi、viem、TanStack Query、Vitest、Testing Library 和 Playwright。不引入 UI 框架。

- [x] **步骤 2：写 slug 与 Receipt 失败测试**

只有来自配置 Factory 的 Launched Event 可以被解析为成功结果。

- [x] **步骤 3：实现 Chain Provider 与生成 ABI 导出**

连接链不是 84532 时禁用写操作，并提供 switch-network。

- [x] **步骤 4：实现 API Client**

Fetch 使用 credentials: include；按 status/code 解析 problem response；不得引用 o1 API。

- [x] **步骤 5：运行并提交**

执行：cd web && npm test -- --run && npm run build

提交信息：feat: establish wallet and API frontend

---

### 任务 12：创建 Launchpad 用户流程

**文件：**
- 新建：web/app/page.tsx
- 新建：web/app/create-launchpad/page.tsx
- 新建：web/components/launchpad-form.tsx、transaction-status.tsx
- 测试：web/app/create-launchpad/page.test.tsx

**接口：**
- 钱包登录 → 调用 Registry → 提交品牌和 TxHash → 跳转 /launchpad/{slug}

- [x] **步骤 1：写 UI 失败测试**

覆盖未连接钱包、错误链、表单错误、拒绝签名、Registry 回滚、API 验证延迟、成功跳转和失败时保留输入。

- [x] **步骤 2：实现 EIP-191 Challenge 登录**

- [x] **步骤 3：实现 Registry 交易和 API 验证**

MVP Treasury 默认使用连接钱包。分别展示 waiting for wallet、confirming onchain、saving branding。

- [x] **步骤 4：验证可访问性**

所有输入有 Label；错误通过 aria-describedby 关联；焦点移动到首个错误项；交易状态使用 aria-live。

- [x] **步骤 5：运行并提交**

执行：cd web && npm test -- --run app/create-launchpad/page.test.tsx

提交信息：feat: create branded launchpads

---

### 任务 13：Tenant 页面与 Token 发行流程

**文件：**
- 新建：web/app/launchpad/[slug]/page.tsx
- 新建：web/app/launchpad/[slug]/create-token/page.tsx
- 新建：web/components/token-launch-form.tsx、token-card.tsx
- 测试：web/app/launchpad/[slug]/create-token/page.test.tsx

**接口：**
- 读取本地 Launchpad/Token API 与 Registry。
- 调用 POST /v1/launches/prepare 获取包含新 configVersion 和最新区块时间+30 分钟 deadline 的未签名 Factory 交易。
- Receipt 立即展示；随后轮询本地交易状态。

- [x] **步骤 1：写 Tenant 渲染和未知 slug 测试**

- [x] **步骤 2：写发行确认页测试**

签名前必须展示固定供应量、永久流动性、1% 协议费、0.5% LaaS、1.5% 总费率、16 秒 anti-snipe、Quote、launchpadId 和 Treasury。

- [x] **步骤 3：实现 Prepare + 钱包广播**

生成随机 bytes32 salt，调用 /v1/launches/prepare，展示服务器返回的 Review，然后将 to、data、value 原样交给钱包。遇到 stale_plan、StaleConfig 或 LaunchExpired 时重新 Prepare，并要求用户重新审阅和签名。

- [x] **步骤 4：实现 Receipt-first 对账**

只解析配置 Factory；立即显示 confirming；使用有上限的指数退避轮询，直到 confirmed 或 reverted。

- [x] **步骤 5：运行并提交**

执行：cd web && npm test -- --run && npm run build

提交信息：feat: launch and browse tenant tokens

---

### 任务 14：Token 页面、Claim 与测试 Swap

**文件：**
- 新建：web/app/token/[address]/page.tsx
- 新建：web/components/fee-balances.tsx、swap-form.tsx
- 测试：web/app/token/[address]/page.test.tsx

**接口：**
- 读取本地 Token 数据和 Escrow 余额。
- 通过 /v1/swaps/quote 与 /v1/swaps/prepare 获取 Quote 和未签名 Universal Router 交易。
- 通过 /v1/claims/fees/prepare 获取未签名 Escrow Claim 交易。
- 仅提供 exact-input 买卖。

- [x] **步骤 1：写费用展示与 Claim 失败测试**

Creator、Protocol、Referrer、LaaS 分开展示；余额为 0 时禁用；Claim Receipt 成功后刷新。

- [x] **步骤 2：写 Swap 安全失败测试**

拒绝零金额、余额不足、Quote 过期和错误链；签名前展示 Price Impact、Slippage、Fee、anti-snipe、Deadline；不提供 exact-output。

- [x] **步骤 3：实现 o1 风格 Quote/Prepare/广播路径**

前端先调用本地 /v1/swaps/quote，再把 quote_id 交给 /v1/swaps/prepare，最后将返回的 to、data、value 原样交给钱包。RPC 使用准确 PoolKey、官方 Base Sepolia Quoter/Router/Permit2 构造并模拟交易；本地服务不签名、不广播，也不调用 o1 生产 API。

- [x] **步骤 4：运行并提交**

执行：cd web && npm test -- --run && npm run build

提交信息：feat: display claim and test launch fees

---

### 任务 15：端到端验证与运行手册

**文件：**
- 新建：web/e2e/launchpad.spec.ts
- 新建：scripts/demo-smoke.sh
- 新建：docs/demo/base-sepolia-runbook.md
- 修改：Makefile、README.md

**接口与产物：**
- 提供 make test、make dev、make smoke。
- 输出可通过区块浏览器复核的完整 Demo 证据。

- [x] **步骤 1：增加 Playwright 本地 Happy Path**

完成钱包认证、创建 Launchpad、保存品牌、发行 Token、Receipt-first 展示、事件索引和正确 Tenant 列表展示。

- [x] **步骤 2：增加真实 Base Sepolia Smoke**

预检 Chain/Deployment，创建 Launchpad、发行 Token、执行 exact-input Swap、验证四类费用、Claim 一笔余额、查询 API 归属。任何地址或 Chain 不匹配立即停止。

- [x] **步骤 3：编写部署与恢复手册**

包含 Migration、部署顺序、Hook Salt 校验、API/Indexer 启动、前端部署、浏览器链接、Indexer 回退和密钥轮换。

- [x] **步骤 4：执行完整本地验证**

执行：make test

预期：Foundry Unit/Fuzz、Go、前端测试与构建、Playwright 全部通过。

- [ ] **步骤 5：执行 Base Sepolia Smoke**

执行：make smoke

记录 Registry Tx、Launch Tx、Pool ID、Swap Tx、费用合计、Claim Tx 和 API 中正确的 launchpadId 绑定。

- [x] **步骤 6：提交**

提交信息：test: verify Base Sepolia MVP end to end

## 验收清单

- 钱包可以创建并访问 /launchpad/{slug}。
- 另一个钱包可以在该 Launchpad 下发行 Token。
- 受信 Factory 事件包含准确 launchpadId。
- Token 只出现在其权威归属的 Tenant 下。
- Pool 正常费率固定为 1%+0.5%，LaaS 不参与 Creator/Platform/Referrer 分配。
- Anti-snipe 在 16 秒后精确达到 1.5%。
- Escrow Credit/Claim 与实际配对资产收费完全对账。
- Hook 永久锁定发行流动性。
- 等待 2 个确认的 Indexer 可幂等重放并处理短重组。
- 浏览器、仓库、日志和数据库不包含私钥或 o1 API Key。
- 使用公开 Base Sepolia 地址和运行手册可以复现完整 Demo。
