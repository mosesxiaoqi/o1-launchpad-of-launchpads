package main

import (
	"flag"
	"fmt"
	"log"

	"o1-launchpad/service/indexer/internal/config"
	"o1-launchpad/service/indexer/internal/logic"
	"o1-launchpad/service/indexer/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
)

var configFile = flag.String("f", "etc/indexer.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	if err := c.Validate(); err != nil {
		log.Fatal(err)
	}

	svcCtx := svc.NewServiceContext(c)
	group := service.NewServiceGroup()
	defer group.Stop()
	group.Add(logic.NewIndexerLogic(svcCtx))
	fmt.Println("Starting launchpad indexer...")
	group.Start()
}
