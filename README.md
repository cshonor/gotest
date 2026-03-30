# Go 完整测试体系示例（新手版）

这个项目是给新手小白的：用一个最小的 HTTP API 项目，演示 Go 项目里常见的 **5 类测试目录**，以及每类测试用什么命令跑。

你会学到：

- **Unit（单元测试）**：只测业务逻辑，不连真实数据库/网络
- **Integration（集成测试）**：连真实依赖（这里用 SQLite），验证“模块和依赖是否能一起工作”
- **E2E（端到端/接口测试）**：把 HTTP 服务完整跑起来，按真实请求方式黑盒验证接口
- **Benchmark（基准测试）**：测性能（ns/op、B/op、allocs/op），用于性能回归对比
- **Full（全量/冒烟）**：Go 里全量通常是 `go test ./...`；这里额外放一个 smoke test 示例帮助理解“全链路最小验证”

## 项目结构

```text
.
├── cmd/api/main.go
├── internal/
│   ├── app/app.go
│   ├── handler/user_handler.go
│   ├── model/user.go
│   ├── repository/
│   │   ├── repository.go
│   │   ├── sqlite_user_repository.go
│   └── service/
│       └── user_service.go
└── tests/
    ├── benchmark/benchmark_test.go
    ├── benchmark/sub_benchmark_test.go
    ├── e2e/api_test.go
    ├── e2e/health_test.go
    ├── full/smoke_test.go
    ├── integration/health_migration_test.go
    ├── integration/sqlite_user_repository_test.go
    ├── unit/table_driven_test.go
    └── unit/user_service_test.go
```

## 本地运行

```powershell
go run ./cmd/api
```

接口：

- `GET /health`
- `POST /api/v1/users` body: `{"name":"alice"}`
- `GET /api/v1/users`

## 执行测试

### 0) 先记住最常用的一条（全量）

在 Go 里，“全量测试”最常见的意思就是：把仓库里所有 package 的测试都跑一遍。

```powershell
go test ./... -v
```

### 1) Unit（单元测试）

**测什么**

- 只测“函数/方法的业务规则”
- 不依赖真实数据库、网络
- 失败时更容易定位到具体规则

**怎么跑**

```powershell
go test ./tests/unit -v
```

### 2) Integration（集成测试）

**测什么**

- 连接真实依赖（这里是 SQLite 文件数据库）
- 重点验证：建表/写入/查询这些“和依赖交互”的逻辑是否正确

**怎么跑**

```powershell
go test ./tests/integration -v
```

### 3) E2E（端到端 / 接口测试）

**测什么**

- 用 `httptest.NewServer` 起一个完整 HTTP 服务
- 用 `http.Get/http.Post` 像真正客户端一样请求接口（黑盒）

**怎么跑**

```powershell
go test ./tests/e2e -v
```

### 4) Full（全量/冒烟示例）

**它是什么**

- 真正“全量”通常是 `go test ./...`
- 这个目录放一个 **smoke test**：跑起路由 + 打一下 `/health`，帮助你理解“全链路最小验证”的写法

**怎么跑**

```powershell
go test ./tests/full -v
```

### 5) Benchmark（基准测试 / 性能测试）

**测什么**

- 你会看到类似：`ns/op`（耗时）、`B/op`（每次操作分配字节）、`allocs/op`（每次操作分配次数）
- 用来做“性能基线”和“回归对比”（改代码前后对比）

**怎么跑**

```powershell
go test ./tests/benchmark -bench=. -benchmem -run=^$
```

可重复多轮稳定性观察（示例）：

```powershell
go test ./tests/benchmark -bench=. -benchmem -run=^$ -count=5
```

## 说明

- 单元测试目录 `tests/unit` 里用 fake repo 来隔离外部依赖（你只需要关心业务规则是否正确）。
- 集成测试目录 `tests/integration` 会用 `t.TempDir()` 创建临时 SQLite 文件，测试结束自动清理。
- E2E 目录 `tests/e2e` 是“接口黑盒测试”的一个常见实现方式（不用真的起 8080 端口）。
- Benchmark 目录 `tests/benchmark` 里包含 service/repo/http 关键路径的基准测试示例。
