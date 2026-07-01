package main

import (
	"log"

	"github.com/jessyekessia/microservices/shipping/config"
	"github.com/jessyekessia/microservices/shipping/internal/adapters/db"
	grpcAdapter "github.com/jessyekessia/microservices/shipping/internal/adapters/grpc"
	"github.com/jessyekessia/microservices/shipping/internal/application/core/api"
)

func main() {
	dbAdapter, err := db.NewAdapter(config.GetDataSourceURL())
	if err != nil {
		log.Fatalf("Failed to connect to database. Error: %v", err)
	}

	application := api.NewApplication(dbAdapter)
	grpcServer := grpcAdapter.NewAdapter(application, config.GetApplicationPort())
	grpcServer.Run()
}
