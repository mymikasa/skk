# golang-coder Skill 实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 创建一个覆盖 Go 编码全生命周期的 opencode skill（golang-coder），包含 7 个领域的参考文档和自动进化机制。

**Architecture:** 混合型 skill 架构 — 轻量主入口（SKILL.md）+ 7 个模块化 reference 文件 + 项目级进化日志。主入口提供触发条件和快速导航，reference 文件采用速查手册格式（规则+示例+反模式），接口设计和并发领域加入决策流程图。

**Tech Stack:** Markdown（opencode skill 格式，front matter + 内容）

---

## 文件结构

```
skills/golang-coder/
├── SKILL.md                          # 主入口：触发条件 + 决策引导 + 领域索引
└── references/
    ├── naming.md                     # 从 go-naming 吸收
    ├── project-structure.md          # 项目结构
    ├── error-handling.md             # 错误处理
    ├── interface-design.md           # 接口设计（含决策流程图）
    ├── concurrency.md                # 并发编程（含决策流程图）
    ├── testing.md                    # 测试实践
    └── performance.md                # 性能优化
docs/
└── golang-coder-evolution.md         # 进化日志（skill 不加载）
```

完成后移除 `skills/go-naming/`。

---

### Task 1: 创建目录结构

**Files:**
- Create: `skills/golang-coder/SKILL.md`（占位）
- Create: `skills/golang-coder/references/` 目录
- Create: `docs/golang-coder-evolution.md`

- [ ] **Step 1: 创建目录和初始文件**

```bash
mkdir -p skills/golang-coder/references
mkdir -p docs
```

- [ ] **Step 2: 创建进化日志初始文件**

创建 `docs/golang-coder-evolution.md`：

```markdown
# golang-coder 进化日志

## 待整理

（暂无条目）

## 已整理

（暂无条目）
```

- [ ] **Step 3: Commit**

```bash
git add docs/golang-coder-evolution.md
git commit -m "feat: init golang-coder evolution log"
```

---

### Task 2: 编写 SKILL.md 主入口

**Files:**
- Create: `skills/golang-coder/SKILL.md`

- [ ] **Step 1: 编写 SKILL.md**

```markdown
---
name: golang-coder
description: 在编写、修改或设计 Go 代码时使用。覆盖命名规范、项目结构、错误处理、接口设计、并发编程、测试实践、性能优化。面向有经验的 Go 开发者，提供速查规则和决策引导。
---

# Go Coder

## 触发条件

当以下任一条件满足时激活：
- 正在创建或编辑 `.go` 文件
- 讨论 Go 包结构、接口设计、并发模式
- 编写 Go 测试或 benchmark
- 排查 Go 性能问题或竞态条件
- 设计 Go API 或错误处理策略

## 快速导航

根据当前任务，查阅对应的参考文档：

| 任务 | 参考文档 |
|------|----------|
| 命名任何 Go 标识符 | references/naming.md |
| 组织项目或包结构 | references/project-structure.md |
| 处理错误、设计错误策略 | references/error-handling.md |
| 设计接口、选择抽象方式 | references/interface-design.md |
| 使用 goroutine/channel/sync | references/concurrency.md |
| 编写测试或 benchmark | references/testing.md |
| 性能优化、减少分配 | references/performance.md |

## 决策引导

遇到以下决策点时：

1. **需要抽象？** → references/interface-design.md（含决策流程图）
2. **并发方案？** → references/concurrency.md（含决策流程图）
3. **其他** → 直接查对应 reference

## 全局原则

不管在哪个领域，始终遵循：
- 简洁优于复杂，组合优于继承
- 优先使用标准库，谨慎引入第三方依赖
- 错误必须处理，永远不要忽略
- 先写正确的代码，再考虑性能

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
```

- [ ] **Step 2: 验证 front matter 格式正确**

检查 `---` 分隔符、name 和 description 字段存在且无语法错误。

- [ ] **Step 3: Commit**

```bash
git add skills/golang-coder/SKILL.md
git commit -m "feat: add golang-coder SKILL.md entry point"
```

---

### Task 3: 迁移 naming.md

**Files:**
- Create: `skills/golang-coder/references/naming.md`（从 `skills/go-naming/SKILL.md` 迁移）

