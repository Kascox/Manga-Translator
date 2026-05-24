# 工作项A2: SQLite数据层

## 基本信息

| 属性 | 值 |
|------|-----|
| 编号 | A2 |
| 阶段 | Phase A - 基础设施 |
| 预估工期 | 1天 |
| 前置依赖 | A1 |
| 状态 | 未开始 |

## 工作描述

使用 `modernc.org/sqlite` 初始化数据库连接、创建Migration脚本（所有表）、实现Data层初始化（`internal/data/data.go`）、数据库文件存储路径管理。

具体工作：
1. 实现 `internal/data/data.go` — 数据库连接初始化与自动迁移
2. 编写Migration脚本：
   - `internal/data/migration/001_init.sql` — 核心表（projects/chapters/images/text_regions/translations）
   - `internal/data/migration/002_glossary.sql` — 术语库表（glossary_entry/glossary_translation）
   - `internal/data/migration/003_knowledge.sql` — 知识库表（character_profiles/style_rules/translation_examples）
3. 实现 `schema_migrations` 表管理与版本检查逻辑
4. 实现数据库文件路径配置（默认 `./data/comic-translator.db`，从config.yaml读取）
5. 实现自动迁移：启动时检查未执行的迁移脚本并按序执行

## 相关文档

| 文档 | 说明 |
|------|------|
| [03-backend-design.md](../docs/03-backend-design.md) | 后端方案-Data层设计与数据库配置 |
| [05-database-schema.md](../docs/05-database-schema.md) | 数据库结构-9张表DDL与索引设计 |

## 技术要点

- 使用 `modernc.org/sqlite` 驱动，连接字符串格式：`file:./data/comic-translator.db?_journal_mode=WAL`
- 每个Migration脚本在单个事务内执行，失败则回滚
- Migration脚本只允许 `CREATE/ALTER`，不允许 `DROP`（保证向前兼容）
- 所有9张表：projects、chapters、images、text_regions、translations、glossary_entry、glossary_translation、character_profiles、style_rules、translation_examples
- 需创建索引（参考05-database-schema.md索引设计章节）
- 启用WAL模式提升并发读性能
- 外键约束需启用 `PRAGMA foreign_keys = ON`

## CheckList

- [ ] 实现DB连接初始化（`internal/data/data.go`）
- [ ] 编写Migration脚本（001_init.sql / 002_glossary.sql / 003_knowledge.sql）
- [ ] 实现自动迁移逻辑（schema_migrations版本检查+按序执行）
- [ ] 编写 `data_test.go` 验证表创建
- [ ] 验证数据库文件创建和读写
