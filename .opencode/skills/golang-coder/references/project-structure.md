# 项目结构

> 来源：Go 标准库、Google Go Style Guide、Uber Go Style Guide

## 核心原则

1. **包按职责划分，每个包做一件事。** 围绕一个领域概念组织，而不是一个技术层。包名描述它提供什么功能，不是包含什么文件。
2. **`internal/` 限制公开范围，`pkg/` 暴露可复用代码。** 不打算被外部导入的代码放 `internal/`；确认需要被其他项目复用的才放 `pkg/`。
3. **依赖方向单向流动：`cmd/` → `internal/` → 标准库/第三方。** 内层不依赖外层，同级包之间避免循环依赖。

## 标准目录布局

| 目录 | 用途 | 必需 |
|------|------|------|
| `cmd/` | 入口 main 包，每个可执行文件一个子目录 | 是 |
| `internal/` | 私有应用代码，其他项目不可导入 | 是 |
| `pkg/` | 可被外部复用的公共代码 | 否 |
| `api/` | API 定义（protobuf、swagger、GraphQL schema） | 否 |
| `configs/` | 配置文件模板 | 否 |
| `docs/` | 设计文档 | 否 |
| `scripts/` | 构建/安装/分析脚本 | 否 |
| `test/` | 额外的集成测试数据 | 否 |

## 包拆分规则

### 按领域拆分（推荐）

```
✅ good — 按领域
internal/
├── user/
│   ├── handler.go
│   ├── service.go
│   └── repository.go
└── order/
    ├── handler.go
    ├── service.go
    └── repository.go
```

### 按技术层拆分（不推荐）

```
❌ bad — 按技术层
internal/
├── handlers/
│   ├── user_handler.go
│   └── order_handler.go
├── services/
│   ├── user_service.go
│   └── order_service.go
└── repositories/
    ├── user_repository.go
    └── order_repository.go
```

按技术层拆分导致修改一个功能要跨三个包，增加耦合和循环依赖风险。

### 避免泛化包名

```
✅ internal/httputil  internal/validation  internal/clock
❌ internal/utils     internal/common      internal/helpers
```

泛化包名会变成垃圾桶。将辅助函数放到使用它的包中，或创建命名具体的小包。

### 包名描述提供什么

```
✅ 包名描述提供的能力：encoding/json、crypto/sha256、net/http
❌ 包名描述包含的内容：models、types、constants、structs
```

## 依赖管理

### go.mod 最佳实践

```
✅ module github.com/example/myapp
   go 1.22

   require (
       github.com/gin-gonic/gin v1.9.1
       github.com/lib/pq v1.10.9
   )
```

```
❌ module myapp          # 没有用完整路径
   require (
       github.com/gin-gonic/gin v1.9.1 // indirect  # 实际直接使用但标记 indirect
   )
```

### replace 指令

仅用于本地开发联调，不要在生产行使用：

```go
// ✅ 本地开发临时替换
replace github.com/example/mylib => ../mylib
```

```go
// ❌ 生产 go.mod 中用 fork 的 replace 绕过版本问题
replace github.com/some/pkg => github.com/fork/pkg v1.0.0
```

### 版本选择策略

- 优先选择有 semver 标签的依赖
- 避免 `go.mod` 中长期使用 `replace` 指令
- 定期运行 `go mod tidy` 清理未使用的依赖
- 用 `go list -m -json all` 审查依赖树

## 完整项目示例

### ✅ 推荐结构

```
myapp/
├── cmd/
│   └── myapp/
│       └── main.go
├── internal/
│   ├── user/
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   └── model.go
│   ├── order/
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   └── model.go
│   └── platform/
│       ├── database/
│       │   └── postgres.go
│       └── server/
│           └── server.go
├── api/
│   └── proto/
│       └── user.proto
├── configs/
│   └── config.yaml
├── go.mod
├── go.sum
└── Makefile
```

`internal/platform/` 存放横切关注点（数据库连接、HTTP 服务器配置），不属于特定业务领域。

### ❌ 反面示例

```
myapp/
├── main.go                          # main 放根目录
├── controllers/                     # 按技术层拆分，且不在 internal/
│   ├── user_controller.go
│   └── order_controller.go
├── models/                          # 泛化包名
├── services/
├── utils/                           # 垃圾桶包
│   ├── string_utils.go
│   └── http_common.go
├── pkg/                             # 空的 pkg/，项目是应用不是库
├── common/                          # 泛化包名
└── go.mod
```

问题：main 在根目录、按技术层拆分、泛化包名、没有 internal/ 保护、空的 pkg/。

## 反模式

| 反模式 | 问题 | 修正 |
|--------|------|------|
| 上帝包 | 一个包做所有事，代码膨胀到数千行 | 按领域拆分为多个包 |
| 循环依赖 | A 导入 B，B 导入 A | 提取共享部分到新包或定义接口解耦 |
| 过深嵌套 | internal/app/handlers/v1/... | 不超过 3 层 |
| internal 下按 MVC 分层 | internal/controllers/, internal/services/ | 在 internal/ 下按领域组织 |
| utils/common 包 | 泛化包名无意义，变成垃圾桶 | 用具体名称或放入使用者包 |
| pkg 下放所有代码 | pkg 失去复用语义 | pkg 只放真正被外部复用的代码 |
