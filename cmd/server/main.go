package main

import (
	"context"
	"flag"
	"fmt"

	"go-dianping/cmd/server/wire"
	"go-dianping/pkg/config"
	"go-dianping/pkg/log"
	"go.uber.org/zap"
)

// @title           Go Dian Ping
// @version         1.0.0
// @description     golang 实现的黑马点评
// @contact.name   llmons
// @contact.url    https://github.com/llmons
// @contact.email  llmons@foxmail.com
// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT
// @host      localhost:8080
// @securityDefinitions.apiKey Bearer
// @in header
// @name Authorization
// @externalDocs.description  GitHub
// @externalDocs.url          https://github.com/llmons/go-dianping
func main() {
	var envConf = flag.String("conf", "config/local.yml", "config path, eg: -conf ./config/local.yml")
	flag.Parse()
	conf := config.NewConfig(*envConf)

	logger := log.NewLog(conf)

	app, cleanup, err := wire.NewWire(conf, logger)
	defer cleanup()
	if err != nil {
		panic(err)
	}
	logger.Info("server start", zap.String("host", fmt.Sprintf("http://%s:%d", conf.GetString("http.host"), conf.GetInt("http.port"))))
	logger.Info("docs addr", zap.String("addr", fmt.Sprintf("http://%s:%d/swagger/index.html", conf.GetString("http.host"), conf.GetInt("http.port"))))
	if err = app.Run(context.Background()); err != nil {
		panic(err)
	}
}
