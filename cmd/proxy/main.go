package main

// move this functionality to a makefile?
import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go-service-example/pkg/proxy"

	// "github.com/go-kit/kit/transport/http"
	"go-service-example/pkg/processor"

	"github.com/oklog/oklog/pkg/group"
	"google.golang.org/grpc"
)

func main() {

	httpAddr := ":8080"
	grpcAddr := ":8081"

	gRPCconn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("could not connect to processor: %v", err)
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
		os.Exit(1)
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

	fmt.Println("HTTP: listening on port ", httpAddr)
	g.Run()
}
