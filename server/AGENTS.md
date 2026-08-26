# Go / go-zero 后端开发规范

本文件约束 `server/` 下后续所有 Go、go-zero API/RPC、Indexer、数据库和链上集成开发。当前代码用于识别技术栈和业务风险，不代表现有实现天然符合本规范。新增代码必须遵守本规范；修改既有代码时，应让被触及区域向本规范收敛。

## 技术栈与架构边界

- Go 版本与依赖以 `go.mod` 为唯一来源。
- 服务框架为 go-zero：HTTP API、zRPC、Logic、ServiceContext 和配置遵循 go-zero 约定。
- 链上交互使用 go-ethereum；金额使用整数、十进制字符串或 `big.Int`。
- 持久化使用 PostgreSQL 和 `database/sql`。
- 后端由三类运行单元组成：HTTP API、Launchpad RPC、链上 Indexer。
- `common/` 放跨服务复用且边界稳定的能力；不得变成无归属的工具箱。
- `service/launchpad/api` 负责 HTTP 协议、会话中间件和 RPC 转换。
- `service/launchpad/rpc` 负责业务编排、数据库与链上读取、无签名交易准备。
- `service/indexer` 负责确认区块后的日志抓取、重组恢复、幂等落库和 checkpoint。

## 工作原则

- 只做完整解决当前任务所需的最小改动，不顺便重构无关包。
- 修改前先追踪入口、Logic、ServiceContext、Model、链客户端和测试的完整调用链。
- 遵循已有包边界和命名；不要为单一实现创建接口、工厂或通用框架，测试替身或真实多实现需求除外。
- 优先使用标准库、go-zero、go-ethereum 和现有依赖；未经明确需求不得新增或升级依赖。
- 每个行为变更和缺陷修复必须有测试；缺陷修复必须包含回归测试。
- 不得删除测试、放宽断言、吞掉错误或降低安全检查来让代码“通过”。
- 新代码不能以现有不理想实现为先例；必要时在任务范围内做最小收敛，并说明剩余风险。

## go-zero 分层职责

### HTTP Handler

- 只负责协议解析、调用 Logic 和输出统一响应。
- 不得直接访问数据库、链节点、合约 binding 或实现业务规则。
- 参数校验以 `.api` 声明和 Logic 的领域校验共同完成，不能只信任自动绑定。
- 身份信息必须从可信中间件 context 获取，不能信任请求体中可伪造的钱包字段。

### API Logic

- 负责 HTTP 语义、身份授权、输入规范化及对 zRPC 的调用。
- 不复制 RPC 中的核心业务规则，不直接拼装链上 calldata。
- 将 RPC/领域错误稳定映射为公开 HTTP problem；不得泄露 SQL、堆栈、RPC URL 或内部实现。

### RPC Logic

- 负责业务用例编排、权限验证、数据库查询和链上读取。
- 使用 `l.ctx`/传入 context 贯穿数据库、RPC 和下游调用。
- 交易准备服务只返回待用户审查和签名的交易，不持有用户私钥，不代替用户广播。
- 复杂规则应放入可单测的领域函数或 `common/` 中已有的明确模块，而不是堆进 Handler。

### ServiceContext

- 只构造和持有进程级依赖：配置、数据库、链客户端、Model、RPC client 和有界内存存储。
- 新增依赖必须明确所有权、生命周期、关闭方式和并发安全。
- 启动失败应 fail fast；请求路径和可恢复后台任务不得使用 `panic`。

### Model

- 负责 SQL 和事务，不包含 HTTP/RPC 展示逻辑。
- 多个必须原子完成的写操作必须使用同一事务。
- Model 返回可识别的领域错误或包装后的底层错误；调用方不得匹配错误字符串。

## 生成代码与唯一来源

以下文件是生成产物，不得手工编辑：

- `service/launchpad/api/internal/types/types.go`；
- `service/launchpad/api/internal/handler/routes.go` 及 goctl 标记的 Handler 骨架；
- `service/launchpad/rpc/pb/`；
- `service/launchpad/rpc/launchpadclient/`；
- `service/launchpad/rpc/internal/server/`；
- `common/chain/bindings/`。

