
package main

import (
	"net"
	"log"
	"google.golang.org/grpc"
	"github.com/peanut-pg/lyanna/tools/lyanna/output/controller"
	hello "github.com/peanut-pg/lyanna/tools/lyanna/output/generate"
)

var server = &controller.Server{}

var port=":12345"

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("failed to listen:%v",err)
	}
	s := grpc.NewServer()
	hello.RegisterHelloServiceServer(s, server)
	s.Serve(lis)
}
