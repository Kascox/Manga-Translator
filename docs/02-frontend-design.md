# 前端方案

## 技术选型

| 技术 | 版本 | 用途 |
|------|------|------|
| React | 18 | 组件化 UI 开发 |
| TypeScript | 5.x | 类型安全 |
| Ant Design | 5 | 企业级 UI 组件库 |
| Fabric.js | 6.x | Canvas 图形库，文字区域选择/拖拽/缩放 |
| Zustand | 4.x | 轻量级状态管理 |
| Vite | 5.x | 构建工具，快速热更新 |
| pnpm | 8.x | 包管理 |

### 选型理由

- **React 18**：成熟生态，Concurrent Mode 支持流畅交互，社区资源丰富
- **Ant Design 5**：开箱即用的企业级组件（Tree、Table、Modal、Form），CSS-in-JS 定制能力强
- **Fabric.js**：专为 Canvas 交互设计，支持对象选择、拖拽、缩放、分组，适合图片标注场景
- **Zustand**：比 Redux 更轻量，无 boilerplate，TypeScript 友好，适合中等复杂度状态管理

## 项目结构

```
web/
├── index.html                  # HTML 入口
├── package.json
├── tsconfig.json
├── tsconfig.node.json
├── vite.config.ts
└── src/
    ├── main.tsx                # 应用入口
    ├── index.css               # 全局样式
    ├── vite-env.d.ts           # Vite 类型声明
    ├── App.tsx                 # 根组件-三栏布局
    ├── components/             # UI 组件
    │   ├── ImageCanvas/        # 图片Canvas组件
    │   │   ├── index.tsx
    │   │   ├── ImageCanvas.tsx       # 主Canvas组件
    │   │   ├── ToolBar.tsx           # 工具栏-框选/套索/缩放
    │   │   ├── RegionOverlay.tsx     # 文字区域覆盖层
    │   │   └── useCanvasInteraction.ts  # Canvas交互Hook
    │   ├── WorkspaceTree/      # 工作区树组件
    │   │   ├── index.tsx
    │   │   └── WorkspaceTree.tsx
    │   ├── TranslationPanel/   # 翻译管理面板
    │   │   ├── index.tsx
    │   │   ├── TranslationPanel.tsx   # 主面板
    │   │   ├── TranslationItem.tsx    # 单条翻译项
    │   │   ├── CharacterTag.tsx       # 角色Tag组件
    │   │   └── ContextMenu.tsx        # 右键菜单
    │   ├── SettingsModal/      # 设置弹窗
    │   │   ├── index.tsx
    │   │   └── SettingsModal.tsx
    │   ├── CharacterModal/     # 角色管理弹窗
    │   │   ├── index.tsx
    │   │   └── CharacterModal.tsx
    │   ├── StyleRuleModal/     # 风格规则弹窗
    │   │   ├── index.tsx
    │   │   └── StyleRuleModal.tsx
    │   └── GlossaryModal/      # 术语编辑弹窗
    │       ├── index.tsx
    │       └── GlossaryModal.tsx
    ├── stores/                 # Zustand 状态管理
    │   ├── useProjectStore.ts  # 项目状态
    │   ├── useCanvasStore.ts   # Canvas状态
    │   ├── useTranslationStore.ts  # 翻译状态
    │   └── useSettingsStore.ts # 设置状态
    ├── services/               # API 服务层
    │   ├── api.ts              # Axios 实例配置
    │   ├── projectService.ts   # 项目管理API
    │   ├── translationService.ts  # 翻译工作流API
    │   ├── glossaryService.ts  # 术语库API
    │   ├── knowledgeService.ts # 知识库API
    │   └── imageService.ts     # 图片加工API
    ├── types/                  # TypeScript 类型定义
    │   ├── project.ts
    │   ├── translation.ts
    │   ├── glossary.ts
    │   └── knowledge.ts
    ├── hooks/                  # 自定义 Hooks
    │   ├── useWebSocket.ts     # 工作流状态推送
    │   └── useContextMenu.ts   # 右键菜单Hook
    └── utils/                  # 工具函数
        ├── canvasUtils.ts      # Canvas坐标变换
        └── formatUtils.ts      # 格式化工具
```

## 组件设计

### ImageCanvas 组件

Fabric.js 集成核心组件，负责图片渲染和文字区域标注交互。

**功能清单：**

| 功能 | 说明 |
|------|------|
| 图片加载与渲染 | 加载原图到 Fabric.js Canvas |
| 框选工具 | 矩形框选文字区域，创建 TextRegion |
| 套索工具 | 自由形状圈选文字区域 |
| 区域拖拽/缩放 | 选中区域后支持拖拽移动和缩放调整 |
| Canvas 缩放平移 | 鼠标滚轮缩放，中键拖拽平移 |
| 区域高亮 | 鼠标悬停高亮显示文字区域 |
| 右键菜单 | 区域上右键弹出操作菜单 |

