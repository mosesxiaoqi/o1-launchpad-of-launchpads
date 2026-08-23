# o1 Launchpad of Launchpads

基于 o1 Launchpad 合约逻辑构建的多租户 Launchpad MVP，目标网络为 Base Sepolia。

仓库已经包含可运行的 MVP：

- [方案 A 架构设计](docs/superpowers/specs/2026-08-23-launchpad-of-launchpads-scheme-a-design.md)
- [方案 A 实施计划](docs/superpowers/plans/2026-08-23-launchpad-of-launchpads-scheme-a.md)
- [Base Sepolia 部署、运行与恢复手册](docs/demo/base-sepolia-runbook.md)

系统由以下部分组成：

- Foundry 智能合约；
- Next.js 前端；
- go-zero REST Gateway；
- Proto/zRPC 业务服务；
- Base Sepolia 事件 Indexer；
- PostgreSQL。

## 本地准备

要求：Go、Foundry、Node.js/npm、PostgreSQL 16、`psql`。复制 `.env.example` 为本地 `.env`，填写现有 Docker PostgreSQL 的 `DATABASE_URL`、Base Sepolia RPC 和随机 `SESSION_SECRET`。真实 `.env` 已被 Git 忽略。

```bash
set -a
source .env
set +a
make migrate
cd web && npm install && npx playwright install chromium
```

`deployments/base-sepolia.json` 中的自有合约地址在实际部署前仍为 `null`；未完成部署时 API/RPC/Indexer 的链上路径会按预期拒绝启动或预检失败。

## 常用命令

```bash
make test   # Foundry、Go、Vitest、ESLint、Next build、Playwright
make dev    # RPC、REST API、Indexer、Next.js 四个开发进程
make smoke  # 对已运行环境执行真实 Base Sepolia 全链路验证
```

`make dev` 使用同一个配置项目启动四个进程：浏览器只访问 REST API；REST Gateway 通过 zRPC 调业务服务；RPC 服务与 Indexer 共享 PostgreSQL 和部署清单，但二者不互相直接通信。

真实链部署和 Smoke 所需变量、执行顺序、证据格式及恢复步骤见运行手册。
