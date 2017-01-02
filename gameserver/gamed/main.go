// gamed is the main entrypoint for the rummy game server.
package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/glog"
	"google.golang.org/grpc"

	"rummy"
	"rummy/gameserver"
)

func main() {
	port := flag.Int("port", 8081, "Port to run gRPC service on")
	flag.Parse()

	endpoint := fmt.Sprintf(":%d", *port)
	lis, err := net.Listen("tcp", endpoint)
	if err != nil {
		glog.Fatalf("failed to listen: %v", err)
	}

	glog.Info("Initializing RPC server")
	grpcServer := grpc.NewServer()
	rummyServer := gameserver.NewRummyServer(endpoint)
	rummy.RegisterRummyServiceServer(grpcServer, rummyServer)
	go grpcServer.Serve(lis)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	glog.Info(<-ch)
	glog.Info("Shutting down")

	grpcServer.Stop()
}