修改 API 时，以 `service/launchpad/api/launchpad.api` 为唯一来源，然后在该目录运行：

```bash
goctl api go -api launchpad.api -dir .
```

修改 RPC 时，以 `service/launchpad/rpc/launchpad.proto` 为唯一来源，然后在该目录运行：

```bash
goctl rpc protoc launchpad.proto --go_out=. --go-grpc_out=. --zrpc_out=.
```

修改合约 ABI 时，先从仓库根目录执行：

```bash
cd contracts && forge build
cd ../server && go generate ./common/chain
```

生成后必须审查 diff，确认没有覆盖手写业务逻辑、意外删除接口或引入无关格式变化。生成器版本变化属于工具链变更，必须明确说明并单独审查。

## Go 编码规范

- 所有修改的 Go 文件必须经过 `gofmt`；import 由 Go 工具管理。
- 使用清晰、惯用、直接的 Go；避免反射、全局可变状态和无必要泛型。
- 不要创建只有一次调用且不能提升可读性的小包装。
- 函数应保持单一职责；错误路径优先返回，减少深层嵌套。
- 接口应定义在消费方，并保持最小方法集；只有真实替换或测试隔离需求才引入。
- 不得提交调试输出、被注释掉的旧代码、无归属 TODO 或生成二进制。
- `go.mod`/`go.sum` 只通过 Go 工具修改；依赖未变化时不要无意义运行 `go mod tidy`。

## 错误处理

- 除进程启动阶段明确不可恢复的配置/依赖错误外，不得 `panic`。
- 不得忽略错误。确实只能 best-effort 的清理操作，应记录或注释说明为何可以忽略。
- 包装错误时加入操作和关键标识，并使用 `%w` 保留原因：

```go
return fmt.Errorf("load launchpad %q: %w", slug, err)
```

- 使用 `errors.Is`/`errors.As` 判断错误，不比较错误字符串。
- 错误只在能够增加上下文或决定处理策略的层记录；避免每层记录一次造成重复日志。
- 对外错误必须稳定、有限且不泄露内部信息；内部错误保留足够诊断上下文。
- `sql.ErrNoRows`、唯一键冲突、认证失败、参数错误、上游不可用和 context 取消必须区分处理。

## Context、超时与资源生命周期

- `context.Context` 作为函数第一个参数传递，不存入长期结构体。
- HTTP、gRPC、数据库、以太坊 RPC 和后台批处理必须传播取消与 deadline。
- 不得用 `context.Background()` 绕过已有请求 context；只有进程启动或明确脱离请求的生命周期可使用。
- 所有网络和数据库操作必须有上游 deadline 或显式超时。
- 循环、重试和批处理中检查 `ctx.Done()`，取消后及时退出。
- `Rows`、HTTP response body、ticker、timer、subscription、数据库和 RPC client 必须由创建者明确关闭。
- 禁止无退出机制、无错误通道或无所有者的 goroutine。

## 并发安全

- 共享 map、内存 quote store、nonce 状态、checkpoint 和客户端状态必须同步保护或保持不可变。
- 持锁期间不得进行数据库、RPC、网络或其他不可控耗时操作。
- channel 由发送方/所有者关闭，接收方不得随意关闭。
- goroutine 必须可通过 context、stop channel 或服务生命周期终止。
- 后台任务的 panic 必须被视为进程级故障或被显式恢复并上报；不能静默死亡。
- 修改并发、缓存、Indexer 或后台循环时，必须运行 race detector。

## HTTP、gRPC 与接口兼容性

- 所有外部输入都在边界处校验：长度、格式、枚举、分页上限、地址、哈希、金额、chain ID 和业务关系。
- API 列表必须分页或有明确硬上限；不得允许用户触发无界查询或无界链上扫描。
- HTTP 状态码、problem code、JSON 字段、gRPC method 和 protobuf 字段视为公开契约。
- protobuf 字段编号不得复用；删除字段时使用 `reserved` 保留编号和名称。
- 新字段默认应向后兼容；改变含义不能只保留旧名称。
- `.api` 与 `.proto` 表示同一业务对象时，字段语义、单位、可选性和错误行为必须一致。
- chain ID 不能从不可信请求覆盖服务配置；多链支持必须显式设计，而非临时放宽校验。

