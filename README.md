# Go 完整测试体系示例

这个项目演示了在 Go 项目里搭建完整测试体系：

- 单元测试（Unit Test）
- 集成测试（Integration Test）
- E2E 接口测试（端到端）
- 基准测试（Benchmark，性能/稳定性评估）

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
    ├── e2e/api_test.go
    ├── integration/sqlite_user_repository_test.go
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

### 1) 单元测试（Unit Test）

```powershell
go test ./tests/unit -v
```

### 2) 集成测试（Integration Test）

```powershell
go test ./tests/integration -v
```

### 3) E2E 接口测试（端到端）

```powershell
go test ./tests/e2e -v
```

### 4) 全量测试

```powershell
go test ./... -v
```

### 5) 基准测试（Benchmark）

```powershell
go test ./tests/benchmark -bench=. -benchmem -run=^$
```

可重复多轮稳定性观察（示例）：

```powershell
go test ./tests/benchmark -bench=. -benchmem -run=^$ -count=5
```

## 说明

- 单元测试重点验证业务规则（例如用户名长度校验），通过 fake repo 隔离外部依赖。
- 集成测试验证 repository 与真实 SQLite 的交互行为（建表、写入、查询）。
- E2E 测试通过 `httptest.NewServer` 启动完整 HTTP 服务，以黑盒方式验证接口行为。
- Benchmark 给出关键路径吞吐和内存分配数据，可用于性能基线与回归比较。
