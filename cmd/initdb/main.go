// cmd/initdb/main.go
// 一次性数据库初始化工具，用于在不启动完整服务的情况下创建数据库文件和表结构。
// 用法：go run ./cmd/initdb/
package main

import (
	"log"

	"github.com/user/manga-translator/internal/conf"
	"github.com/user/manga-translator/internal/data"
)

func main() {
	cfg, err := conf.Load("")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	_, cleanup, err := data.NewData(&cfg.Database)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	defer cleanup()

	log.Printf("数据库初始化成功：%s", cfg.Database.DSN)
}
