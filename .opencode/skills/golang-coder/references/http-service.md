# Go HTTP 服务参考

## 最小可运行服务（标准库）

```go
package main

import (
    "context"
    "log/slog"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /health", healthHandler)
    mux.HandleFunc("GET /users/{id}", getUserHandler)

    srv := &http.Server{
        Addr:         ":8080",
        Handler:      mux,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    // 优雅关闭
    go func() {
        quit := make(chan os.Signal, 1)
        signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
        <-quit

        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()
        if err := srv.Shutdown(ctx); err != nil {
            slog.Error("shutdown failed", "err", err)
        }
    }()

    slog.Info("server starting", "addr", srv.Addr)
    if err := srv.ListenAndServe(); err != http.ErrServerClosed {
        slog.Error("server error", "err", err)
        os.Exit(1)
    }
}
```

---

## 中间件模式

```go
type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
    for i := len(middlewares) - 1; i >= 0; i-- {
        h = middlewares[i](h)
    }
    return h
}

// 日志中间件
func Logging(logger *slog.Logger) Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            next.ServeHTTP(w, r)
            logger.Info("request",
                "method", r.Method,
                "path", r.URL.Path,
                "duration", time.Since(start),
            )
        })
    }
}
```

---

## JSON 响应助手

```go
func writeJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    if err := json.NewEncoder(w).Encode(v); err != nil {
        slog.Error("encode response", "err", err)
    }
}

func readJSON(r *http.Request, v any) error {
    dec := json.NewDecoder(r.Body)
    dec.DisallowUnknownFields()
    return dec.Decode(v)
}
```

---

## 依赖注入推荐模式

```go
// 将依赖挂在 handler struct 上，而非全局变量
type UserHandler struct {
    svc    UserService
    logger *slog.Logger
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id") // Go 1.22+
    user, err := h.svc.GetUser(r.Context(), id)
    if err != nil {
        writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
        return
    }
    writeJSON(w, http.StatusOK, user)
}
```