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
