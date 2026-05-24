# 工作项A5: webview集成

## 基本信息

| 属性 | 值 |
|------|-----|
| 编号 | A5 |
| 阶段 | Phase A - 基础设施 |
| 预估工期 | 1天 |
| 前置依赖 | A4 |
| 状态 | 未开始 |

## 工作描述

集成 `go-webview2`、启动时自动打开桌面窗口加载 `localhost:8080`、窗口标题/尺寸配置、进程生命周期管理（webview关闭时退出HTTP Server）。

具体工作：
1. 安装 `go-webview2` 依赖
2. 实现 `cmd/server/main.go` 中的webview窗口创建
3. 配置窗口标题（默认"漫画翻译工具"）和默认尺寸（1280x800）
4. 实现Server+webview协同启动/关闭：
   - 先启动Kratos HTTP Server
   - Server就绪后创建WebView窗口
   - WebView窗口关闭时触发Server优雅关闭
5. 实现WebView2运行时检测：若未安装则提示用户安装
6. 开发模式支持：环境变量 `DEV_MODE=1` 时跳过WebView，直接浏览器访问

## 相关文档

| 文档 | 说明 |
|------|------|
| [03-backend-design.md](../docs/03-backend-design.md) | 后端方案-webview桌面壳集成 |

## 技术要点

- `go-webview2` 仅Windows可用，编译需加 `CGO_ENABLED=1`（但SQLite用modernc无需CGO）
- WebView2 Runtime 在 Windows 10/11 默认已安装，需做运行时检测
- 窗口关闭回调中调用 `server.Stop()` 实现优雅关闭
- 开发模式下可禁用WebView，直接在浏览器访问 `localhost:5173`（Vite Dev Server）
- 窗口标题栏使用原生样式，内容区全屏加载前端
- 窗口尺寸从 `config.yaml` 的 `server.window` 段读取

## CheckList

- [ ] 安装go-webview2依赖
- [ ] 实现webview窗口创建
- [ ] 配置窗口标题和默认尺寸
- [ ] 实现Server+webview协同启动/关闭
- [ ] WebView2运行时检测与提示
- [ ] 手动测试桌面窗口
