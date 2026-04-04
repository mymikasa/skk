# 测试实践

> 来源：Go 标准库、Google Go Style Guide、Uber Go Style Guide

## 核心原则

1. **表驱动测试是首选模式**——一个测试函数覆盖多组输入输出，结构清晰、易扩展
2. **测试应该快速、确定、无副作用**——相同的输入永远得到相同的结果，可在任意环境重复运行
3. **测试公开行为，不测试实现细节**——重构不应破坏测试，公开 API 是唯一的契约

## 速查规则

| 场景 | 做法 |
|------|------|
| 基本测试 | 表驱动 `[]struct{...}` + `t.Run` |
| 并行测试 | `t.Parallel()`（注意闭包捕获 `tt`） |
| 测试辅助函数 | `t.Helper()` 标记 |
| 基准测试 | `BenchmarkXxx(b *testing.B)` |
| 设置/清理 | `t.Cleanup(func(){})` |
| 跳过测试 | `t.Skip()` / `t.Skipf()` |
| mock 依赖 | 接口 + 手写 fake（优先于 mock 框架） |
| 子测试隔离 | 每个子测试独立 setup |
| 测试覆盖率 | `go test -cover` / `-coverprofile` |

## 表驱动测试

✅ **推荐写法**——用匿名结构体切片组织用例，`t.Run` 创建子测试：

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive", 1, 2, 3},
        {"negative", -1, -2, -3},
        {"zero", 0, 0, 0},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Add(tt.a, tt.b)
            if got != tt.expected {
                t.Errorf("Add(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.expected)
            }
        })
    }
}
```

要点：

- 每个用例有 `name` 字段，`t.Run` 输出可读的子测试名
- 测试结构体字段用有意义命名，避免 `input1`、`input2`
- 用 `t.Errorf` 报告期望与实际的差异，包含足够上下文

## t.Parallel 注意事项

✅ **正确**——在闭包内重新绑定循环变量：

```go
func TestParse(t *testing.T) {
    tests := []struct {
        name  string
        input string
        want  int
    }{
        {"decimal", "10", 10},
        {"hex", "0xff", 255},
    }
    for _, tt := range tests {
        tt := tt
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            got, err := Parse(tt.input)
            if err != nil {
                t.Fatalf("Parse(%q) error: %v", tt.input, err)
            }
            if got != tt.want {
                t.Errorf("Parse(%q) = %d, want %d", tt.input, got, tt.want)
            }
        })
    }
}
```

❌ **错误**——闭包直接捕获循环变量，并行执行时所有子测试读到同一个值：

```go
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        t.Parallel()
        got := Parse(tt.input)
        if got != tt.want {
            t.Errorf("got %d, want %d", got, tt.want)
        }
    })
}
```

`tt := tt` 必须在 `t.Run` 的闭包内部、`t.Parallel()` 之前声明，否则并行子测试共享外层循环变量的最终值。

## t.Helper

测试辅助函数用 `t.Helper()` 标记，让 `t.Errorf` / `t.Fatalf` 报告调用者的文件和行号，而非辅助函数内部的位置。

✅ **正确**：

```go
func assertEqual(t *testing.T, got, want int) {
    t.Helper()
    if got != want {
        t.Errorf("got %d, want %d", got, want)
    }
}

func TestCalc(t *testing.T) {
    assertEqual(t, Calc(2, 3), 5)
}
```

失败时错误指向 `TestCalc` 的调用行，而非 `assertEqual` 内部。

## 基准测试

✅ **BenchmarkXxx**——`b.N` 由测试框架自动调整：

```go
func BenchmarkFormat(b *testing.B) {
    data := []byte(`{"key":"value"}`)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = Format(data)
    }
}

