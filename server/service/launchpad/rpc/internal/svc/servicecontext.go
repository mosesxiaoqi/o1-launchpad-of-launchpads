package svc

import (
	"context"
	"database/sql"
	"time"

	"o1-launchpad/common/chain"
	"o1-launchpad/common/database"
	"o1-launchpad/common/model"
	"o1-launchpad/service/launchpad/rpc/internal/config"
)

type ServiceContext struct {
	Config     config.Config
	DB         *sql.DB
	Chain      *chain.Client
	Model      *model.Model
	Deployment chain.Deployment
}

func NewServiceContext(c config.Config) *ServiceContext {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := database.Open(ctx, c.Database)
	if err != nil {
		panic(err)
	}
	deployment, err := chain.LoadDeployment(c.DeploymentFile, c.Chain.ChainId)
	if err != nil {
		db.Close()
		panic(err)
	}
	return &ServiceContext{
		Config:     c,
		DB:         db,
		Chain:      chain.NewClient(c.Chain, deployment),
		Model:      model.New(db),
		Deployment: deployment,
	}
}
