# Go 错误处理参考

## 错误处理三层模型

```
调用层        →  检查 errors.Is / errors.As
中间层        →  fmt.Errorf("操作: %w", err) 添加上下文
底层          →  返回原始错误或哨兵错误
```

---

## 哨兵错误（Sentinel Errors）

用于调用方需要精确匹配的场景：

```go
// 定义（在包级别）
var (
    ErrNotFound   = errors.New("not found")
    ErrPermission = errors.New("permission denied")
)

// 使用
if errors.Is(err, ErrNotFound) {
    // 处理未找到的情况
}
```

---

## 自定义错误类型

需要携带额外信息时：

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error on field %q: %s", e.Field, e.Message)
}

// 检查类型
var ve *ValidationError
if errors.As(err, &ve) {
    fmt.Println("字段:", ve.Field)
}
```

---

## 错误 Wrapping

```go
// ✅ 使用 %w 保留错误链，允许 errors.Is/As 穿透
return fmt.Errorf("query user %d: %w", id, err)

// ✅ Go 1.20+ 合并多个错误
return errors.Join(err1, err2)

// ❌ 使用 %v 会丢失错误链
return fmt.Errorf("query user %d: %v", id, err) // 不能被 errors.Is 匹配
```

---

## HTTP 错误处理模式

```go
type AppError struct {
    Code    int    // HTTP 状态码
    Message string // 用户可见消息
    Err     error  // 内部错误（不对外暴露）
}

func (e *AppError) Error() string { return e.Message }
func (e *AppError) Unwrap() error { return e.Err }

func handler(w http.ResponseWriter, r *http.Request) {
    if err := doSomething(r.Context()); err != nil {
        var appErr *AppError
        if errors.As(err, &appErr) {
            http.Error(w, appErr.Message, appErr.Code)
            return
        }
        http.Error(w, "internal server error", http.StatusInternalServerError)
    }
}
```