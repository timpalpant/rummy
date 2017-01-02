// gamed is the main entrypoint for the rummy game server.
// It runs a gRPC server, as well as a JSON reverse-proxy for that
// server using gRPC-gateway.
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"rummy"
	"rummy/gameserver"
)

func runProxy(ctx context.Context, grpcEndpoint string, port int) error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := rummy.RegisterRummyServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
	if err != nil {
		return err
	}

	jsonEndpoint := fmt.Sprintf(":%d", port)
	http.ListenAndServe(jsonEndpoint, mux)
	return nil
}

func main() {
	port := flag.Int("port", 8081, "Port to run gRPC service on")
	proxyPort := flag.Int("proxyport", 8082, "Port to run JSON proxy on")
	seed := flag.Int64("seed", 1, "Seed for random shuffling")
	flag.Parse()

	rand.Seed(*seed)

	endpoint := fmt.Sprintf(":%d", *port)
	lis, err := net.Listen("tcp", endpoint)
	if err != nil {
		glog.Fatalf("failed to listen: %v", err)
	}

	glog.Info("Initializing RPC server")
	grpcServer := grpc.NewServer()
	rummyServer := gameserver.NewRummyServer()
	rummy.RegisterRummyServiceServer(grpcServer, rummyServer)
	go grpcServer.Serve(lis)

	glog.Info("Starting JSON proxy")
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go runProxy(ctx, endpoint, *proxyPort)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	glog.Info(<-ch)
	glog.Info("Shutting down")

	grpcServer.Stop()
}