## 身份认证与会话

- 钱包认证 challenge 必须绑定预期 domain、URI、chain ID、地址、nonce、签发时间和过期时间。
- nonce 必须一次性、随机、短期有效，并在验证成功时原子消费。
- 签名恢复后使用规范化地址比较；正确处理 `v` 值和签名长度，不接受畸形编码。
- 防止 challenge 在不同域、链、地址或会话之间重放。
- Session 必须使用强随机密钥和恒定时间 MAC 比较，验证版本、过期时间和 subject。
- 不得记录 challenge 完整签名、session token、HMAC secret、私钥或授权头。
- Handler 不得信任请求体中的 wallet 作为当前身份；必须与认证 context 对照。

## 数据库与迁移

- 所有 SQL 使用参数绑定；不得拼接用户输入。动态排序/字段只能来自固定白名单。
- `QueryContext`、`ExecContext`、`BeginTx` 必须使用传入 context。
- `rows.Close()` 与 `rows.Err()` 必须处理。
- 事务闭包内任一步失败都必须回滚；提交错误必须返回。
- 链事件、派生业务记录、交易状态和 checkpoint 需要一致时，必须在同一事务中提交。
- 唯一约束和外键用于维护不可破坏的不变量，不要只依赖应用层预检查。
- 金额使用 PostgreSQL `numeric`/文本或无损整数表达，不使用 `float`。
- 查询必须考虑索引、锁、数据规模和分页；不得在请求路径做全表扫描。
- 已在环境执行的迁移不得原地修改。使用新的递增迁移文件演进 schema。
- 迁移必须考虑滚动发布兼容、已有数据回填、锁表时间、失败恢复和降级策略。
- 未经明确授权，不得对真实数据库执行迁移、回填或删除。

## 链上读取与交易准备

- 地址必须使用 go-ethereum 类型解析并严格校验；不要把 `HexToAddress` 对无效字符串的零值结果当作成功。
- 交易哈希、block hash、PoolId、bytes32 和 calldata 必须验证长度和编码。
- 金额、供应量、费用和余额不得使用浮点数；跨边界使用十进制字符串，内部使用整数或新建的 `big.Int`。
- `big.Int` 是可变对象；不得共享后修改，也不得把调用方指针直接保存到长期状态。
- `bind.CallOpts` 必须带 Context；需要一致视图时绑定明确 block number。
- RPC 返回的 receipt、status、block number、log 地址、topic 和事件参数必须全部验证。
- 交易确认必须检查成功状态和所需确认数；“RPC 已看到”不等于最终确认。
- 服务只准备无签名交易时，必须返回明确的 chain ID、from、to、data、value、deadline 和人类可读 review。
- 后端不得持有用户私钥、替用户签名或自动广播，除非未来任务明确设计托管签名系统并完成安全评审。
- 不得仅依赖前端展示保护高价值交易；后端必须自行验证配置版本、地址和业务约束。

## Indexer、确认数与链重组

- Indexer 必须可重启、可重试、幂等，重复处理同一日志不能产生重复业务记录。
- 只处理达到配置确认数的区块；确认数是安全参数，不得为提速随意降低。
- checkpoint 至少绑定 chain ID、合约地址、next block 和 last block hash。
- 每次继续扫描前验证父区块/上次区块 hash；不匹配时按配置 lookback 回滚并重新索引。
- 回滚必须删除或修正受影响区块派生的数据和交易状态，不能只移动 checkpoint。
- 事件落库与 checkpoint 前进必须原子；否则崩溃可能造成漏事件。
- 使用 `(chain_id, contract, tx_hash, log_index)` 或等价稳定键去重。
- 必须处理 removed logs、空区间、RPC 限流、批次过大、节点落后和暂时缺块。
- 批次大小、轮询间隔和 reorg lookback 必须有合理上限并来自经过校验的配置。
- ABI、事件 topic 或部署地址变化时，Indexer 必须同步升级并评估历史重放。