- [ ] **Step 1: 复制并调整命名文件**

将 `skills/go-naming/SKILL.md` 的内容复制到 `skills/golang-coder/references/naming.md`，移除 front matter 中的 `name` 和 `description`，改为简单注释说明来源。

保持文件中所有内容不变：速查表、包/类型/接口/函数/变量/常量/接收器/错误命名、缩写词表、决策流程图。

文件头部：

```markdown
# 命名规范

> 来源：Go 标准库、Google Go Style Guide、Uber Go Style Guide

（保持原有内容不变，从 "# Go 命名" 开始的所有内容）
```

- [ ] **Step 2: 移除旧 skill 目录**

```bash
rm -rf skills/go-naming/
```

- [ ] **Step 3: Commit**

```bash
git add skills/golang-coder/references/naming.md
git add skills/go-naming/
git commit -m "feat: migrate go-naming into golang-coder references"
```

---

### Task 4: 编写 project-structure.md

**Files:**
- Create: `skills/golang-coder/references/project-structure.md`

- [ ] **Step 1: 编写内容**

文件遵循统一格式：核心原则 → 速查规则（表格）→ 代码示例（✅/❌）→ 反模式。

内容要点：
1. **核心原则**（3 条）：
   - 包按职责划分，每个包做一件事
   - `internal/` 限制公开范围，`pkg/` 暴露可复用代码
   - 依赖方向：从 `cmd/` → `internal/` → 标准库/第三方

2. **标准目录布局**（表格）：

| 目录 | 用途 | 必需 |
|------|------|------|
| `cmd/` | 入口 main 包 | 是 |
| `internal/` | 私有应用代码 | 是 |
| `pkg/` | 可被外部复用的公共代码 | 否 |
| `api/` | API 定义（proto, swagger） | 否 |
| `configs/` | 配置文件模板 | 否 |
| `docs/` | 设计文档 | 否 |

3. **包拆分规则**（✅/❌ 示例）：
   - 按领域拆分 vs 按技术层拆分
   - 避免泛化包名 `utils`、`common`、`helpers`
   - 包名应描述提供什么，而非包含什么

4. **依赖管理**：
   - `go.mod` 最佳实践
   - 何时用 `replace` 指令
   - 版本选择策略

5. **反模式**：
   - 上帝包（一个包做所有事）
   - 循环依赖
   - 过深目录嵌套（超过 3 层）
   - `internal/` 下再按 MVC 分层

每个文件预计 150-250 行。

- [ ] **Step 2: 审查内容**

检查：规则是否明确、示例是否正确、反模式是否有说服力。

- [ ] **Step 3: Commit**

```bash
git add skills/golang-coder/references/project-structure.md
git commit -m "feat: add project-structure reference"
```

---

### Task 5: 编写 error-handling.md

**Files:**
- Create: `skills/golang-coder/references/error-handling.md`

- [ ] **Step 1: 编写内容**

内容要点：
1. **核心原则**（3 条）：
   - 错误是值，用常规代码处理
   - 添加上下文，不要只传递
   - 用 `errors.Is`/`errors.As` 检查，不要用 `==` 或类型断言

2. **速查规则**（表格）：

| 场景 | 做法 | 示例 |
|------|------|------|
| 包装错误 | `fmt.Errorf("...: %w", err)` | `fmt.Errorf("open config: %w", err)` |
| 检查错误 | `errors.Is` / `errors.As` | `errors.Is(err, os.ErrNotExist)` |
| 哨兵错误 | `var Err... = errors.New(...)` | `var ErrNotFound = errors.New("not found")` |
| 自定义错误类型 | `type ...Error struct{}` + `Error() string` | `type SyntaxError struct{ Line int }` |
| panic | 仅用于不可恢复的编程错误 | `panic("unreachable")` |

3. **代码示例**（✅/❌）：
   - 错误包装 + unwrap 链
   - 哨兵错误定义和使用
   - 自定义错误类型
   - 错误边界：处理 vs 传递

4. **反模式**：
   - `_` 忽略错误
   - 丢失上下文（只返回 err 不 wrap）
   - 过度嵌套的 if err != nil
   - 在库代码中使用 panic

