# UserService Design

## Overview

使用 Gin + SQLite 实现用户管理服务，包含用户的增删改查操作。采用经典三层架构（Handler → Service → Repository），Repository 通过 interface 抽象以方便后续替换存储实现。

## Project Structure

```
skk/
├── cmd/
│   └── server/
│       └── main.go              # 依赖注入根节点，组装并启动服务
├── internal/
│   ├── domain/
│   │   └── user.go              # User 模型定义（纯 struct，无依赖）
│   ├── handler/
│   │   └── user.go              # Gin HTTP handler（参数绑定 + 响应序列化）
│   ├── service/
│   │   └── user.go              # 业务逻辑层，接受 Repository 接口注入
│   ├── repository/
│   │   ├── user.go              # UserRepository interface 定义
│   │   └── sqlite/
│   │       └── user.go          # SQLite 实现
│   └── config/
│       └── config.go            # 配置加载（DB 路径、端口等）
├── go.mod
├── go.sum
└── Makefile
```

依赖方向：`handler → service → repository → domain`

- `domain` 包不能 import 任何内部包
- `repository` 层定义 interface，在 `service` 层接受注入
- `handler` 层不包含业务逻辑，只做请求解析和响应序列化
- `cmd/server/main.go` 负责组装所有依赖（依赖注入根节点）

## Domain Model

```go
type User struct {
    ID        int64     `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Age       int       `json:"age"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

## Repository Layer

### Interface (`repository/user.go`)

```go
type UserRepository interface {
    Create(ctx context.Context, user *domain.User) (*domain.User, error)
    GetByID(ctx context.Context, id int64) (*domain.User, error)
    Update(ctx context.Context, user *domain.User) (*domain.User, error)
    Delete(ctx context.Context, id int64) error
    List(ctx context.Context, offset, limit int) ([]*domain.User, error)
}
```

### SQLite Implementation (`repository/sqlite/user.go`)

- 使用 `database/sql` + `modernc.org/sqlite`（纯 Go 实现，无 CGO 依赖）
- 建表语句在初始化时自动执行
- SQLite 表结构：

```sql
CREATE TABLE IF NOT EXISTS users (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    name       TEXT    NOT NULL,
    email      TEXT    NOT NULL UNIQUE,
    age        INTEGER NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);
```

## Service Layer

职责：
- 业务校验：Email 格式、Name 非空、Age 范围（0-150）
- 错误包装：使用 `fmt.Errorf("service: create user: %w", err)`
- 定义哨兵错误供 handler 层判断 HTTP 状态码

哨兵错误：

```go
var (
    ErrUserNotFound = errors.New("user not found")
    ErrUserExists   = errors.New("user already exists")
    ErrInvalidInput = errors.New("invalid input")
)
```

## Handler Layer (Gin)

| Method  | Path        | Description | Request Body      | Success Code |
|---------|-------------|-------------|-------------------|-------------|
| POST    | /users      | Create user | User JSON         | 201         |
| GET     | /users/:id  | Get by ID   | -                 | 200         |
| PUT     | /users/:id  | Update user | User JSON         | 200         |
| DELETE  | /users/:id  | Delete user | -                 | 204         |
| GET     | /users      | List users  | query: page, size | 200         |

响应格式统一为 JSON。错误响应结构：

```json
{
  "error": "error message"
}
```

分页参数：`page`（默认 1）、`size`（默认 20，最大 100）。

## Error Handling Strategy

- Repository 层：返回原始数据库错误或 `ErrUserNotFound`
- Service 层：包装错误并附加上下文，返回哨兵错误
- Handler 层：根据 service 返回的错误类型映射 HTTP 状态码
  - `ErrUserNotFound` → 404
  - `ErrUserExists` → 409
  - `ErrInvalidInput` → 400
  - 其他 → 500

## Testing Strategy

### Repository Tests (`repository/sqlite/user_test.go`)
- 使用临时 SQLite 内存数据库（`:memory:`）
- 测试 CRUD 操作的正确性
- 测试唯一约束（email 重复）

### Service Tests (`service/user_test.go`)
- Mock `UserRepository` 接口
- 测试业务校验逻辑（非法 name、email、age）
- 测试哨兵错误返回

### Handler Tests (`handler/user_test.go`)
- 使用 `httptest.NewRecorder()` + Gin test mode
- 测试 HTTP 状态码和响应体
- 测试参数绑定和错误响应

## Dependencies

- `github.com/gin-gonic/gin` — HTTP 框架
- `modernc.org/sqlite` — SQLite 驱动（纯 Go，无 CGO）

## Configuration

```go
type Config struct {
    DBPath string // SQLite 数据库文件路径，默认 "data/skk.db"
    Port   string // HTTP 服务端口，默认 ":8080"
}
```

配置通过环境变量读取，支持 `SKK_DB_PATH` 和 `SKK_PORT`。
