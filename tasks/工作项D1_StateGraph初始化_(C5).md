# 工作项D1: StateGraph初始化

## 基本信息

| 属性 | 值 |
|------|-----|
| 编号 | D1 |
| 阶段 | Phase D - 工作流 |
| 预估工期 | 1天 |
| 前置依赖 | C5 |
| 状态 | 未开始 |

## 工作描述

使用langgraphgo创建翻译工作流StateGraph(`internal/workflow/engine.go`)：定义状态结构、注册所有节点、定义边和条件路由、编译Graph。

具体工作：
1. 定义`TranslationState`结构体，包含工作流执行过程中需要的全部状态字段（图片ID、TextRegion列表、角色标注结果、术语匹配结果、翻译结果、质量校验结果等）
2. 创建`StateGraph`实例，传入状态结构定义
3. 注册所有工作流节点（OCR/术语匹配/RAG/翻译/校验/图片加工等，占位实现）
4. 定义节点间的边连接，包括START→OCR、OCR→确认区域等直线边
5. 定义条件边，特别是质量校验节点根据通过/不通过分别路由到不同节点
6. 编译Graph为Runnable，供后续工作流执行调用

## 相关文档

| 文档 | 说明 |
|------|------|
| [04-agent-workflow.md](../docs/04-agent-workflow.md) | StateGraph设计、节点定义、边定义 |

## 技术要点

- langgraphgo的StateGraph API使用方式：`graph.NewStateGraph(stateSchema, ...)` 创建实例
- 节点注册：`graph.AddNode("node_name", nodeFunc)` 注册节点函数
- 边定义：`graph.AddEdge("from", "to")` 定义确定边，`graph.AddConditionalEdge("from", conditionFunc, map)` 定义条件边
- 状态结构设计需覆盖所有节点的输入输出，使用可变字段（指针/切片）表示可选状态
- 编译：`graph.Compile()` 生成可执行的Runnable
- 条件边中质量校验的路由逻辑：通过→用户审核，不通过→重试LLM翻译

## CheckList

- [ ] 定义TranslationState结构体
- [ ] 创建StateGraph实例
- [ ] 注册OCR/术语匹配/RAG/翻译/校验/图片加工节点(占位)
- [ ] 定义节点间的边连接
- [ ] 定义条件边(质量校验通过/不通过)
- [ ] 编译Runnable
- [ ] 编写engine_test.go(图构建验证)
