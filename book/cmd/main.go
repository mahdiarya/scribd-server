package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	book "scribd/book/internal/features"
	eventstoregateway "scribd/book/internal/gateway/eventstore/grpc"
	grpcHandler "scribd/book/internal/handler/grpc"
	httpHandler "scribd/book/internal/handler/http"
	"scribd/book/internal/repository/memory"
	gen "scribd/gen/proto/v1"
	"scribd/pkg/discovery"
	"scribd/pkg/discovery/consul"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const serviceName = "book"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", 8484)); err != nil {
		panic(err)
	}
	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Fatalf("Failed to report healthy state: %v", err)
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)

	repo := memory.New()
	eventGateway := eventstoregateway.New(registry)
	ctrl := book.New(repo, eventGateway)

	//startHttp(ctrl)
	startGrpc(ctrl)
}

func startGrpc(ctrl *book.Application) {
	h := grpcHandler.New(ctrl)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8181))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	reflection.Register(srv)
	gen.RegisterBookServiceServer(srv, h)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Failed to start the gRPC server: %v", err)
	}
}

func startHttp(ctrl *book.Application) {
	h := httpHandler.New(ctrl)

	router := gin.Default()
	h.HandleRoutes(router)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}