**Fabric.js 集成要点：**

```typescript
// Canvas 初始化与清理
// - useEffect 中初始化 Fabric.Canvas
// - 组件卸载时 dispose 释放资源
// - 监听 window resize 自适应 Canvas 尺寸

// 文字区域对象模型
// - 每个文字区域对应一个 Fabric.Rect（框选模式）或 Fabric.Path（套索模式）
// - 区域对象携带 metadata: { regionId, status, isIgnored }
// - 选中区域时显示控制手柄（缩放/旋转）
```

**工具栏设计：**

| 工具 | 图标 | 快捷键 | 说明 |
|------|------|--------|------|
| 选择 | pointer | V | 默认模式，选中/拖拽/缩放区域 |
| 框选 | rect-select | R | 矩形框选创建新文字区域 |
| 套索 | lasso-select | L | 自由形状圈选创建新文字区域 |
| 缩放+ | zoom-in | Ctrl+= | 放大 Canvas |
| 缩放- | zoom-out | Ctrl+- | 缩小 Canvas |
| 适应 | fit | Ctrl+0 | 适应窗口大小 |

### WorkspaceTree 组件

基于 Ant Design Tree 组件的工作区文件导航。

**功能清单：**

| 功能 | 说明 |
|------|------|
| 三级树结构 | 项目 → 章节 → 图片 |
| 拖拽导入 | 拖拽图片文件到章节节点导入 |
| 右键菜单 | 新建章节、删除、重命名等 |
| 状态图标 | 图片节点显示翻译状态图标 |
| 选中联动 | 选中图片节点 → 加载到 Canvas + 翻译面板 |

**数据结构：**

```typescript
interface TreeNode {
  key: string;
  title: string;
  type: 'project' | 'chapter' | 'image';
  status?: 'pending' | 'translating' | 'completed';
  children?: TreeNode[];
}
```

### TranslationPanel 组件

翻译管理核心面板，展示原文/译文对照列表。

**功能清单：**

| 功能 | 说明 |
|------|------|
| 翻译列表 | 原文/译文对照，按文字区域顺序排列 |
| 角色Tag | 每条翻译前显示角色Tag |
| 行内编辑 | 点击译文进入编辑模式 |
| 右键菜单 | 加入术语库、重新翻译、编辑角色等 |
| 状态标记 | 待翻译/翻译中/待审核/已确认等状态 |
| 批量操作 | 全选、批量确认、批量重新翻译 |

**TranslationItem 子组件：**

```typescript
interface TranslationItemProps {
  id: number;
  originalText: string;
  translatedText: string;
  status: 'pending' | 'translating' | 'reviewing' | 'modified' | 'confirmed';
  character?: {
    id: number;
    name: string;
    color: string;
  };
  onEdit: (text: string) => void;
  onConfirm: () => void;
  onContextMenu: (e: React.MouseEvent) => void;
}
```

**CharacterTag 子组件：**

- 显示角色名称，带颜色标记
- 点击弹出角色选择/创建浮层
- 自动推断结果用虚线框标识，手动指定用实线框
- 未知角色显示"?"，点击可选择或创建

### SettingsModal 组件

基于 Ant Design Modal + Form 的设置弹窗。

**设置项分组：**

| 分组 | 设置项 |
|------|--------|
| API 配置 | DeepSeek API Key、GPT-image API Key、模型选择 |
| 翻译设置 | 源语言、目标语言、自动/手动模式 |
| 知识库 | 角色档案开关、风格规则开关、翻译范例RAG开关、术语库开关 |
| 输出设置 | 图片格式、质量、输出目录 |

### 角色管理弹窗

管理当前项目的角色档案。

**功能清单：**

| 功能 | 说明 |
|------|------|
| 角色列表 | 展示项目所有角色，按出场频率排序 |
| 新建角色 | 填写角色名称、群体、风格属性 |
| 编辑角色 | 修改人称代词、句式特点、语气词、口头禅等 |
| 合并角色 | 将误拆分的角色合并 |
| 删除角色 | 删除角色并重新分配翻译条目 |

### 风格规则弹窗

管理项目级/群体级/角色级风格规则。

**功能清单：**

| 功能 | 说明 |
|------|--------|
| 规则列表 | 按作用域分组展示（作品级/群体级/角色级） |
| 新建规则 | 选择作用域，填写规则内容 |
| 优先级调整 | 同一作用域内调整规则优先级 |
| 规则预览 | 预览规则在 Prompt 中的组装效果 |

