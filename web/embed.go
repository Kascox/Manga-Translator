package web

import "embed"

// Dist 包含前端构建产物（web/dist/）。
//
//go:embed dist
var Dist embed.FS
