package main

import (
	"log"
	"net"

	"github.com/peanut-pg/lyanna/tools/lyanna/output/controller"
	hello "github.com/peanut-pg/lyanna/tools/lyanna/output/generate"
	"google.golang.org/grpc"
)

var server = &controller.Server{}

var port = ":12345"

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("failed to listen:%!v(MISSING)", err)
	}
	s := grpc.NewServer()
	hello.RegisterHelloServiceServer(s, server)
	s.Serve(lis)
}
