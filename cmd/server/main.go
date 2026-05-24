package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func main() {
	logger := log.DefaultLogger
	httpSrv := http.NewServer(http.Address(":8080"))
	app := kratos.New(
		kratos.Name("manga-translator"),
		kratos.Server(httpSrv),
		kratos.Logger(logger),
	)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
