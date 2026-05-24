# 工作项E3: 图片Canvas

## 基本信息

| 属性 | 值 |
|------|-----|
| 编号 | E3 |
| 阶段 | Phase E - 前端 |
| 预估工期 | 2天 |
| 前置依赖 | E1 |
| 状态 | 未开始 |

## 工作描述

实现中间图片Canvas组件(ImageCanvas)：Fabric.js集成、图片加载/缩放、框选工具、套索工具、文本区域标注显示、工具栏。

具体工作：
1. 集成Fabric.js Canvas，在useEffect中初始化，组件卸载时dispose释放资源
2. 实现图片加载与自适应显示：加载原图到Canvas，自动适应窗口大小
3. 实现缩放：鼠标滚轮缩放+底部缩放滑块，支持适应窗口(Ctrl+0)
4. 实现前一张/后一张切换按钮
5. 实现框选工具(R快捷键)：矩形框选创建新文字区域，生成Fabric.Rect对象
6. 实现套索工具(L快捷键)：自由形状圈选创建新文字区域，生成Fabric.Path对象
7. 实现文本区域标注渲染：API返回的TextRegion渲染为半透明矩形+文字标签
8. 实现标注hover高亮：鼠标悬停时高亮显示对应区域
9. 实现坐标系转换：Canvas坐标↔图片坐标双向转换
10. 实现自动识别模式：显示API返回的所有区域标注
11. 底部状态栏：显示文件名/识别区域数/翻译进度

## 相关文档

| 文档 | 说明 |
|------|------|
| [02-frontend-design.md](../docs/02-frontend-design.md) | ImageCanvas组件设计、Fabric.js集成要点、工具栏设计 |
| [01-product-interaction.md](../docs/01-product-interaction.md) | 文本选择与识别模式、自动识别/圈选模式交互 |

## 技术要点

- Fabric.js初始化：`new fabric.Canvas('canvasId')`，监听window resize自适应Canvas尺寸
- 文字区域对象模型：每个区域对应Fabric.Rect（框选）或Fabric.Path（套索），携带metadata: `{ regionId, status, isIgnored }`
- 框选工具：监听Canvas的mouse:down/mouse:move/mouse:up事件，绘制临时矩形，释放时创建Fabric.Rect
- 套索工具：记录鼠标路径点，释放时创建Fabric.Path闭合多边形
- 坐标系转换：Canvas显示坐标 = 图片坐标 × zoom + panOffset，需在交互时正确换算
- 区域高亮：鼠标hover时修改对象opacity/strokeWidth，配合useCanvasStore的selectedRegionId联动翻译面板
- 缩放控制：使用Canvas.setZoom()和Canvas.viewportTransform实现，限制缩放范围(0.1~5.0)
- 自动识别模式：OCR完成后接收TextRegion列表，批量创建Fabric.Rect标注到Canvas上
- 工具栏：选择(V)/框选(R)/套索(L)/缩放+/缩放-/适应，使用useCanvasStore的tool状态管理

## CheckList

- [ ] 集成Fabric.js Canvas
- [ ] 实现图片加载与自适应显示
- [ ] 实现缩放(滚轮+滑块)
- [ ] 实现前一张/后一张切换
- [ ] 实现框选工具(矩形选区)
- [ ] 实现套索工具(自由多边形)
- [ ] 实现文本区域标注渲染(半透明矩形+文字)
- [ ] 实现标注hover高亮
- [ ] 实现坐标系转换(Canvas坐标↔图片坐标)
- [ ] 实现自动识别模式(显示API返回的所有区域)
- [ ] 底部状态栏(文件名/识别数/翻译进度)
