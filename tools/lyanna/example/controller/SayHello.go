
package controller

import (
	"context"
	
	""/users/zhaofan/code/go_project/lyanna/tools/lyanna/example/generate/hello"
	
)
type SayHelloController struct{}

func (s *SayHelloController) CheckParams(ctx context.Context, r*hello.HelloRequest) (err error) {
	return
}

func (s *SayHelloController) SayHello(ctx context.Context, r*hello.HelloRequest) (resp*hello.HelloResponse, err error) {
	return
}


