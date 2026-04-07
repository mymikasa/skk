# Go 测试规范参考

## Table-Driven Tests（标准模式）

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name    string
        a, b    int
        want    int
        wantErr bool
    }{
        {"正数相加", 1, 2, 3, false},
        {"负数", -1, -2, -3, false},
        {"溢出", math.MaxInt, 1, 0, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Add(tt.a, tt.b)
            if (err != nil) != tt.wantErr {
                t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
            }
            if got != tt.want {
                t.Errorf("Add() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

---

## Mock 接口

推荐使用 `github.com/stretchr/testify/mock` 或手写 mock：

```go
// 手写 mock（简单场景首选）
type mockUserRepo struct {
    users map[int]*User
    err   error
}

func (m *mockUserRepo) GetByID(_ context.Context, id int) (*User, error) {
    if m.err != nil {
        return nil, m.err
    }
    return m.users[id], nil
}

// 测试中使用
func TestUserService_GetUser(t *testing.T) {
    repo := &mockUserRepo{
        users: map[int]*User{1: {ID: 1, Name: "Alice"}},
    }
    svc := NewUserService(repo)
    user, err := svc.GetUser(context.Background(), 1)
    // ...
}
```

---

## HTTP Handler 测试

```go
func TestGetUserHandler(t *testing.T) {
    handler := &UserHandler{svc: &mockUserService{...}}

    req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
    w := httptest.NewRecorder()

    handler.GetUser(w, req)

    res := w.Result()
    if res.StatusCode != http.StatusOK {
        t.Errorf("status = %d, want %d", res.StatusCode, http.StatusOK)
    }
}
```

---

## 测试辅助规范

```go
// 使用 t.Helper() 标记辅助函数，让错误指向调用行
func assertEqual(t *testing.T, got, want any) {
    t.Helper()
    if got != want {
        t.Errorf("got %v, want %v", got, want)
    }
}

// 使用 t.Cleanup 替代 defer（更安全）
func TestWithDB(t *testing.T) {
    db := setupTestDB(t)
    t.Cleanup(func() { db.Close() })
}

// 跳过慢速测试
func TestIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test in short mode")
    }
}
```

---

## 基准测试

```go
func BenchmarkProcess(b *testing.B) {
    data := generateTestData()
    b.ResetTimer()           // 不计入数据准备时间
    b.ReportAllocs()         // 报告内存分配

    for b.Loop() {           // Go 1.24+，旧版用 range b.N
        Process(data)
    }
}
```