# 工作项D5: Checkpointer

## 基本信息

| 属性 | 值 |
|------|-----|
| 编号 | D5 |
| 阶段 | Phase D - 工作流 |
| 预估工期 | 1天 |
| 前置依赖 | A2 |
| 状态 | 未开始 |

## 工作描述

实现工作流状态持久化Checkpointer：基于SQLite存储工作流快照、支持恢复中断的工作流、清理过期快照。

具体工作：
1. 实现SQLiteCheckpointer，实现langgraphgo的Checkpointer接口
2. 实现Save：将工作流状态序列化为JSON存入SQLite，记录run_id、thread_id、当前节点名称、时间戳
3. 实现Load：从SQLite反序列化恢复工作流状态，支持按run_id或thread_id查询
4. 实现List：列出指定工作流（thread_id）的所有快照，按时间排序
5. 实现清理策略：保留最近N个快照（默认保留每个线程最近10个），定期清理超过30天的历史记录

## 相关文档

| 文档 | 说明 |
|------|------|
| [04-agent-workflow.md](../docs/04-agent-workflow.md) | SQLite Checkpointer持久化章节 |

## 技术要点

- SQLite表结构：run_id(TEXT)、thread_id(TEXT)、node_name(TEXT)、state(BLOB/JSON)、created_at(DATETIME)、updated_at(DATETIME)
- 序列化方式：TranslationState结构体→JSON编码→存入state字段
- 使用modernc.org/sqlite（纯Go，无CGO）操作SQLite
- 存储策略：中断点自动保存、节点执行完成后保存增量、工作流完成后保留最终状态
- 清理策略：可配置保留天数（默认30天）和每线程保留快照数（默认10），清理任务可定时执行或手动触发
- Checkpointer接口需匹配langgraphgo定义：`Put(ctx, threadID, checkpoint)` / `Get(ctx, threadID)` / `List(ctx, threadID)`

## CheckList

- [ ] 实现SQLiteCheckpointer(实现langgraphgo Checkpointer接口)
- [ ] 实现Save(序列化状态为JSON存入SQLite)
- [ ] 实现Load(反序列化恢复状态)
- [ ] 实现List(列出某工作流的所有快照)
- [ ] 实现清理策略(保留最近N个快照)
- [ ] 编写checkpointer_test.go
