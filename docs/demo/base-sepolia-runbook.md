# Base Sepolia 部署、运行与恢复手册

本文用于把方案 A MVP 部署到 Base Sepolia，并产出可由 Basescan 和本地 API 交叉复核的演示证据。命令均从仓库根目录执行；私钥、RPC 密钥和 Session Secret 只放本地 `.env`，不得写入部署清单、日志或 Git。

## 1. 前置条件

- Foundry、Go、Node.js/npm、PostgreSQL 16、`psql`、`curl`、`jq`。
- 已有 Docker PostgreSQL 中创建独立数据库 `o1_launchpad_mvp`。
- 部署钱包持有足够的 Base Sepolia ETH。
- `.env.example` 已复制为 `.env` 并填入真实值。

加载本地环境：

```bash
set -a
source .env
set +a
```

首先确认 RPC 网络，任何不等于 `84532` 的结果都必须停止：

```bash
cast chain-id --rpc-url "$BASE_SEPOLIA_RPC_URL"
```

## 2. Migration

数据库不由本仓库 Docker Compose 管理，直接复用现有 PostgreSQL：

```bash
make migrate
psql "$DATABASE_URL" -c '\dt'
```

Migration 使用 `IF NOT EXISTS`，可以重复执行。部署前备份数据库；不要在生产式演示环境手工修改 `launchpad_id`、Token 或 Pool 归属。

## 3. 合约部署顺序

`DeployBaseSepolia.s.sol` 严格执行以下顺序：

1. 部署一次性 `HookCreate2Deployer`。
2. 挖出满足 Uniswap v4 permission bits 的 salt 并预测 Hook 地址。
3. 使用预测 Hook 地址部署 `FeeEscrow`，解除 Escrow/Hook 循环依赖。
4. 通过 CREATE2 部署 `LaunchHookV2`，脚本立即校验预测地址和权限位。
5. 部署 `LaunchpadRegistry` 与 `MultiTenantLaunchpadFactory`。
6. 一次性执行 `hook.setFactoryOnce(factory, escrow)`。

执行：

```bash
cd contracts
forge script script/DeployBaseSepolia.s.sol:DeployBaseSepolia \
  --rpc-url "$BASE_SEPOLIA_RPC_URL" --broadcast -vvvv
cd ..
```

记录输出中的 HookDeployer、Registry、Factory、Hook、FeeEscrow、salt 和第一笔部署区块。把它们与 Quote、LaaS Treasury、Start Tick 一起写入 `deployments/base-sepolia.json`，不得保留 `null`。

Hook 地址低 14 位必须满足 `0x2acc`。以下检查应输出 `hook permission bits OK`：

```bash
HOOK=$(jq -r '.hook' deployments/base-sepolia.json)
SUFFIX=${HOOK: -4}
(( (16#$SUFFIX & 0x3fff) == 0x2acc )) && echo 'hook permission bits OK'
```

继续核对链上 wiring：

```bash
cast call "$(jq -r '.hook' deployments/base-sepolia.json)" 'factory()(address)' --rpc-url "$BASE_SEPOLIA_RPC_URL"
cast call "$(jq -r '.hook' deployments/base-sepolia.json)" 'feeEscrow()(address)' --rpc-url "$BASE_SEPOLIA_RPC_URL"
cast call "$(jq -r '.feeEscrow' deployments/base-sepolia.json)" 'hook()(address)' --rpc-url "$BASE_SEPOLIA_RPC_URL"
cast call "$(jq -r '.factory' deployments/base-sepolia.json)" 'registry()(address)' --rpc-url "$BASE_SEPOLIA_RPC_URL"
```

四项必须逐字匹配部署清单。源码验证可在配置 `BASESCAN_API_KEY` 后对每个已部署地址运行 `forge verify-contract`；构造参数必须来自同一次 broadcast 产物，不要手工猜测。

## 4. 服务启动

安装依赖并构建一次：

```bash
cd web
npm install
npx playwright install chromium
cd ..
make test
```

开发演示可直接执行：

```bash
make dev
```

这会启动四个进程：zRPC `:8080`、REST Gateway `:8888`、Indexer 和 Next.js `:3000`。启动次序建议为 PostgreSQL → zRPC → REST Gateway → Indexer → 前端。REST Gateway 与 Indexer 不直接通信：两者通过 zRPC/数据库和同一份链上事实最终汇合。

就绪检查：

```bash
curl -fsS http://127.0.0.1:8888/health/ready | jq
curl -fsS http://127.0.0.1:8888/v1/config | jq
```

前端生产构建使用 `NEXT_PUBLIC_API_URL` 和 `NEXT_PUBLIC_RPC_URL`：

```bash
cd web
npm run build
npm run start
```

## 5. 真实 Base Sepolia Smoke