## 状态管理设计

使用 Zustand 管理全局状态，按职责拆分为多个 Store。

### useProjectStore

```typescript
interface ProjectState {
  // 数据
  projects: Project[];
  currentProject: Project | null;
  chapters: Chapter[];
  currentChapter: Chapter | null;
  images: Image[];
  currentImage: Image | null;

  // 操作
  loadProjects: () => Promise<void>;
  createProject: (name: string, path: string) => Promise<Project>;
  selectProject: (id: number) => void;
  selectChapter: (id: number) => void;
  selectImage: (id: number) => void;
}
```

### useCanvasStore

```typescript
interface CanvasState {
  // 数据
  tool: 'select' | 'rect-select' | 'lasso';
  zoom: number;
  panOffset: { x: number; y: number };
  selectedRegionId: number | null;
  regions: TextRegion[];

  // 操作
  setTool: (tool: CanvasState['tool']) => void;
  setZoom: (zoom: number) => void;
  selectRegion: (id: number | null) => void;
  addRegion: (region: TextRegion) => void;
  updateRegion: (id: number, updates: Partial<TextRegion>) => void;
  removeRegion: (id: number) => void;
}
```

### useTranslationStore

```typescript
interface TranslationState {
  // 数据
  translations: Translation[];
  loading: boolean;
  workflowStatus: 'idle' | 'running' | 'paused' | 'completed';

  // 操作
  startTranslation: (imageId: number, mode: 'auto' | 'manual') => Promise<void>;
  confirmTranslation: (id: number) => Promise<void>;
  editTranslation: (id: number, text: string) => Promise<void>;
  retranslate: (id: number) => Promise<void>;
  batchConfirm: (ids: number[]) => Promise<void>;

  // 工作流回调
  onOCRComplete: (regions: TextRegion[]) => void;
  onTranslationComplete: (translations: Translation[]) => void;
  onWorkflowComplete: () => void;
}
```

### useSettingsStore

```typescript
interface SettingsState {
  // 数据
  apiKeys: { deepseek?: string; openai?: string };
  sourceLang: string;
  targetLang: string;
  knowledgeBase: {
    characterProfiles: boolean;
    styleRules: boolean;
    translationExamples: boolean;
    glossary: boolean;
  };

  // 操作
  updateSettings: (updates: Partial<SettingsState>) => Promise<void>;
  loadSettings: () => Promise<void>;
}
```

## API 服务层设计

### Axios 实例配置

```typescript
// services/api.ts
const api = axios.create({
  baseURL: 'http://localhost:8080/api',
  timeout: 30000,
  headers: { 'Content-Type': 'application/json' },
});

// 请求拦截器：附加认证信息
// 响应拦截器：统一错误处理
```

### API 服务模块

| 模块 | 文件 | 主要方法 |
|------|------|----------|
| 项目管理 | projectService.ts | listProjects, createProject, getProject, deleteProject |
| 翻译工作流 | translationService.ts | startTranslation, confirmOCR, editTranslation, confirmTranslation, retranslate |
| 术语库 | glossaryService.ts | listEntries, createEntry, updateEntry, deleteEntry |
| 知识库 | knowledgeService.ts | listCharacters, createCharacter, updateCharacter, listStyleRules, createStyleRule |
| 图片加工 | imageService.ts | processImage, getProcessStatus |

## 与后端联调方案

### 通信方式

| 场景 | 方式 | 说明 |
|------|------|------|
| 常规 CRUD | HTTP REST | 请求-响应模式 |
| 工作流状态推送 | WebSocket / SSE | 工作流节点状态变更实时推送 |
| 图片传输 | HTTP + Base64 / FormData | 图片上传与下载 |

### 联调流程

1. **开发环境**：Vite dev server (port 5173) → Kratos HTTP Server (port 8080)，通过 Vite proxy 转发 API 请求
2. **生产环境**：Kratos 通过 `go:embed` 嵌入前端静态资源，统一端口提供服务
3. **桌面壳**：go-webview2 加载 `http://localhost:8080`，与本地 Kratos Server 通信

### Vite 代理配置

```typescript
// vite.config.ts
export default defineConfig({
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
});
```

### 错误处理策略

| 错误类型 | 处理方式 |
|----------|----------|
| 网络错误 | 全局 toast 提示，支持重试 |
| API 业务错误 | 根据 code 字段展示对应中文提示 |
| 工作流超时 | 前端轮询 + WebSocket 超时检测 |
| 图片加载失败 | 降级显示占位图，支持重新加载 |
