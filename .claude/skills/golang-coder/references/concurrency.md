# Go 并发模式参考

## 基本原则

> "不要通过共享内存来通信，而要通过通信来共享内存。"

- goroutine 很轻量（初始 ~2KB 栈），但不是无限免费的
- channel 用于传递数据所有权；mutex 用于保护共享状态
- 每个 goroutine 都必须有明确的退出路径

---

## 模式一：Context 取消

```go
func worker(ctx context.Context) error {
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            // 执行工作
        }
    }
}
```

---

## 模式二：errgroup 并发任务组

```go
import "golang.org/x/sync/errgroup"

func fetchAll(ctx context.Context, ids []int) ([]*User, error) {
    g, ctx := errgroup.WithContext(ctx)
    users := make([]*User, len(ids))

    for i, id := range ids {
        i, id := i, id // 捕获循环变量（Go 1.22 前必须）
        g.Go(func() error {
            u, err := fetchUser(ctx, id)
            if err != nil {
                return fmt.Errorf("fetch user %d: %w", id, err)
            }
            users[i] = u
            return nil
        })
    }

    if err := g.Wait(); err != nil {
        return nil, err
    }
    return users, nil
}
```

---

## 模式三：Worker Pool

```go
func workerPool(ctx context.Context, jobs <-chan Job, workerNum int) <-chan Result {
    results := make(chan Result, workerNum)

    var wg sync.WaitGroup
    for range workerNum {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobs {
                select {
                case results <- process(ctx, job):
                case <-ctx.Done():
                    return
                }
            }
        }()
    }

    go func() {
        wg.Wait()
        close(results)
    }()

    return results
}
```

---

## 模式四：Pipeline

```go
func generate(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            out <- n
        }
    }()
    return out
}

func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            out <- n * n
        }
    }()
    return out
}
```

---

## 常见陷阱

```go
// ❌ goroutine 泄漏：channel 发送方退出但接收方还在等待
ch := make(chan int)
go func() { ch <- 1 }() // 如果没人读，goroutine 永久阻塞

// ✅ 使用带缓冲 channel 或确保接收方存在

// ❌ 在循环中 defer（资源释放被推迟到函数返回）
for _, f := range files {
    r, _ := os.Open(f)
    defer r.Close() // 所有 defer 在函数结束时才执行
}

// ✅ 提取为独立函数或手动关闭
for _, f := range files {
    func() {
        r, _ := os.Open(f)
        defer r.Close()
        // 处理 r
    }()
}
```