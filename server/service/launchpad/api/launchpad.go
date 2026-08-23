package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"o1-launchpad/service/launchpad/api/internal/config"
	"o1-launchpad/service/launchpad/api/internal/handler"
	"o1-launchpad/service/launchpad/api/internal/problem"
	"o1-launchpad/service/launchpad/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/launchpad-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	if err := c.Validate(); err != nil {
		log.Fatal(err)
	}

	server := rest.MustNewServer(c.RestConf, rest.WithCors(c.Auth.AllowedOrigin))
	defer server.Stop()
	httpx.SetErrorHandlerCtx(func(_ context.Context, err error) (int, any) {
		code, body := problem.From(err)
		return code, body
	})

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
