# 工作项C4: Embedding服务

## 基本信息

| 属性 | 值 |
|------|-----|
| 编号 | C4 |
| 阶段 | Phase C - 核心业务 |
| 预估工期 | 1天 |
| 前置依赖 | C3 |
| 状态 | 未开始 |

## 工作描述

实现Embedding生成服务（`pkg/embedding/`）：调用DeepSeek embedding API、批量Embedding生成、限流控制、缓存策略。

具体工作：
1. 定义 `EmbeddingProvider` 接口（`pkg/embedding/embedding.go`）：
   - `Embed(ctx, text string) ([]float32, error)` — 单文本embedding
   - `EmbedBatch(ctx, texts []string) ([][]float32, error)` — 批量embedding
2. 实现DeepSeek Embedding调用：
   - 基于langchaingo的Embedding接口
   - 支持OpenAI兼容接口（DeepSeek使用OpenAI兼容API）
   - 配置API Key、Base URL、模型名
3. 实现批量Embedding：
   - 分批处理（每批最多16条，避免API限制）
   - 批间等待避免触发限流
4. 实现API限流（rate limiter）：
   - 使用 `golang.org/x/time/rate` 实现
   - 默认每秒2次请求
5. 实现结果缓存：
   - 使用内存map缓存（文本hash → embedding向量）
   - 避免重复计算相同文本的embedding

## 相关文档

| 文档 | 说明 |
|------|------|
| [04-agent-workflow.md](../docs/04-agent-workflow.md) | 工作流设计-Embedding使用场景 |
| [06-knowledge-base.md](../docs/06-knowledge-base.md) | 知识库系统-Embedding生成策略 |

## 技术要点

- DeepSeek Embedding API兼容OpenAI格式，使用langchaingo的OpenAI Embedding适配器
- Embedding维度取决于模型：DeepSeek 768维 / text-embedding-3-small 1536维
- 文本预处理：去除特殊字符，统一空白符
- 向量归一化：L2归一化，确保余弦相似度有效
- 限流控制：使用 `golang.org/x/time/rate` 令牌桶算法
- 缓存策略：使用 `sync.Map` 实现并发安全缓存，key为文本SHA256哈希
- 错误处理：API调用失败时实现指数退避重试（最多3次）

## CheckList

- [ ] 定义EmbeddingProvider接口
- [ ] 实现DeepSeek Embedding调用
- [ ] 实现批量Embedding（分批处理）
- [ ] 实现API限流（rate limiter）
- [ ] 实现结果缓存（避免重复计算）
- [ ] 编写 `embedding_test.go`（Mock API）
