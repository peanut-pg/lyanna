package main

var routerTemplate = `
package router
import (
	"context"
	{{if not .Prefix}}
		{{.Package.Name}}
	{{else}}
		{{.Prefix}}/{{.Package.Name}}
	{{end}}
)
type RouterServer struct{}

{{range .Meta.RPCs}}
func (s *RouterServer) {{.Name}}(ctx context.Context, r*{{$.Package.Name}}.{{.RequestType}}))(resp*{{.$Package.Name}}.{{.ReturnsType}}){
	inst := &SayHelloController{}
	err = inst.CheckParams(ctx, r)
	if err != nil {
		return
	}
	resp, err = inst.Run(ctx, r)
	return
}
{{end}}
`
