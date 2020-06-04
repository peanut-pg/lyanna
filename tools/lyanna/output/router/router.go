
package router
import (
	"context"
	
	"hello"
	
	
	"controller"
	
)
type RouterServer struct{}


func (s *RouterServer) SayHello(ctx context.Context, r*hello.HelloRequest))(resp*hello.HelloResponse){
	ctrl := &SayHelloController{}
	err = ctrl.CheckParams(ctx, r)
	if err != nil {
		return
	}
	resp, err = ctrl.Run(ctx, r)
	return
}

