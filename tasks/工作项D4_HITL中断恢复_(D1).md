# 工作项D4: HITL中断恢复

## 基本信息

| 属性 | 值 |
|------|-----|
| 编号 | D4 |
| 阶段 | Phase D - 工作流 |
| 预估工期 | 1天 |
| 前置依赖 | D1 |
| 状态 | 未开始 |

## 工作描述

实现Human-in-the-loop中断恢复机制：在用户确认角色/审核翻译结果等环节设置InterruptBefore、状态恢复续跑、用户编辑后注入修改。

具体工作：
1. 配置InterruptBefore节点列表：在`confirm_regions_node`（用户确认文字区域）、`confirm_character_node`（用户确认角色标注）、`review_translation_node`（用户审核翻译结果）前设置中断
2. 实现中断时状态返回前端：工作流暂停时将当前进度和待确认数据序列化，通过API返回给前端展示
3. 实现恢复接口：接收用户操作结果后Resume工作流，将用户输入注入状态继续执行
4. 实现UpdateState：用户编辑译文后更新工作流状态中的翻译结果字段
5. 处理并发中断场景：同一工作流不会同时存在多个中断点，确保状态一致性

## 相关文档

| 文档 | 说明 |
|------|------|
| [04-agent-workflow.md](../docs/04-agent-workflow.md) | HITL中断/恢复机制章节 |

## 技术要点

- langgraphgo的InterruptBefore配置：编译Graph时指定`graph.WithInterruptBefore("node_name")`选项
- 中断恢复流程：工作流执行到中断点→暂停→保存状态到Checkpointer→通知前端→用户操作→前端发送恢复请求→加载状态→注入用户输入→继续执行
- 中断状态数据结构：包含当前节点名称、工作流进度、待确认数据（如OCR结果/角色推断结果/翻译结果）
- UpdateState机制：用户编辑译文后调用`graph.UpdateState()`更新状态中对应字段，而非重新执行节点
- 并发控制：同一thread_id的工作流实例同时只有一个活跃执行，通过Checkpointer的thread_id隔离
- 中断通知方式：工作流暂停后通过API返回中断信息，前端可轮询或通过WebSocket/SSE接收通知

## CheckList

- [ ] 配置InterruptBefore节点列表(用户确认角色/用户审核翻译)
- [ ] 实现中断时状态返回前端(当前进度+待确认数据)
- [ ] 实现恢复接口(接收用户修改后Resume)
- [ ] 实现UpdateState(用户编辑译文后更新状态)
- [ ] 处理并发中断场景
- [ ] 编写hitl_test.go
