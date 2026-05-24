# 工作项A1: Go模块初始化

## 基本信息

| 属性 | 值 |
|------|-----|
| 编号 | A1 |
| 阶段 | Phase A - 基础设施 |
| 预估工期 | 1天 |
| 前置依赖 | 无 |
| 状态 | 未开始 |

## 工作描述

初始化Go模块（manga-translator）、创建Kratos项目目录结构（cmd/internal/pkg）、执行 `go mod init`、安装核心依赖（kratos/sqlite/langgraphgo/yaml）。

具体工作：
1. 执行 `go mod init github.com/user/manga-translator` 初始化Go模块
2. 按照后端方案中的项目结构创建完整目录树：
   - `cmd/server/` — 应用入口
   - `internal/conf/` — 配置定义
   - `internal/server/` — HTTP Server
   - `internal/service/` — 服务层
   - `internal/biz/` — 业务逻辑层
   - `internal/data/` — 数据访问层（含 `migration/` 子目录）
   - `internal/workflow/` — 工作流引擎
   - `pkg/llm/` — LLM适配层
   - `pkg/imageproc/` — 图片加工适配层
   - `pkg/embedding/` — Embedding服务
   - `pkg/vectorstore/` — 向量存储
   - `configs/` — 配置文件
3. 安装核心依赖：
   - `github.com/go-kratos/kratos/v2` — HTTP框架
   - `modernc.org/sqlite` — 纯Go SQLite驱动
   - `github.com/tmc/langchaingo` — LLM集成
   - `gopkg.in/yaml.v3` — YAML配置解析
4. 创建 `cmd/server/main.go` 骨架文件

## 相关文档

| 文档 | 说明 |
|------|------|
| [03-backend-design.md](../docs/03-backend-design.md) | 后端方案-项目结构与依赖选型 |

## 技术要点

- 模块名建议使用 `manga-translator`，与项目定位一致
- Kratos v2 使用 `github.com/go-kratos/kratos/v2` 导入路径
- `modernc.org/sqlite` 无需CGO，交叉编译友好
- `langchaingo` 提供LLM Provider抽象，支持OpenAI兼容接口
- 目录结构遵循Kratos推荐分层：cmd / internal / pkg / configs

## CheckList

- [ ] 创建go.mod（`go mod init`）
- [ ] 创建目录结构（cmd/internal/pkg/configs）
- [ ] 安装核心依赖（kratos/sqlite/langchaingo/yaml）
- [ ] 验证go build通过
- [ ] 创建.gitignore（排除data/、*.db、*.exe等）