- [ ] **Step 2: 审查内容**
- [ ] **Step 3: Commit**

```bash
git add skills/golang-coder/references/error-handling.md
git commit -m "feat: add error-handling reference"
```

---

### Task 6: 编写 interface-design.md（含决策流程图）

**Files:**
- Create: `skills/golang-coder/references/interface-design.md`

- [ ] **Step 1: 编写内容**

内容要点：
1. **核心原则**（3 条）：
   - 越小越好，单方法接口是理想状态
   - 接口定义在使用方，不在实现方
   - 隐式实现——不需要声明 `implements`

2. **速查规则**（表格）：

| 场景 | 做法 |
|------|------|
| 单方法抽象 | 接口，命名用 `-er` 后缀 |
| 多方法契约 | 小接口（2-3 个方法） |
| 多种实现且需要运行时切换 | 接口 |
| 编译期类型安全，单一实现 | 泛型 |
| 一次性行为参数 | 函数类型 `func(...) (...)` |

3. **决策流程图**（graphviz dot）：

```
需要抽象？
├─ 是否只有单个行为需要抽象？
│  ├─ 是 → 函数类型（func type）
│  └─ 否 → 需要多种实现？
│     ├─ 是 → 接口
│     └─ 否 → 是否需要编译期类型安全？
│        ├─ 是 → 泛型
│        └─ 否 → 具体类型即可，不需要抽象
```

4. **代码示例**（✅/❌）：
   - 小接口设计：`io.Reader`、`io.Writer` 模式
   - 接口组合：embed 多个小接口
   - 使用方定义接口

5. **反模式**：
   - 大接口（超过 5 个方法）
   - 在实现方定义接口
   - 返回接口（应返回具体类型）
   - 接口工厂模式

- [ ] **Step 2: 审查内容**
- [ ] **Step 3: Commit**

```bash
git add skills/golang-coder/references/interface-design.md
git commit -m "feat: add interface-design reference with decision flowchart"
```

---

### Task 7: 编写 concurrency.md（含决策流程图）

**Files:**
- Create: `skills/golang-coder/references/concurrency.md`

- [ ] **Step 1: 编写内容**

内容要点：
1. **核心原则**（3 条）：
   - 不要通过共享内存来通信，通过通信来共享内存
   - goroutine 必须有明确的生命周期管理
   - 共享可变状态必须保护

2. **速查规则**（表格）：

| 场景 | 做法 |
|------|------|
| 保护共享状态 | `sync.Mutex` / `sync.RWMutex` |
| 等待一组 goroutine | `sync.WaitGroup` |
| 一次性初始化 | `sync.Once` |
| 对象复用 | `sync.Pool` |
| goroutine 间通信 | channel |
| 取消/超时 | `context.WithCancel`/`WithTimeout` |
| 多生产者多消费者 | fan-out/fan-in channel 模式 |
| 流水线处理 | pipeline channel 模式 |

3. **决策流程图**（graphviz dot）：

```
并发问题？
├─ 需要保护共享状态？
│  ├─ 读多写少 → sync.RWMutex
│  ├─ 读写均等 → sync.Mutex
│  └─ 可以用通信代替？ → channel
├─ 需要 goroutine 间协调？
│  ├─ 传递数据 → channel
│  │  ├─ 多生产者 → fan-out
│  │  ├─ 多消费者 → fan-in
│  │  └─ 流水线 → pipeline
│  └─ 仅同步 → sync.WaitGroup / sync.Once
└─ 需要取消/超时？ → context
```

4. **代码示例**（✅/❌）：
   - channel 基本用法
   - fan-out/fan-in 模式
   - pipeline 模式
   - context 传播和取消
   - goroutine 生命周期管理

5. **反模式**：
   - goroutine 泄漏（未处理 context 取消）
   - 无缓冲 channel 死锁
   - 竞态条件
   - 忽略 context 传播

- [ ] **Step 2: 审查内容**
- [ ] **Step 3: Commit**

```bash
git add skills/golang-coder/references/concurrency.md
git commit -m "feat: add concurrency reference with decision flowchart"
```

---

### Task 8: 编写 testing.md

