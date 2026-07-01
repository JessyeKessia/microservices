package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/jessyekessia/microservices-proto/golang/shipping"
	"github.com/jessyekessia/microservices/shipping/config"
	"github.com/jessyekessia/microservices/shipping/internal/ports"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Adapter struct {
	api  ports.APIPort
	port int
	shipping.UnimplementedShippingServer
}

func NewAdapter(api ports.APIPort, port int) *Adapter {
	return &Adapter{api: api, port: port}
}

func (a Adapter) Run() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		log.Fatalf("failed to listen on port %d, error: %v", a.port, err)
	}

	grpcServer := grpc.NewServer()
	shipping.RegisterShippingServer(grpcServer, a)

	if config.GetEnv() == "development" {
		reflection.Register(grpcServer)
	}

	log.Printf("gRPC shipping server running on port %d", a.port)

	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve grpc: %v", err)
	}
}
