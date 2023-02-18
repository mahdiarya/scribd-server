package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"

	grpcHandler "scribd/eventstore/internal/handler/grpc"
	"scribd/eventstore/internal/repository/memory"
	"scribd/eventstore/pkg/natsutil"
	gen "scribd/gen/proto/v1"
)

func getServer() *grpcHandler.Server {
	repository := memory.New()
	natsComponent := natsutil.NewNATSComponent("eventstore-service")
	natsComponent.ConnectToServer(nats.DefaultURL)
	server := &grpcHandler.Server{
		Repository: repository,
		Nats:       natsComponent,
	}
	return server
}
func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 7979))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	server := getServer()
	gen.RegisterEventStoreServer(grpcServer, server)
	log.Printf("server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
