# golang-coder Skill 设计文档

## 概述

创建一个面向有经验的 Go 开发者的 opencode skill，覆盖 Go 编码全生命周期。采用混合型架构：主入口轻量决策引导 + 模块化 reference 文件（速查手册格式），在接口设计和并发领域加入决策流程图。合并并吸收现有 `go-naming` skill。

参考来源：Go 标准库风格、Google Go Style Guide、Uber Go Style Guide。

## 文件结构

```
skills/
├── golang-coder/                        # 新 skill
│   ├── SKILL.md                         # 主入口
│   └── references/
│       ├── naming.md                    # 从 go-naming 吸收
│       ├── project-structure.md
│       ├── error-handling.md
│       ├── interface-design.md
│       ├── concurrency.md
│       ├── testing.md
│       └── performance.md
├── go-naming/                           # 旧 skill，完成后移除
│   └── SKILL.md
docs/
└── golang-coder-evolution.md            # 进化日志（skill 不加载）
```

## SKILL.md 主入口

### 触发条件

- 正在创建或编辑 `.go` 文件
- 讨论 Go 包结构、接口设计、并发模式
- 编写 Go 测试或 benchmark
- 排查 Go 性能问题或竞态条件
- 设计 Go API 或错误处理策略

### 快速导航

| 任务 | 参考文档 |
|------|----------|
| 命名任何 Go 标识符 | references/naming.md |
| 组织项目或包结构 | references/project-structure.md |
| 处理错误、设计错误策略 | references/error-handling.md |
| 设计接口、选择抽象方式 | references/interface-design.md |
| 使用 goroutine/channel/sync | references/concurrency.md |
| 编写测试或 benchmark | references/testing.md |
| 性能优化、减少分配 | references/performance.md |

### 决策引导

主入口提供轻量决策点，引导到正确的 reference：

1. **需要抽象？** → references/interface-design.md（含决策流程图）
2. **并发方案？** → references/concurrency.md（含决策流程图）
3. **其他** → 直接查对应 reference

### 全局原则

- 简洁优于复杂，组合优于继承
- 优先使用标准库，谨慎引入第三方依赖
- 错误必须处理，永远不要忽略
- 先写正确的代码，再考虑性能

### 进化机制

SKILL.md 中声明自动记录规则：

- 编码过程中发现新模式/反模式时，自动追加到 `docs/golang-coder-evolution.md`
- 不主动读取此文件，仅在写入时打开
- 如果 `docs/` 目录不存在，自动创建
- 用户说「整理进化记录」时，读取文件，将待整理条目融入对应 reference

## Reference 文件详细规格

### 统一格式

每个 reference 文件遵循：

1. **核心原则**（2-3 条，简短）
2. **速查规则**（表格形式）
3. **代码示例**（✅ good / ❌ bad 对比）
4. **反模式**（常见错误 + 修正）
5. **决策流程图**（仅在需要的领域：interface-design、concurrency）

每个文件预计 150-250 行。

### 1. naming.md

来源：直接从现有 `skills/go-naming/SKILL.md` 吸收，内容不变。

包含：速查表、包/类型/接口/函数/变量/常量/接收器/错误命名、缩写词表、决策流程图。

### 2. project-structure.md

- 标准目录布局：`cmd/`、`internal/`、`pkg/`、`api/`、`configs/`
- 包拆分原则：按领域 vs 按层，何时用 `internal`
- 依赖管理：`go mod` 最佳实践、版本选择
- 文件组织：按职责命名、避免 `utils`/`common`
- 反模式：上帝包、循环依赖、过深嵌套

### 3. error-handling.md

- 错误包装：`fmt.Errorf("%w", err)` + `errors.Is/As`
- 哨兵错误 vs 自定义错误类型的选择
- 错误边界：何时处理、何时传递、何时 wrap
- panic 的使用边界
- 反模式：`_` 忽略错误、丢失上下文、过度嵌套

### 4. interface-design.md

- 小接口哲学：单方法接口、`io.Reader` 模式
- 接口定义在使用方
- 组合优于继承：embed 接口
- 隐式实现的优势
- 决策流程图：需要抽象时 → 接口 vs 泛型 vs 函数类型 vs 简单回调
- 反模式：大接口、在实现方定义接口、返回接口

### 5. concurrency.md

- goroutine 生命周期管理
- channel 模式：fan-in/fan-out、pipeline、cancel、timeout
- sync 原语选择：Mutex vs RWMutex vs WaitGroup vs Once vs Pool
- context 传播和取消
- 决策流程图：共享状态 → 用锁 vs channel；并发模式 → 选哪种 channel 模式
- 反模式：goroutine 泄漏、无缓冲 channel 死锁、竞态条件

### 6. testing.md

- 表驱动测试模式
- 测试组织：单元/集成/benchmark 分层
- `t.Helper()`、`t.Parallel()`、test fixtures
- 子测试和测试隔离
- mock 策略：接口 mock vs fakes vs 生产替代
- 反模式：测试间耦合、脆弱测试、忽略 error 检查

### 7. performance.md

- 减少分配：切片预分配、字符串拼接、`sync.Pool`
- 逃逸分析：值类型 vs 指针的选择
- profiling：`pprof` 使用流程
- 常见热点：`[]byte`↔`string` 转换、`interface{}` 开销
- 反模式：过早优化、热路径日志、无缓冲 channel

## 进化机制

### evolution 文件

位置：`docs/golang-coder-evolution.md`（不在 skills/ 目录内，不被 skill 加载）

### 条目格式

```markdown
### [待整理]

- **[日期] 领域-简述**
  - 场景：描述遇到的具体场景
  - 建议：应该补充什么规则/示例到哪个 reference
  - 状态：待整理

### [已整理]

- **[日期] 领域-简述** → 整理至 xxx.md
```

### 工作流程

1. **AI 自动写入**：编码中发现新模式/反模式，AI 自行判断后追加到 evolution 文件（只写不读）
2. **用户触发整理**：用户说「整理进化记录」时，AI 读取文件，批量整理进对应 reference
3. **整理后更新**：已整理条目移入「已整理」区域

### 自动创建

如果 `docs/` 目录不存在，写入时自动创建。

## 实施步骤

1. 创建 `skills/golang-coder/` 目录结构
2. 编写 `SKILL.md` 主入口
3. 将 `go-naming/SKILL.md` 内容迁移到 `references/naming.md`（调整 front matter）
4. 逐个编写其余 6 个 reference 文件
5. 创建初始的 `docs/golang-coder-evolution.md`
6. 更新 AGENTS.md（如有必要）
7. 移除旧 `skills/go-naming/` 目录
