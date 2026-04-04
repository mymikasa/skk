# 性能优化

> 来源：Go 标准库、Google Go Style Guide、Uber Go Style Guide

## 核心原则

1. **先 profile，再优化**——不要猜测性能瓶颈
2. **减少堆分配是最有效的优化手段**
3. **值类型 vs 指针**：理解逃逸分析，不必要的指针会导致堆分配

## 速查规则

| 场景 | 做法 |
|------|------|
| 切片大小已知 | `make([]T, 0, n)` 预分配容量 |
| map 大小已知 | `make(map[K]V, n)` 预分配 |
| 字符串拼接 | `strings.Builder` |
| 对象复用 | `sync.Pool` |
| `[]byte` ↔ `string` 转换 | 热路径用 `unsafe`（只读场景）或统一用 `[]byte` |
| 热路径日志 | 先检查日志级别再格式化 |
| JSON 处理 | 流式 `json.Decoder` 优于 `json.Unmarshal` 整个 body |
| 大数组传参 | 用切片或指针，不要值拷贝 |

## 减少分配

### 切片预分配

✅ 预分配容量：

```go
items := make([]Item, 0, len(ids))
for _, id := range ids {
    items = append(items, Item{ID: id})
}
```

❌ 不指定容量，多次扩容：

```go
var items []Item
for _, id := range ids {
    items = append(items, Item{ID: id})
}
```

### Map 预分配

✅ 已知数量时预分配 ❌ 空创建：

```go
m := make(map[string]int, len(keys))
for _, k := range keys {
    m[k] = 0
}
```

```go
m := make(map[string]int)
for _, k := range keys {
    m[k] = 0
}
```

### 字符串拼接

✅ `strings.Builder`（记得 `Grow`） ❌ `+` 或 `fmt.Sprintf`：

```go
var b strings.Builder
b.Grow(len(parts) * 16)
for _, p := range parts {
    b.WriteString(p)
}
result := b.String()
```

```go
result := ""
for _, p := range parts {
    result += p
}
result = fmt.Sprintf("%s%s%s", a, b, c)
```

### sync.Pool 复用临时对象

✅ 复用减少 GC 压力 ❌ 每次新建：

```go
var bufPool = sync.Pool{
    New: func() any {
        return new(bytes.Buffer)
    },
}

func process(data []byte) {
    buf := bufPool.Get().(*bytes.Buffer)
    defer func() {
        buf.Reset()
        bufPool.Put(buf)
    }()
    buf.Write(data)
}
```

```go
func process(data []byte) {
    buf := new(bytes.Buffer)
    buf.Write(data)
}
```

## 逃逸分析

编译器决定变量分配在**栈**（快、自动释放）还是**堆**（需要 GC）。逃逸到堆意味着额外分配开销、GC 压力、缓存不友好。常见逃逸场景：

- 返回局部变量的指针 → 逃逸到堆
- 赋值给 `interface{}` 变量 → 逃逸
- 闭包捕获局部变量 → 可能逃逸
- 切片扩容超出容量 → 新底层数组堆分配

查看逃逸决策：

```sh
go build -gcflags="-m" ./...
```

✅ 小结构体用值类型，留在栈上：

```go
type Point struct{ X, Y float64 }

func Add(a, b Point) Point {
    return Point{X: a.X + b.X, Y: a.Y + b.Y}
}
```

❌ 不必要的指针导致堆分配：

```go
func NewPoint(x, y float64) *Point {
    return &Point{X: x, Y: y}
}
```

✅ 大结构体或需要共享修改时用指针合理：

```go
type Buffer struct {
    data []byte
    off  int
}

func (b *Buffer) Write(p []byte) {
    b.data = append(b.data, p...)
}
```

## Profiling 流程

### HTTP 服务启用 pprof

```go
import _ "net/http/pprof"

func main() {
    go http.ListenAndServe("localhost:6060", nil)
}
```

```sh
go tool pprof http://localhost:6060/debug/pprof/profile
go tool pprof http://localhost:6060/debug/pprof/heap
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

### CLI 工具用 runtime/pprof

```go
func main() {
    f, _ := os.Create("cpu.prof")
    defer f.Close()
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()

    fMem, _ := os.Create("mem.prof")
    defer fMem.Close()
    defer pprof.WriteHeapProfile(fMem)
}
```

### Benchmark 带内存分析

```sh
go test -bench=. -benchmem ./...
go test -bench=BenchmarkFunc -benchmem -count=5 ./...
```

### 分析 profile

```sh
go tool pprof -http=:8080 cpu.prof     # Web UI
go tool pprof -top cpu.prof            # 文本 top
go tool pprof -list=FuncName cpu.prof  # 逐行分析
```

## 常见热点

### `[]byte` ↔ `string` 转换

每次 `[]byte(s)` 或 `string(b)` 都分配新内存并拷贝。热路径用 `unsafe` 只读转换：

```go
func unsafeString(b []byte) string {
    return *(*string)(unsafe.Pointer(&b))
}
```

> `unsafe` 转换结果必须只读，修改会导致未定义行为。

### `interface{}` / `any` 装箱

赋值给 `interface{}` 产生堆分配，泛型可避免：

```go
func Max[T constraints.Ordered](a, b T) T {
    if a > b {
        return a
    }
    return b
}
```

### 反射开销

反射比直接调用慢 10-100 倍，热路径避免 `reflect`，优先代码生成或泛型。多 goroutine 争抢同一 `sync.Mutex` 成为瓶颈时：分片锁、`sync.RWMutex`（读多写少）、atomic 无锁结构。`defer` 有少量开销（Go 1.20+ 已大幅优化），一般场景优先用 `defer` 保证正确性。

## 反模式

| 反模式 | 问题 | 修正 |
|--------|------|------|
| 过早优化 | 浪费时间，代码变复杂 | 先 profile 确认瓶颈 |
| 热路径日志 `log.Infof("result: %v", data)` | Sprint 在日志级别不够时也执行 | `if log.Level() >= Info { log.Info(...) }` |
| 频繁 `[]byte(string)` 转换 | 每次转换都分配+拷贝 | 统一类型，或用 unsafe 只读转换 |
| 热路径中用 `interface{}` | 装箱导致堆分配+GC 压力 | 用泛型或具体类型 |
| 不必要的指针 | 逃逸到堆 | 小值类型直接用值 |
| 热路径中用 `defer` | defer 有少量开销 | 热路径函数可以不用 defer（谨慎） |
| 无缓冲 channel | 阻塞到对端就绪 | 明确是否需要同步语义 |
| 字符串拼接用 `+` | 每次拼接分配新字符串 | 用 `strings.Builder` |
