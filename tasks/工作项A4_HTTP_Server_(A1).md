# 工作项A4: HTTP Server

## 基本信息

| 属性 | 值 |
|------|-----|
| 编号 | A4 |
| 阶段 | Phase A - 基础设施 |
| 预估工期 | 1天 |
| 前置依赖 | A1 |
| 状态 | 未开始 |

## 工作描述

使用Kratos创建HTTP Server（`internal/server/http.go`）、配置CORS、注册健康检查路由（`/api/health`）、配置静态资源服务（`go:embed`）、路由分组（`/api/`前缀）。

具体工作：
1. 创建 `internal/server/http.go` — Kratos HTTP Server 配置
2. 配置监听 `localhost:8080`
3. 实现CORS中间件（开发环境允许 `localhost:5173`）
4. 注册健康检查路由 `GET /api/health`
5. 配置 `go:embed` 静态资源服务（`web/dist/` 目录）
6. 实现路由分组：`/api/` 前缀下的所有API路由
7. 配置请求超时：默认30s，翻译类接口120s
8. 配置请求体大小限制：50MB（图片上传场景）

## 相关文档

| 文档 | 说明 |
|------|------|
| [03-backend-design.md](../docs/03-backend-design.md) | 后端方案-HTTP Server设计与路由注册 |

## 技术要点

- Kratos HTTP Server 使用 `kratoshttp.NewServer()` 创建
- 静态资源通过 `//go:embed web/dist` 嵌入，使用 `http.FS()` 包装
- CORS配置使用 `github.com/go-kratos/kratos/v2/middleware/cors`
- 路由分组使用 `srv.Route("/api")` 创建子路由器
- 请求超时通过 `kratoshttp.Timeout()` 选项配置
- 前端SPA路由需配置fallback：非API请求返回 `index.html`
- 生产模式通过 `go:embed` 嵌入前端，开发模式代理到Vite Dev Server

## CheckList

- [ ] 创建HTTP Server配置（`internal/server/http.go`）
- [ ] 实现CORS中间件
- [ ] 注册健康检查路由
- [ ] 配置 `go:embed` 静态资源
- [ ] 编写Server启动/关闭测试
- [ ] 验证 `localhost:8080` 可访问
