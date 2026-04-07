---
name: golang-coder
description: 编写、审查、重构和调试 Go 语言代码。当用户需要用 Go 实现功能、构建服务、设计 API、优化性能、编写测试，或请求 Go 代码审查、提问 Go 语言特性、遇到 Go 相关编译/运行错误时，必须使用此 skill。适用场景包括但不限于：HTTP 服务、CLI 工具、并发程序、数据库操作、微服务、grpc、配置管理、日志系统。只要对话中出现"golang"、"go语言"、".go 文件"、"goroutine"、"go mod"等关键词，立即触发此 skill。
---

# Golang Coder Skill

本 skill 指导使用惯用、生产级的 Go 代码，覆盖从代码生成到工程规范的全流程。

---

## 核心原则

在编写任何代码前，牢记 Go 的设计哲学：

- **简洁优于巧妙**：优先选择可读性高的写法，而非"聪明"的写法
- **显式优于隐式**：错误必须显式处理，不得忽略；依赖必须明确传递
- **组合优于继承**：通过 interface 和 struct embedding 实现复用
- **标准库优先**：能用标准库解决的，不引入第三方依赖
- **并发是工具，不是默认选项**：只在真正需要并发时才使用 goroutine

---

## 代码生成流程

### 1. 理解需求
- 明确输入/输出类型、边界条件、性能要求
- 判断是否需要并发、是否涉及 I/O、是否需要测试覆盖
- 主动询问：Go 版本、是否已有 `go.mod`、目标平台

### 2. 设计结构
- 选择合适的包结构（见 `references/project-structure.md`）
- 先定义接口，再实现具体类型
- 确定错误处理策略（哨兵错误 vs 自定义类型 vs `fmt.Errorf` wrapping）

### 3. 编写代码

生成代码时，**始终遵循**以下规范：

#### 错误处理
```go
// ✅ 正确：显式处理，附加上下文
if err != nil {
    return fmt.Errorf("parse config: %w", err)
}

// ❌ 错误：忽略错误
result, _ := strconv.Atoi(s)

// ❌ 错误：panic 代替错误处理（除 main/init 初始化外）
if err != nil {
    panic(err)
}
```

#### 命名规范
```go
// 包名：小写单词，不用下划线
package httpserver   // ✅
package http_server  // ❌

// 接口名：通常以 -er 结尾，单方法接口
type Reader interface { Read([]byte) (int, error) }

// 导出名：PascalCase；非导出名：camelCase
// 缩写全大写：HTTPClient、URLParser、userID
```

#### 接口设计
```go
// ✅ 接口定义在使用方，而非实现方
// 接口尽量小（1~3 个方法）
type Storage interface {
    Get(ctx context.Context, key string) ([]byte, error)
    Set(ctx context.Context, key string, val []byte) error
}
```

#### Context 传递
```go
// ✅ Context 必须是函数第一个参数，命名为 ctx
func FetchUser(ctx context.Context, id int) (*User, error)

// ❌ 不能将 context 存储在 struct 中（除非是 Request 级别）
```

#### 并发安全
```go
// ✅ 使用 sync.Mutex 保护共享状态
type Cache struct {
    mu    sync.RWMutex
    items map[string]any
}

// ✅ 使用 errgroup 管理一组 goroutine
g, ctx := errgroup.WithContext(ctx)
g.Go(func() error { return fetchA(ctx) })
g.Go(func() error { return fetchB(ctx) })
if err := g.Wait(); err != nil { ... }

// ⚠️ 避免 goroutine 泄漏：每个 goroutine 必须有退出路径
```

### 4. 常见场景模板

需要特定场景的完整示例，参阅对应参考文件：
- HTTP 服务 → `references/http-service.md`
- 并发模式 → `references/concurrency.md`
- 错误处理进阶 → `references/error-handling.md`
- 测试规范 → `references/testing.md`

---

## 代码审查规则

审查 Go 代码时，检查以下问题（按优先级排序）：

### 🔴 严重问题（必须修复）
- [ ] goroutine 泄漏（channel 未关闭、无退出条件的 goroutine）
- [ ] 未加锁访问共享变量（data race）
- [ ] 忽略 `error` 返回值（`_, _ = f()` 形式）
- [ ] context 未传递导致无法取消
- [ ] 裸 `recover()` 吞掉 panic 信息

