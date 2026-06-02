package main

import (
	"log"

	"github.com/JessyeKessia/microservices/order/config"

	"github.com/JessyeKessia/microservices/order/internal/adapters/db"

	grpcAdapter "github.com/JessyeKessia/microservices/order/internal/adapters/grpc"

	"github.com/JessyeKessia/microservices/order/internal/application/core/api"
)

func main() {

	dbAdapter, err := db.NewAdapter(
		config.GetDataSourceURL(),
	)

	if err != nil {
		log.Fatalf(
			"Failed to connect to database. Error: %v",
			err,
		)
	}

	application := api.NewApplication(
		dbAdapter,
	)

	grpcAdapter := grpcAdapter.NewAdapter(
		application,
		config.GetApplicationPort(),
	)

	grpcAdapter.Run()
}
