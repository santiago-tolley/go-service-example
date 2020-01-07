package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go-service-example/pb"
	"go-service-example/pkg/processor"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/oklog/oklog/pkg/group"
	"google.golang.org/grpc"
)

func main() {
	grpcAddr := ":8081"

	var (
		service    = processor.NewProcessorServer()
		endpoints  = processor.MakeEndpoints(service)
		grpcServer = processor.NewGRPCServer(endpoints)
	)

	var g group.Group

	grpcListener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatal("could not set up gRPC listner: ", err)
	}

	g.Add(func() error {
		baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
		pb.RegisterProcessorServer(baseServer, grpcServer)
		return baseServer.Serve(grpcListener)
	}, func(error) {
		grpcListener.Close()
	})

	cancelInterrupt := make(chan struct{})
	g.Add(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			return fmt.Errorf("received signal %s", sig)
		case <-cancelInterrupt:
			return nil
		}
	}, func(error) {
		close(cancelInterrupt)
	})

	fmt.Println("gRPC: listening on port ", grpcAddr)
	g.Run()
}