### 🟡 重要问题（应当修复）
- [ ] 错误信息未 wrap（`fmt.Errorf("...: %w", err)`）
- [ ] 过度使用 `interface{}` / `any`，应使用泛型或具体类型
- [ ] `init()` 函数中有副作用（如网络请求、文件读写）
- [ ] 循环中捕获变量引用而非值（Go 1.22 前的经典陷阱）
- [ ] 结构体未使用 `sync.Mutex` 的零值特性（不要用指针 mutex）

### 🔵 建议改进（可选）
- [ ] 函数超过 50 行，考虑拆分
- [ ] 注释缺失（导出符号必须有注释）
- [ ] 使用 `time.Sleep` 做节流，改用 `time.Ticker`
- [ ] `defer` 在循环中使用，可能延迟资源释放

---

## 项目结构快速参考

详细说明见 `references/project-layout.md`。

---

## 常用命令速查

```bash
# 模块管理
go mod init github.com/org/repo
go mod tidy                          # 整理依赖
go get github.com/pkg/errors@v0.9.1 # 添加依赖

# 构建
go build ./...
go build -ldflags="-X main.version=1.0.0" -o bin/app ./cmd/app

# 测试
go test ./...                        # 运行所有测试
go test -race ./...                  # 开启竞态检测
go test -bench=. -benchmem ./...     # 基准测试
go test -coverprofile=cover.out ./...
go tool cover -html=cover.out        # 查看覆盖率

# 代码质量
go vet ./...                         # 静态分析
gofmt -w .                           # 格式化
golangci-lint run                    # 综合 lint（需安装）

# 性能分析
go tool pprof http://localhost:6060/debug/pprof/profile
```

---

## Go 版本特性说明

生成代码时，**默认使用 Go 1.21+**，可使用：
- `slog`（结构化日志标准库）
- `slices` / `maps` 泛型工具包
- `errors.Join`
- `context.WithoutCancel`

若用户指定了较低版本，回退到对应兼容写法，并说明版本差异。

---

## 进化

本 skill 支持持续进化。当在编码过程中遇到以下情况时，自动将发现追加到 `docs/golang-coder-evolution.md`：

- 现有 reference 未覆盖的新模式或反模式
- 现有规则不适用特定场景，需要补充说明
- 重复出现的问题模式，值得记录

规则：

- 不主动读取此文件，仅在写入新条目时打开
- 如果 `docs/` 目录不存在，自动创建
- 当用户说「整理进化记录」时，读取文件，将待整理条目融入对应 reference

条目格式：

```markdown
- **[日期] 领域-简述**
  - 场景：描述遇到的具体场景
  - 建议：应该补充什么规则/示例到哪个 reference
  - 状态：待整理
```

---
## 全局原则

不管在哪个领域，始终遵循：

- 简洁优于复杂，组合优于继承
- 优先使用标准库，谨慎引入第三方依赖
- 错误必须处理，永远不要忽略
- 先写正确的代码，再考虑性能

---

## 参考文件索引

| 文件 | 内容 | 何时读取 |
|------|------|---------|
| `references/concurrency.md` | goroutine 模式、channel 设计、errgroup、worker pool | 涉及并发、异步处理时 |
| `references/error-handling.md` | 错误类型设计、sentinel error、errors.Is/As | 设计错误处理策略时 |
| `references/project-structure.md` | 目录结构规范、包命名、分层架构 | 新建项目或讨论项目结构时 |
| `references/http-service.md` | HTTP 服务模板、中间件、优雅关闭 | 构建 HTTP/REST 服务时 |
| `references/testing.md` | table-driven test、mock、testify、集成测试 | 编写或审查测试代码时 |
| `references/performance.md ` | 性能优化、减少分配 | 做性能优化时 |
| `references/naming.md` | 命名任何 Go 标识符 | 为变量命名时 |

## 快速导航

根据当前任务，查阅对应的参考文档：

| 任务 | 参考文档 |
|------|----------|
| 命名任何 Go 标识符 | references/naming.md |
| 组织项目或包结构 | references/project-structure.md |
| 处理错误、设计错误策略 | references/error-handling.md |
| 设计接口、选择抽象方式 | references/interface-design.md |

## 决策引导

遇到以下决策点时：

1. **需要抽象？** → references/interface-design.md（含决策流程图）
2. **其他** → 直接查对应 reference