func BenchmarkFormatParallel(b *testing.B) {
    data := []byte(`{"key":"value"}`)
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            _, _ = Format(data)
        }
    })
}
```

要点：

- 耗时的 setup 放在 `b.ResetTimer()` 之前
- 并发基准测试用 `b.RunParallel`
- 运行：`go test -bench=. -benchmem`

## Mock 策略

### 优先级

1. **接口 + 手写 fake**（推荐）——简单、可维护、无外部依赖
2. **仅复杂场景使用 mock 框架**（如 testify/mock）——交互复杂、状态多时

### ✅ Fake 实现

定义接口隔离依赖，手写一个轻量 fake：

```go
type UserRepository interface {
    GetByID(ctx context.Context, id string) (*User, error)
    Save(ctx context.Context, u *User) error
}

type FakeUserRepository struct {
    users map[string]*User
    err   error
}

func NewFakeUserRepository() *FakeUserRepository {
    return &FakeUserRepository{users: make(map[string]*User)}
}

func (f *FakeUserRepository) GetByID(_ context.Context, id string) (*User, error) {
    if f.err != nil {
        return nil, f.err
    }
    u, ok := f.users[id]
    if !ok {
        return nil, ErrNotFound
    }
    return u, nil
}

func (f *FakeUserRepository) Save(_ context.Context, u *User) error {
    if f.err != nil {
        return f.err
    }
    f.users[u.ID] = u
    return nil
}
```

测试中使用：

```go
func TestGetUser(t *testing.T) {
    repo := NewFakeUserRepository()
    repo.users["abc"] = &User{ID: "abc", Name: "Alice"}

    svc := NewUserService(repo)
    u, err := svc.GetUser(context.Background(), "abc")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if u.Name != "Alice" {
        t.Errorf("got %q, want %q", u.Name, "Alice")
    }
}
```

### Fake vs Mock 何时使用

| 场景 | 选择 |
|------|------|
| 依赖返回数据，测试处理逻辑 | Fake（实现接口的简化版本） |
| 需要验证调用顺序、调用次数 | Mock 框架 |
| 同一接口多个测试复用 | Fake（一次实现，到处复用） |
| 一次性、简单的存根 | 测试内匿名结构体实现接口 |

## 测试组织

### 文件布局

```
internal/user/
├── service.go
├── service_test.go
├── handler.go
├── handler_test.go
└── testdata/
    └── golden_response.json
```

### 单元测试

- 与被测代码**同目录**，文件名 `*_test.go`
- 默认同包测试（白盒），可访问未导出符号
- 需要黑盒时用 `package user_test`（导入 `package` 路径）

### 集成测试

- 放在项目根目录的 `test/` 目录下，或使用 build tags：
  ```go
  //go:build integration
  ```
- 运行：`go test -tags=integration ./...`

### Test Fixtures

- 测试数据放在 `testdata/` 目录，`go test` 自动忽略该目录
- 读取：`os.ReadFile("testdata/input.json")`
- Golden file 模式：保存期望输出，测试时对比

### Setup / Teardown

```go
func TestMain(m *testing.M) {
    db := setupDB()
    code := m.Run()
    db.Close()
    os.Exit(code)
}
```

单个测试用 `t.Cleanup`：

```go
func TestWithTempDir(t *testing.T) {
    dir := t.TempDir()
    t.Cleanup(func() {
        os.RemoveAll(dir)
    })
}
```

## 反模式

| 反模式 | 问题 | 修正 |
|--------|------|------|
| 测试间共享状态 | 测试顺序影响结果 | 每个测试独立 setup |
| 依赖时间/随机数 | 测试不确定 | 注入时间源，固定随机种子 |
| 依赖外部服务 | CI 环境不稳定 | 用 fake/mock 替代 |
| 测试私有函数 | 重构时大量测试失效 | 通过公开 API 测试 |
| 过度使用 mock 框架 | mock 与实现紧耦合 | 手写 fake，验证行为而非调用次数 |
| 不检查错误返回 | `_ = SomeFunc()` | 用 `t.Fatal` 处理 setup 错误 |
| 忽略 `t.Error` 后继续 | 测试已经失败还继续 | 用 `t.Fatal` 在 setup 阶段提前退出 |
