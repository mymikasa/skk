# Go 项目结构参考

## 标准目录布局


```
myapp/
├── cmd/                        # 可执行程序入口
│   ├── server/
│   │   └── main.go             # HTTP 服务入口
│   └── worker/
│       └── main.go             # 后台任务入口
│
├── internal/                   # 私有代码（不能被外部 import）
│   ├── config/                 # 配置加载与校验
│   │   └── config.go
│   ├── handler/                # HTTP handler（薄层，只做参数绑定和响应）
│   │   └── user.go
│   ├── service/                # 业务逻辑层
│   │   └── user.go
│   ├── repository/             # 数据访问层（接口 + 实现）
│   │   ├── user.go             # UserRepository interface 定义
│   │   └── postgres/
│   │       └── user.go         # PostgreSQL DAO 实现（SQL 查询、事务管理）
│   └── domain/                 # 核心领域模型（无依赖）
│       └── user.go
│
├── pkg/                        # 可对外暴露的公共包（谨慎放这里）
│   └── pagination/
│
├── api/                        # API 定义文件
│   ├── openapi.yaml
│   └── proto/
│
├── migrations/                 # 数据库迁移
├── scripts/                    # 构建/部署脚本
├── testdata/                   # 测试固定数据
│
├── go.mod
├── go.sum
├── Makefile
└── README.md
```
```

---

## 分层依赖规则

```
handler → service → repository → domain
   ↓         ↓           ↓
 (薄层)   (业务逻辑)  (数据访问)  (纯模型，无依赖)
```

**关键原则**：
- `domain` 包不能 import 任何内部包
- `repository` 层定义 interface，在 `service` 层接受注入
- `handler` 层不包含业务逻辑，只做请求解析和响应序列化
- `cmd/` 的 `main.go` 负责组装所有依赖（依赖注入根节点）

---

## 包命名规则

```go
// ✅ 小写，简短，无下划线，无复数
package user
package config
package postgres

// ❌ 避免
package userService
package user_repo
package utils    // 太泛，考虑拆分为具体功能包
package common   // 同上
package helpers  // 同上
```

---

## 小型项目（单服务）

不必强制分层，可以扁平结构：

```
myapp/
├── main.go
├── handler.go
├── service.go
├── storage.go
├── config.go
└── go.mod
```

当文件超过 5-8 个或单文件超过 300 行时，再考虑拆包。

---

## Makefile 常用目标

```makefile
.PHONY: build test lint run

build:
	go build -ldflags="-X main.version=$(shell git describe --tags)" \
	  -o bin/server ./cmd/server

test:
	go test -race -coverprofile=coverage.out ./...

lint:
	golangci-lint run ./...

run:
	go run ./cmd/server

tidy:
	go mod tidy
```