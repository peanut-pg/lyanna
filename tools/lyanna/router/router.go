
package router
import (
	"context"
	
	"/users/zhaofan/code/go_project/lyanna/tools/lyanna/generate/hello"
	
	
	"/users/zhaofan/code/go_project/lyanna/tools/lyanna/controller"
	
)
type RouterServer struct{}


func (r *RouterServer) SayHello(ctx context.Context, r*hello.HelloRequest))(resp*hello.HelloResponse,err error){
	ctrl := &controller.SayHelloController{}
	err = ctrl.CheckParams(ctx, r)
	if err != nil {
		return
	}
	resp, err = ctrl.Run(ctx, r)
	return
}

