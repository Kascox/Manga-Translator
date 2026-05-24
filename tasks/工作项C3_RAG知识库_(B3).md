# 工作项C3: RAG知识库

## 基本信息

| 属性 | 值 |
|------|-----|
| 编号 | C3 |
| 阶段 | Phase C - 核心业务 |
| 预估工期 | 1天 |
| 前置依赖 | B3 |
| 状态 | 未开始 |

## 工作描述

实现Biz层知识库管理：chromem-go集合CRUD、TranslationExample存储与向量检索、CharacterProfile管理、StyleRule管理。

具体工作：
1. 实现 `internal/biz/knowledge.go` — 知识库管理业务逻辑
2. chromem-go DB初始化：按project隔离集合（集合名加 `project_{id}_` 前缀）
3. `AddTranslationExample`：
   - 创建TranslationExample数据库记录
   - 调用Embedding API生成向量
   - 将向量写入chromem-go索引
   - 同时将embedding写入SQLite的BLOB字段作为备份
4. `SearchSimilarExamples`：
   - 生成查询文本的embedding
   - 在chromem-go中按 `character_id` 过滤检索
   - 按余弦相似度排序，返回Top-K结果
   - 默认 top_k=5，similarity_threshold=0.7
5. `CharacterProfile` CRUD：管理角色风格档案
6. `StyleRule` CRUD：管理风格规则，含优先级排序（角色级>群体级>作品级）
7. 实现首次加载项目时从SQLite恢复chromem-go索引的逻辑

## 相关文档

| 文档 | 说明 |
|------|------|
| [06-knowledge-base.md](../docs/06-knowledge-base.md) | 知识库系统-四层结构与RAG检索机制 |

## 技术要点

- chromem-go集合名建议加项目前缀：`project_1_translation_examples`，避免多项目数据混淆
- Embedding生成需调用外部API（DeepSeek/OpenAI），需处理网络错误和重试
- 检索参数：top_k=5，similarity_threshold=0.7，character_filter=true
- StyleRule优先级排序：角色级(priority 20+) > 群体级(priority 10+) > 作品级(priority 0+)
- 同一 `rule_type` 的高优先级规则覆盖低优先级，不同 `rule_type` 累加
- SQLite备份恢复：首次打开项目时，从 `translation_examples.embedding` 字段恢复chromem-go索引
- 角色推断相关逻辑（character_signatures集合）在本工作项中仅初始化，实际推断在后续工作项实现

## CheckList

- [ ] 实现chromem-go DB初始化（按project隔离集合）
- [ ] 实现AddTranslationExample（含embedding生成）
- [ ] 实现SearchSimilarExamples（按character+相似度）
- [ ] 实现CharacterProfile CRUD
- [ ] 实现StyleRule CRUD（含优先级排序）
- [ ] 编写 `knowledge_biz_test.go`
