# 错误处理

> 来源：Go 标准库、Google Go Style Guide、Uber Go Style Guide

## 核心原则

1. **错误是值，用常规代码处理，不用异常机制。** Go 的 `error` 是一个接口，错误处理就是普通的条件判断和值传递。
2. **添加上下文，不要只传递原始错误。** 用 `fmt.Errorf("...: %w", err)` 包装错误，让调用链上的每一层都能提供有意义的定位信息。
3. **用 `errors.Is`/`errors.As` 检查，不要用字符串匹配或 `==`。** `%w` 包装后的错误链无法被 `==` 匹配，`errors.Is` 和 `errors.As` 会遍历整条链。

## 速查规则

| 场景 | 做法 | 示例 |
|------|------|------|
| 包装错误 | `fmt.Errorf("...: %w", err)` | `fmt.Errorf("open config: %w", err)` |
| 检查特定错误 | `errors.Is(err, target)` | `errors.Is(err, os.ErrNotExist)` |
| 提取错误类型 | `errors.As(err, &target)` | `errors.As(err, &pathErr)` |
| 定义哨兵错误 | `var Err... = errors.New(...)` | `var ErrNotFound = errors.New("not found")` |
| 自定义错误类型 | `type ...Error struct{}` + `Error() string` | `type SyntaxError struct{ Line int }` |
| 处理不可恢复错误 | `panic(...)` | 仅用于编程错误 |

## 错误包装

✅ 用 `%w` 包装错误，保留原始错误链：

```go
func ReadConfig(path string) ([]byte, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("read config %s: %w", path, err)
    }
    return data, nil
}
```

✅ 调用方可以用 `errors.Is` 检查包装链中的任意错误：

```go
_, err := ReadConfig("config.toml")
if errors.Is(err, os.ErrNotExist) {
    fmt.Println("config file not found, using defaults")
}
```

❌ 用 `%v` 丢失原始错误，无法 unwrap：

```go
return nil, fmt.Errorf("read config %s: %v", path, err)
```

✅ 多层包装形成清晰的错误链：

```go
// 第一层
return fmt.Errorf("parse config: %w", err)

// 第二层
return fmt.Errorf("load app config: %w", err)

// 最终错误信息：
// load app config: parse config: open config.toml: no such file or directory
```

❌ 用 `+` 拼接错误字符串，丢失链信息：

```go
return nil, errors.New("read config: " + err.Error())
```

## 哨兵错误 vs 自定义错误类型

### 哨兵错误

适用于：简单的值比较，表示特定条件，不需要携带额外信息。

✅

```go
var (
    ErrNotFound   = errors.New("not found")
    ErrConflict   = errors.New("conflict")
    ErrPermission = errors.New("permission denied")
)

func GetUser(id string) (*User, error) {
    user, ok := db[id]
    if !ok {
        return nil, ErrNotFound
    }
    return user, nil
}

func HandleRequest() {
    _, err := GetUser("abc")
    if errors.Is(err, ErrNotFound) {
        http.Error(w, "user not found", http.StatusNotFound)
        return
    }
}
```

### 自定义错误类型

适用于：需要携带额外上下文信息，调用方需要根据错误中的字段做决策。

✅

```go
type SyntaxError struct {
    Line    int
    Column  int
    Message string
}

func (e *SyntaxError) Error() string {
    return fmt.Sprintf("syntax error at %d:%d: %s", e.Line, e.Column, e.Message)
}

func Parse(input string) (*AST, error) {
    return nil, &SyntaxError{Line: 10, Column: 5, Message: "unexpected token"}
}

func HandleParse() {
    _, err := Parse(input)
    var syntaxErr *SyntaxError
    if errors.As(err, &syntaxErr) {
        fmt.Printf("parse error at line %d\n", syntaxErr.Line)
    }
}
```

## 错误边界

### 何时处理

你可以恢复或降级时，在当前层处理。

✅

```go
func ReadConfig(path string) ([]byte, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        if errors.Is(err, os.ErrNotExist) {
            return DefaultConfig, nil
        }
        return nil, fmt.Errorf("read config: %w", err)
    }
    return data, nil
}
```

❌ 在无法做出合理决策的层处理：