## 配置、密钥与日志

- 启动时验证所有必要配置：数据库、RPC、chain ID、confirmations、部署文件、合约地址、session secret 和批处理参数。
- 配置错误应 fail fast；不得用不安全默认值继续运行。
- 密钥、数据库 DSN、RPC 凭证、session secret 和令牌不得进入源码、测试快照或日志。
- 日志使用结构化字段，包含必要的 request ID、chain ID、block、tx hash、服务和操作。
- 不要记录完整请求体、签名、session 或敏感 calldata。
- 指标标签不得使用钱包、交易哈希、slug 等高基数值，除非监控系统明确允许。
- 后台任务失败必须可观察；不得只打印后继续无限快速重试。

## 测试规范

- 测试优先使用表驱动和小型 fake 接口，不启动不必要的完整服务。
- 测试必须确定性，不依赖真实时间、随机公网 RPC 或测试执行顺序。
- 时间、RPC、数据库和外部依赖应通过最小接口或已有注入点控制。
- 缺陷修复的回归测试在没有修复时应失败。
- SQL/事务测试覆盖提交、回滚、唯一键冲突、无记录和 context 取消。
- API/RPC 测试覆盖校验、认证、错误映射、兼容字段和分页。
- 链上测试覆盖无效地址/哈希、失败 receipt、确认数、事件不匹配和 RPC 错误。
- Indexer 测试覆盖首次同步、空批次、重复日志、重启、重组、checkpoint 原子性和取消退出。
- 并发改动覆盖竞争、重复消费、Stop 幂等和 goroutine 退出。

开发时运行最小相关包测试：

```bash
cd server
go test ./common/chain ./common/model -count=1
```

完成前至少运行：

```bash
cd server
gofmt -w <修改的 Go 文件>
go vet ./...
go test ./... -count=1
```

涉及并发、Indexer、缓存或后台任务时再运行：

```bash
go test -race ./... -count=1
```

从仓库根目录可运行后端门禁：

```bash
make test-server
```

不得声称未实际执行的命令已经通过。受环境限制无法运行数据库、RPC 或 race 测试时，必须说明原因和剩余风险。

## Code Review 规则

审查 Go 后端差异时，至少回答：

1. 改动属于 Handler、Logic、ServiceContext、Model、链客户端还是 Indexer？职责是否放对层？
2. 是否手工修改了生成文件，或遗漏 `.api`、`.proto`、ABI 的重新生成？
3. context、timeout、资源关闭和 goroutine 退出是否完整？
4. 错误是否保留原因、正确映射且不泄露内部信息？
5. 数据库操作是否需要事务，事件与 checkpoint 是否原子？
6. 是否能抵抗重复请求、重复日志、进程重启和链重组？
7. 地址、哈希、chain ID、金额、签名和身份是否在可信边界验证？
8. 是否使用浮点金额、错误共享 `big.Int` 或依赖未确认交易？
9. 是否改变公开 API、protobuf、数据库 schema、事件解析或配置兼容性？
10. 测试是否覆盖失败、边界和恢复路径，而不只是 happy path？

高风险问题未解决、测试仍失败或 CI 仍 pending 时不得描述为完成。

## 完成条件

- 改动范围与任务一致，无无关重构或未授权依赖升级。
- 生成文件均由对应源文件和工具生成，生成 diff 已审查。
- 修改文件已 `gofmt`，`go vet ./...` 与 `go test ./... -count=1` 已通过。
- 并发敏感改动已通过 `go test -race`，或明确说明无法执行的原因。
- API、RPC、数据库、ABI、配置和索引兼容性影响已检查。
- 链重组、确认数、幂等、事务和认证风险已按改动范围测试。
- 未提交密钥、凭证、生成二进制或调试输出。
- CI 必须到达最终状态；pending、跳过或 flaky 不能描述为通过。