当前自动 Swap 路径要求部署清单的 Quote 为原生 ETH 零地址。ERC-20 Quote 还需要补 Permit2 授权 UI，Smoke 会在广播前拒绝该配置。Smoke 必须访问启用 HTTPS 的 API；认证 Cookie 带 `Secure`，不能通过明文 `http://127.0.0.1:8888` 完成真实链流程。

不要把原始私钥放进进程参数。先导入 Foundry keystore，并把密码放入权限受限、Git 忽略的本地文件：

```bash
cast wallet import o1-smoke --interactive
chmod 600 /absolute/path/to/password-file
export SMOKE_KEYSTORE="$HOME/.foundry/keystores/o1-smoke"
export SMOKE_KEYSTORE_PASSWORD_FILE=/absolute/path/to/password-file
export SMOKE_API_URL=https://api-demo.example.com
export SMOKE_AUTH_ORIGIN=https://demo.example.com
```

`SMOKE_AUTH_ORIGIN` 必须与 Gateway 的 `FRONTEND_ORIGIN` 相同，代表用户实际访问的 HTTPS 前端 origin；服务端会忽略客户端自报的认证域，并始终使用这一配置作为签名信任根。`SMOKE_API_URL` 是 HTTPS REST Gateway origin。若前后端不同源，Gateway 前必须配置只允许该前端 origin 的 CORS、允许 credentials，并保留 `Set-Cookie`。更简单的部署方式是用同一 HTTPS 反向代理把 `/v1`、`/health` 转发到 Gateway，其余路径转发到 Next.js。

另准备两个彼此不同、且不同于部署钱包和 LaaS Treasury 的地址，以便四种费用不会在同一 Escrow recipient 下合并：

```bash
export SMOKE_LAUNCHPAD_TREASURY=0x...
export SMOKE_REFERRER_ADDRESS=0x...
export SMOKE_SWAP_AMOUNT_WEI=1000000000000000
make smoke
```

脚本按顺序执行 Chain/Deployment 预检、合约 bytecode 与 wiring 校验、EIP-191 登录、Registry 创建、品牌保存、Launch、等待两确认后的 Indexer 归属、16 秒 anti-snipe 窗口、exact-input Swap、四方非零费用核对和 Creator Claim。每笔 API prepare 交易都会核对 chain/from/to/value/deadline，先用 `eth_call` 模拟，再交给 keystore 签名广播。任何地址、链、归属或余额不匹配都会立即停止。

成功后证据写入被 Git 忽略的 `artifacts/base-sepolia-smoke-*.json`，包括 Registry Tx、Launch Tx、Pool ID、Token、launchpadId、Swap Tx、四方费用及合计/查询区块、Claim Tx 和 Basescan 链接。演示前逐项打开：

```text
https://sepolia.basescan.org/tx/<TX_HASH>
https://sepolia.basescan.org/address/<CONTRACT_ADDRESS>
```

## 6. Indexer 回退与恢复

Indexer 只读取配置 Factory 的 `Launched` 事件，等待 2 个确认，并保存 `(next_block, last_block_hash)`。短重组会自动比较上一块哈希；不一致时最多回退 `ReorgLookback=64` 个区块，删除孤块 Token、把对应交易标为 `orphaned`，再幂等重放。

故障处理顺序：

1. 停止 Indexer，保留 API/RPC 只读能力。
2. 核对 RPC 当前 Chain ID、部署清单 Factory 与 `deploymentBlock`。
3. 备份 PostgreSQL，查询 `indexer_checkpoints`、`token_launches`、`chain_transactions`。
4. 优先恢复正确 RPC 后重启，让内置哈希检查自动回退。
5. 超过 64 区块或清单错误时，从备份恢复；不要直接把 Token 改挂到另一个 Tenant。确需人工回放时，先在副本数据库验证，再把错误区间的 Token/Checkpoint 按同一事务回退至已确认区块。

API 中 Token 的 `launchpad_id` 必须与 Factory `Launched` 事件一致。数据库内容与链冲突时，以受信 Factory 的已确认事件为准。

## 7. 密钥轮换与泄露响应

- `SESSION_SECRET`：生成新随机值、重启 REST Gateway；现有 Session 全部失效，用户重新签名登录。
- RPC/WS Key：在供应商控制台轮换，更新本地/部署平台 Secret 后滚动重启 RPC 与 Indexer。
- Smoke/部署私钥：此架构部署后没有可升级管理员通道；把剩余测试 ETH 转到新钱包，停用旧 Key，并更新之后的 Smoke 钱包。不要把 Key 复制到证据文件。
- 若 Git 或日志出现私钥，立即视为已泄露并轮换；删除历史不是轮换的替代品。

最后执行：

```bash
make secrets-check
git status --short
```

确认 `.env`、broadcast、artifacts 和私钥均未进入待提交文件。
