
package main

import (
	"net"
	"log"
	"google.golang.org/grpc"
	
	
	"/users/zhaofan/code/go_project/lyanna/tools/lyanna/router"
	

	
	"/users/zhaofan/code/go_project/lyanna/tools/lyanna/generate/hello"
	
)

var server = &router.RouterServer{}

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
