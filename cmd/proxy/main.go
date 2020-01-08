package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go-service-example/pkg/proxy"

	"go-service-example/pkg/processor"

	kitlog "github.com/go-kit/kit/log"
	"github.com/oklog/oklog/pkg/group"
	"google.golang.org/grpc"
)

func main() {

	httpAddr := ":8080"
	grpcAddr := ":8081"
	logger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stdout))
	errLogger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))

	gRPCconn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	if err != nil {
		errLogger.Log("message", "could not set up gRPC connection to processor", "addr", grpcAddr, "error", err)
	}
	client := processor.NewGRPCClient(gRPCconn)

	var (
		service     = proxy.NewProxyServer(client)
		endpoints   = proxy.MakeEndpoints(service)
		httpHandler = proxy.NewHTTPHandler(endpoints)
	)

	var g group.Group
	httpListener, err := net.Listen("tcp", httpAddr)
	if err != nil {
		errLogger.Log("message", "could not set up HTTP listner", "error", err)
	}
	g.Add(func() error {
		return http.Serve(httpListener, httpHandler)
	}, func(error) {
		httpListener.Close()
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

	logger.Log("HTTP", "listening", "addr", httpAddr)
	g.Run()
}