**Files:**
- Create: `skills/golang-coder/references/testing.md`

- [ ] **Step 1: 编写内容**

内容要点：
1. **核心原则**（3 条）：
   - 表驱动测试是首选模式
   - 测试应该快速、确定、无副作用
   - 测试公开行为，不测试实现细节

2. **速查规则**（表格）：

| 场景 | 做法 |
|------|------|
| 基本测试 | 表驱动 `[]struct{...}` + `t.Run` |
| 并行测试 | `t.Parallel()`（注意闭包变量） |
| 测试辅助 | `t.Helper()` 标记辅助函数 |
| 基准测试 | `BenchmarkXxx(b *testing.B)` |
| 设置/清理 | `t.Cleanup(func(){})` |
| mock 依赖 | 接口 + 手写 fake（优先于 mock 框架） |
| 子测试隔离 | 每个子测试独立 setup |

3. **代码示例**（✅/❌）：
   - 表驱动测试完整模板
   - `t.Parallel()` 正确使用（闭包捕获 `tt`）
   - benchmark 示例
   - fake 实现模式

4. **反模式**：
   - 测试间共享状态
   - 脆弱测试（依赖时间、随机数、外部服务）
   - 忽略 `t.Error` 返回值
   - 过度使用 mock 框架

- [ ] **Step 2: 审查内容**
- [ ] **Step 3: Commit**

```bash
git add skills/golang-coder/references/testing.md
git commit -m "feat: add testing reference"
```

---

### Task 9: 编写 performance.md

**Files:**
- Create: `skills/golang-coder/references/performance.md`

- [ ] **Step 1: 编写内容**

内容要点：
1. **核心原则**（3 条）：
   - 先 profile，再优化——不要猜测
   - 减少分配是最有效的优化手段
   - 值类型 vs 指针：理解逃逸分析

2. **速查规则**（表格）：

| 场景 | 做法 |
|------|------|
| 切片大小已知 | `make([]T, 0, n)` 预分配 |
| map 大小已知 | `make(map[K]V, n)` 预分配 |
| 字符串拼接 | `strings.Builder` |
| 对象复用 | `sync.Pool` |
| `[]byte`↔`string` | `unsafe` 转换（只读场景）或直接用 `[]byte` |
| 热路径日志 | 先检查日志级别再格式化 |

3. **代码示例**（✅/❌）：
   - 切片预分配
   - strings.Builder
   - sync.Pool 使用
   - pprof 基本使用流程

4. **反模式**：
   - 过早优化
   - 热路径日志（`log.Infof("result: %v", largeSlice)`）
   - 频繁 `[]byte`↔`string` 转换
   - 热路径中 `interface{}` 开销
   - 无缓冲 channel 阻塞

- [ ] **Step 2: 审查内容**
- [ ] **Step 3: Commit**

```bash
git add skills/golang-coder/references/performance.md
git commit -m "feat: add performance reference"
```

---

### Task 10: 最终验证和清理

**Files:**
- Verify: `skills/golang-corer/` 完整性
- Remove: `skills/go-naming/`（如尚未移除）

- [ ] **Step 1: 验证目录结构完整**

```bash
find skills/golang-coder -type f | sort
```

预期输出：
```
skills/golang-coder/SKILL.md
skills/golang-coder/references/concurrency.md
skills/golang-coder/references/error-handling.md
skills/golang-coder/references/interface-design.md
skills/golang-coder/references/naming.md
skills/golang-coder/references/performance.md
skills/golang-coder/references/project-structure.md
skills/golang-coder/references/testing.md
```

- [ ] **Step 2: 确认旧 skill 已移除**

```bash
ls skills/go-naming/ 2>/dev/null && echo "STILL EXISTS" || echo "REMOVED"
```

预期：`REMOVED`

- [ ] **Step 3: 验证每个文件格式正确**

检查所有 `.md` 文件：
- SKILL.md 有正确的 front matter（`---`、`name`、`description`）
- reference 文件有标题和内容
- 无 TBD/TODO 占位符

- [ ] **Step 4: 最终 commit**

```bash
git add -A
git commit -m "feat: complete golang-coder skill with all references and evolution mechanism"
```
