# 工作项B4: API接口定义

## 基本信息

| 属性 | 值 |
|------|-----|
| 编号 | B4 |
| 阶段 | Phase B - 数据模型 |
| 预估工期 | 1天 |
| 前置依赖 | A4 |
| 状态 | 未开始 |

## 工作描述

定义所有RESTful API的请求/响应Go struct、在HTTP Server注册路由表、实现空壳Handler（返回mock数据）。

具体工作：
1. 定义API请求/响应结构体（`internal/service/` 目录）：
   - **Project API**：`CreateProjectReq`/`CreateProjectResp`/`ListProjectsResp`/`GetProjectResp`/`DeleteProjectResp`/`GetChaptersResp`/`GetImagesResp`
   - **Translation API**：`StartTranslationReq`/`StartTranslationResp`/`ConfirmOCRReq`/`GetTranslationStatusResp`/`EditTranslationReq`/`ConfirmTranslationResp`/`RetranslateResp`/`BatchConfirmReq`
   - **Glossary API**：`CreateGlossaryReq`/`UpdateGlossaryReq`/`ListGlossaryResp`
   - **Knowledge API**：`CreateCharacterReq`/`UpdateCharacterReq`/`CreateStyleRuleReq`/`UpdateStyleRuleReq`/`CreateExampleReq`
   - **Settings API**：`GetSettingsResp`/`UpdateSettingsReq`
   - **Image API**：`ProcessImageReq`/`ProcessImageResp`/`ProcessStatusResp`
2. 在 `internal/server/http.go` 注册完整路由表
3. 实现空壳Handler（返回mock/空数据，后续工作项填充实际逻辑）

## 相关文档

| 文档 | 说明 |
|------|------|
| [03-backend-design.md](../docs/03-backend-design.md) | 后端方案-API接口列表与路由注册 |

## 技术要点

- 路由表参考03-backend-design.md的完整API列表：
  - 项目管理：`GET/POST /api/projects`、`GET/DELETE /api/projects/:id`、`GET /api/projects/:id/chapters`、`GET /api/chapters/:id/images`
  - 翻译工作流：`POST /api/translate/start`、`POST /api/translate/:runId/confirm-ocr`、`GET /api/images/:id/translations`、`PUT /api/translations/:id`、`POST /api/translations/:id/confirm`、`POST /api/translations/:id/retranslate`、`POST /api/translations/batch-confirm`
  - 术语库：`GET /api/projects/:id/glossary`、`POST /api/glossary/entries`、`PUT /api/glossary/entries/:id`、`DELETE /api/glossary/entries/:id`
  - 知识库：角色CRUD + 风格规则CRUD + 翻译范例
  - 图片加工：`POST /api/images/:id/process`、`GET /api/image-tasks/:taskId`
- 使用Kratos的 `r.Handle()` / `r.HandleFunc()` 注册路由
- 请求参数绑定使用Kratos的 `binding` 包
- 响应统一格式：`{"code": 0, "data": {...}, "message": ""}`

## CheckList

- [ ] 定义Project API结构体（Create/List/Get/Delete）
- [ ] 定义Translation API结构体（Extract/Batch/Confirm）
- [ ] 定义Glossary API结构体（CRUD）
- [ ] 定义Knowledge API结构体
- [ ] 定义Settings API结构体
- [ ] 注册所有路由
- [ ] 编写路由注册测试
