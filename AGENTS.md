# AGENTS.md — SKK

Go project. MIT licensed, copyright 2026 Mikasa.
Always response in Chinese.

## Repository Layout

```
.
├── .gitignore
├── LICENSE
├── go.mod          # module definition (to be created)
├── go.sum
├── cmd/            # main packages / entry points
│   └── skk/
│       └── main.go
├── internal/       # private application code
├── pkg/            # public library code (if any)
└── AGENTS.md
```

## Build / Lint / Test Commands

### Build
```sh
go build ./...
```

Build a specific binary:
```sh
go build -o bin/skk ./cmd/skk
```

### Test
```sh
go test ./...
```

Run a single test by name:
```sh
go test -run TestFunctionName ./path/to/package/
```

Run a single test file:
```sh
go test -run TestFunctionName ./path/to/package/ -v
```

Run with verbose output:
```sh
go test -v ./...
```

Run with race detector:
```sh
go test -race ./...
```

Run with coverage:
```sh
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Lint
```sh
go vet ./...
```

If golangci-lint is configured:
```sh
golangci-lint run
golangci-lint run ./path/to/package/
```

### Format
```sh
gofmt -w .
goimports -w .
```

### Tidy dependencies
```sh
go mod tidy
```

## Code Style Guidelines

### Formatting

- Use `gofmt` (or `goimports`) — no custom formatting configuration.
- Tabs for indentation (Go standard).
- No trailing whitespace.
- Files should end with a newline.

### Imports

- Use `goimports` to manage imports automatically.
- Group imports into three sections separated by blank lines:
  1. Standard library packages
  2. Third-party packages
  3. Local / project packages
- Do not use dot imports (`.`) except in tests.
- Avoid blank imports (`_`) unless required (e.g., driver registration).

Example:
```go
import (
    "context"
    "fmt"

    "google.golang.org/grpc"

    "github.com/mikasa/skk/internal/foo"
)
```

### Naming Conventions

- **Packages**: lowercase, single word, no underscores (e.g., `http`, `json`, `foo`).
- **Files**: `snake_case.go` (e.g., `user_handler.go`).
- **Types/Interfaces**: `PascalCase` (e.g., `UserService`, `Reader`).
- **Functions/Methods**: `PascalCase` (exported), `camelCase` (unexported).
- **Constants**: `PascalCase` (exported), `camelCase` (unexported). Group in `const` blocks.
- **Variables**: same as functions — `PascalCase` / `camelCase`.
- **Acronyms**: keep consistent casing — `HTTPClient`, `userID`, `apiURL`.
- **Interfaces**: prefer single-method interfaces named by the method + "-er" suffix (e.g., `Reader`, `Stringer`). Do not prefix with `I` or `Interface`.
- **Test files**: `foo_test.go` in the same package. Prefer same-package tests (white-box) unless black-box testing is specifically needed.

### Types

- Prefer defining types in the package where they are used.
- Use type aliases only when there is a clear reason.
- Prefer small, focused interfaces over large ones.
- Use `type` for domain-specific types to improve readability:
  ```go
  type UserID string
  ```
- Prefer structs over multiple positional function parameters.

### Error Handling

- **Never ignore errors** — always check and handle them.
- Do not use `_` for error return values.
- Wrap errors with context using `fmt.Errorf("doing X: %w", err)` for chain inspection.
- Use `errors.Is()` and `errors.As()` for error checking, not `==` or type assertions.
- Define sentinel errors as package-level `var`:
  ```go
  var ErrNotFound = errors.New("not found")
  ```
- Return errors early; avoid deep nesting.
- In tests, use `t.Fatal` / `t.Fatalf` for setup errors and `t.Error` / `t.Errorf` for test assertion failures.

### Comments

- Package comments: start with `// Package foo ...` in one file per package.
- Exported identifiers must have doc comments starting with the identifier name.
- Comments should be complete sentences ending with a period.
- Do not add comments that merely restate the code.

### Concurrency

- Prefer channels for communication; share memory by communicating.
- Use `sync.Mutex` / `sync.RWMutex` when shared state is necessary.
- Always protect shared mutable state.
- Use `context.Context` as the first parameter in functions that perform I/O or may block.
- Prefer `context.Context` over custom cancellation mechanisms.

### Testing

- Table-driven tests are the preferred pattern:
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
- Use `t.Parallel()` where safe.
- Use `t.Helper()` in test utility functions.
- Prefer standard library testing; introduce testify or similar only if already in use.

### General

- Follow [Effective Go](https://go.dev/doc/effective_go) and [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).
- Keep functions short and focused.
- Prefer composition over inheritance (embed structs instead).
- Use `defer` for cleanup (closing files, unlocking mutexes, etc.).
- Avoid `panic` in library code — return errors instead.
- Avoid `init()` functions unless absolutely necessary.
- Run `go vet` and `golangci-lint` before committing.
npx skills add https://github.com/jeffallan/claude-skills --skill golang-pro


npx skills add https://github.com/cxuu/golang-skills --skill go-naming