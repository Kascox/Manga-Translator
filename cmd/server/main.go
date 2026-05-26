package main

import (
	"log"

	"github.com/go-kratos/kratos/v2"
	kratoslog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"

	"github.com/user/manga-translator/internal/conf"
	"github.com/user/manga-translator/internal/data"
)

func main() {
	// 1. 加载配置
	cfg, err := conf.Load("")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 2. 初始化数据库（自动建表 + 迁移）
	_, dbCleanup, err := data.NewData(&cfg.Database)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	defer dbCleanup()

	// 3. 启动 HTTP Server
	logger := kratoslog.DefaultLogger
	httpSrv := http.NewServer(http.Address(cfg.Server.HTTP.Addr))
	app := kratos.New(
		kratos.Name("manga-translator"),
		kratos.Server(httpSrv),
		kratos.Logger(logger),
	)
	if err := app.Run(); err != nil {
		log.Fatal(err) //nolint:gocritic // log.Fatal is intentional; defer is a best-effort cleanup
	}
}
