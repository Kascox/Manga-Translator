# Manga-Translator

基于AI的漫画翻译工具，支持OCR文字识别、智能翻译、术语管理和图片加工，实现漫画翻译全流程自动化。

## 技术栈

| 层级 | 技术 |
|------|------|
| 桌面壳 | go-webview2 |
| 前端 | React 18 + TypeScript + Ant Design 5 |
| 后端 | Go + Kratos v2 |
| 工作流 | langgraphgo |
| LLM集成 | langchaingo |
| 结构化存储 | modernc.org/sqlite |
| 向量存储 | chromem-go |

## 项目结构

    cmd/server/          # 应用入口
    internal/            # 内部模块
      ├── conf/          # 配置定义
      ├── server/        # HTTP Server
      ├── service/       # 服务层
      ├── biz/           # 业务逻辑层
      ├── data/          # 数据访问层
      └── workflow/      # 工作流引擎
    pkg/                 # 可复用包
      ├── llm/           # LLM适配层
      ├── imageproc/     # 图片加工适配层
      ├── embedding/     # Embedding服务
      └── vectorstore/   # 向量存储
    configs/             # 配置文件
    web/                 # 前端项目（Vite + React）

## 开发环境

- Go 1.22+
- Node.js 18+
- pnpm

## 快速开始

```bash
# 启动后端
go run ./cmd/server/

# 启动前端开发服务器
cd web && pnpm dev
```