# 工作项C5: LLM适配层

## 基本信息

| 属性 | 值 |
|------|-----|
| 编号 | C5 |
| 阶段 | Phase C - 核心业务 |
| 预估工期 | 1天 |
| 前置依赖 | A3 |
| 状态 | 未开始 |

## 工作描述

实现统一LLM Provider接口（`pkg/llm/`）：基于langchaingo封装、DeepSeek实现、Prompt模板管理、流式响应支持。

具体工作：
1. 定义 `LLMProvider` 接口（`pkg/llm/provider.go`）：
   - `Chat(ctx, messages []Message) (*Response, error)` — 同步对话
   - `ChatStream(ctx, messages []Message) (<-chan StreamChunk, error)` — 流式对话
   - `SetModel(model string)` — 切换模型
2. 实现 `Message` / `Response` / `StreamChunk` 结构体
3. 实现DeepSeek Provider（`pkg/llm/deepseek.go`）：
   - 基于langchaingo的OpenAI兼容接口（DeepSeek API兼容OpenAI格式）
   - 从 `config.yaml` 读取 API Key、Base URL、模型名
   - 支持system/user/assistant消息角色
4. 实现Prompt模板加载与渲染（`pkg/llm/prompt.go`）：
   - 支持从 `configs/prompts/` 目录加载模板文件
   - 模板语法使用Go `text/template`
   - 支持变量替换（角色风格、术语、范例等）
5. 实现流式响应处理：
   - 解析SSE事件流
   - 逐chunk回调或channel输出

## 相关文档

| 文档 | 说明 |
|------|------|
| [04-agent-workflow.md](../docs/04-agent-workflow.md) | 工作流设计-LLM适配层章节 |
| [03-backend-design.md](../docs/03-backend-design.md) | 后端方案-配置管理与LLM配置 |

## 技术要点

- langchaingo的OpenAI兼容接口：`langchaingo/llms/openai`，通过 `WithBaseURL()` 指定DeepSeek端点
- DeepSeek API配置：`base_url: "https://api.deepseek.com"`，`model: "deepseek-chat"`
- API Key通过环境变量间接读取：`api_key_env: "DEEPSEEK_API_KEY"`
- 流式响应使用langchaingo的 `llms.WithStreamingFunc()` 选项
- Prompt模板示例：
  - `system_translate.tmpl` — 翻译系统Prompt
  - `system_quality_check.tmpl` — 质量校验Prompt
  - `few_shot_example.tmpl` — Few-shot示例格式
- 模板渲染时注入：角色风格描述、风格规则列表、术语对照表、翻译范例

## CheckList

- [ ] 定义LLMProvider接口（Chat/ChatStream）
- [ ] 实现DeepSeek Provider（OpenAI兼容接口）
- [ ] 实现Prompt模板加载与渲染
- [ ] 实现流式响应处理
- [ ] 配置从config.yaml读取API Key/model
- [ ] 编写 `provider_test.go`（Mock）