```go
func ReadFile(path string) ([]byte, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        log.Fatal(err)
        return nil, err
    }
    return data, nil
}
```

### 何时传递

调用者能更好地决定如何处理时，包装后传递。

✅

```go
func LoadUser(id string) (*User, error) {
    row := db.QueryRow("SELECT * FROM users WHERE id = ?", id)
    var u User
    if err := row.Scan(&u.ID, &u.Name); err != nil {
        return nil, fmt.Errorf("scan user %s: %w", id, err)
    }
    return &u, nil
}
```

❌ 不包装直接传递，丢失上下文：

```go
func LoadUser(id string) (*User, error) {
    row := db.QueryRow("SELECT * FROM users WHERE id = ?", id)
    var u User
    if err := row.Scan(&u.ID, &u.Name); err != nil {
        return nil, err
    }
    return &u, nil
}
```

### 何时 wrap

在每一层调用都添加有意义的上下文后传递。

✅

```go
func (s *Service) HandleRequest(req Request) error {
    user, err := s.repo.LoadUser(req.UserID)
    if err != nil {
        return fmt.Errorf("handle request %s: %w", req.ID, err)
    }
    return nil
}
```

## panic 的使用边界

✅ 仅用于不可恢复的编程错误（如索引越界、nil 解引用、逻辑不变量被违反）：

```go
func MustCompile(pattern string) *regexp.Regexp {
    re, err := regexp.Compile(pattern)
    if err != nil {
        panic(`regexp: Compile(` + quote(pattern) + `): ` + err.Error())
    }
    return re
}
```

❌ 不要在库代码中用 panic 处理可恢复的错误：

```go
func Connect(addr string) *Conn {
    conn, err := net.Dial("tcp", addr)
    if err != nil {
        panic(err)
    }
    return conn
}
```

❌ 不要用 panic 做 control flow：

```go
func Process(items []string) []string {
    var result []string
    for _, item := range items {
        if item == "" {
            panic("empty item")
        }
        result = append(result, item)
    }
    return result
}
```

✅ 在顶层用 `recover()` 捕获 panic，防止整个程序崩溃：

```go
func safeHandler(fn func()) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic recovered: %v", r)
            debug.PrintStack()
        }
    }()
    fn()
    return nil
}

func main() {
    if err := safeHandler(func() {
        doWork()
    }); err != nil {
        log.Printf("fatal: %v", err)
        os.Exit(1)
    }
}
```

## 反模式

| 反模式 | 问题 | 修正 |
|--------|------|------|
| `_ = doSomething()` | 忽略错误 | 检查并处理每个错误 |
| `return err` 不包装 | 丢失调用上下文 | `return fmt.Errorf("do X: %w", err)` |
| `if err != nil { log.Fatal(err) }` | 库代码中直接退出 | 返回错误让调用者决定 |
| `strings.Contains(err.Error(), "...")` | 脆弱的字符串匹配 | 用 `errors.Is`/`errors.As` |
| `err == specificErr` | 不支持包装链 | 用 `errors.Is(err, specificErr)` |
| 深层嵌套 if err | 难以阅读 | early return，减少嵌套 |
| 库代码中 panic | 影响调用者稳定性 | 返回 error |

### 忽略错误

❌

```go
file, _ := os.Open("config.json")
```

✅

```go
file, err := os.Open("config.json")
if err != nil {
    return fmt.Errorf("open config: %w", err)
}
defer file.Close()
```

### 深层嵌套

❌

```go
func Process(data []byte) error {
    if len(data) > 0 {
        if isValid(data) {
            if result, err := transform(data); err == nil {
                if err := save(result); err != nil {
                    return err
                }
            } else {
                return err
            }
        } else {
            return ErrInvalid
        }
    } else {
        return ErrEmpty
    }
    return nil
}
```

✅ early return 减少嵌套：

```go
func Process(data []byte) error {
    if len(data) == 0 {
        return ErrEmpty
    }
    if !isValid(data) {
        return ErrInvalid
    }
    result, err := transform(data)
    if err != nil {
        return fmt.Errorf("transform data: %w", err)
    }
    if err := save(result); err != nil {
        return fmt.Errorf("save result: %w", err)
    }
    return nil
}
```
