# 工作项D2: OCR节点

## 基本信息

| 属性 | 值 |
|------|-----|
| 编号 | D2 |
| 阶段 | Phase D - 工作流 |
| 预估工期 | 1天 |
| 前置依赖 | C5 |
| 状态 | 未开始 |

## 工作描述

实现OCR识别节点(`internal/workflow/nodes/ocr.go`)：调用DeepSeek Vision API识别图片文字、解析返回的文字区域坐标、生成TextRegion列表。

具体工作：
1. 实现OCR节点函数，符合langgraphgo的Node接口签名
2. 从状态中获取当前图片路径，读取图片并转为Base64
3. 构建Vision API请求，包含图片Base64和OCR专用Prompt（要求返回文字内容及坐标）
4. 解析API响应，提取文字区域列表（文字内容+坐标位置）
5. 生成TextRegion列表写入工作流状态
6. 处理多语言识别场景（日/英/韩等）
7. 实现错误处理与重试机制（失败重试3次，间隔递增1s/3s/5s）

## 相关文档

| 文档 | 说明 |
|------|------|
| [04-agent-workflow.md](../docs/04-agent-workflow.md) | OCR识别节点章节、LLM适配层设计 |

## 技术要点

- DeepSeek Vision API调用：使用langchaingo的Vision模型接口，传入图片Base64
- OCR Prompt设计：需明确要求模型返回结构化结果（文字内容+坐标位置），可参考JSON格式输出
- TextRegion结构：包含文字内容(`text`)、区域坐标(`region_coords`，矩形四角坐标)、置信度等字段
- 多语言识别：Prompt中指定源语言或让模型自动检测
- 错误重试：使用指数退避策略，API调用失败时重试，返回格式异常时跳过该区域并记录日志
- 节点函数签名需匹配langgraphgo的`NodeFunc`类型：`func(ctx context.Context, state *TranslationState) (*TranslationState, error)`

## CheckList

- [ ] 实现OCR节点函数签名(符合langgraphgo Node接口)
- [ ] 构建Vision API请求(图片base64+prompt)
- [ ] 解析API响应(文字内容+坐标)
- [ ] 生成TextRegion列表写入状态
- [ ] 处理多语言识别(日/英/韩)
- [ ] 错误处理与重试
- [ ] 编写ocr_test.go(Mock Vision API)
