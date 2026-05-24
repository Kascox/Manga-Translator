# 工作项B3: RAG模型

## 基本信息

| 属性 | 值 |
|------|-----|
| 编号 | B3 |
| 阶段 | Phase B - 数据模型 |
| 预估工期 | 1天 |
| 前置依赖 | A2 |
| 状态 | 未开始 |

## 工作描述

实现CharacterProfile/StyleRule/TranslationExample的Go结构体、chromem-go向量集合初始化、Repository接口定义。

具体工作：
1. 定义知识库数据模型Go结构体：
   - `CharacterProfile` — id/project_id/name/group_name/pronoun_preference/sentence_style/tone_words/catchphrase/vocab_level/error_pattern/description/created_at/updated_at
   - `StyleRule` — id/project_id/scope/target_id/rule_type/rule_content/priority/created_at/updated_at
   - `TranslationExample` — id/project_id/character_id/source_text/translated_text/context/embedding/created_at
2. 定义Repository接口：
   - `CharacterRepo` — ListByProject/Get/Create/Update/Delete
   - `StyleRuleRepo` — ListByProject(scope过滤)/Create/Update/Delete
   - `TranslationExampleRepo` — ListByCharacter/Create/Delete
3. 初始化chromem-go向量集合：
   - `translation_examples` 集合 — RAG检索相似翻译
   - `character_signatures` 集合 — 角色自动推断
4. 实现chromem-go DB初始化与集合创建逻辑

## 相关文档

| 文档 | 说明 |
|------|------|
| [05-database-schema.md](../docs/05-database-schema.md) | 数据库结构-知识库3张表与向量存储设计 |
| [06-knowledge-base.md](../docs/06-knowledge-base.md) | 知识库系统-四层结构与RAG检索机制 |

## 技术要点

- chromem-go使用 `chromem.NewDB()` 创建，持久化到 `./data/chromem` 目录
- `translation_examples` 集合：内容字段为 `source_text`，元数据含 `project_id`/`character_id`/`translated_text`
- `character_signatures` 集合：内容字段为角色口头禅+语气词组合，元数据含 `project_id`/`group_name`
- Embedding维度取决于模型（DeepSeek 768维 / OpenAI 1536维），需配置化
- `TranslationExample.embedding` 为BLOB字段，存储embedding向量备份（用于从SQLite恢复chromem-go索引）
- `StyleRule.scope` 支持 project/group/character 三种作用域
- `StyleRule.target_id` 为可空字段，scope=project时为NULL

## CheckList

- [ ] 定义CharacterProfile结构体和Repository
- [ ] 定义StyleRule结构体和Repository
- [ ] 定义TranslationExample结构体和Repository
- [ ] 初始化chromem-go集合（translation_examples/character_signatures）
- [ ] 编写 `rag_repo_test.go`
