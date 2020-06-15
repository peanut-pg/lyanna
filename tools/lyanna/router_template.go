package main

var routerTemplate = `
package router
import (
	"context"
	{{if not .Prefix}}
	"generate/{{.Package.Name}}"
	{{else}}
	"{{.Prefix}}/generate/{{.Package.Name}}"
	{{end}}
	{{if not .Prefix}}
	"controller"
	{{ else }}
	"{{.Prefix}}/controller"
	{{end}}
)
type RouterServer struct{}

{{range .RPCs}}
func (r *RouterServer) {{.Name}}(ctx context.Context, r*{{$.Package.Name}}.{{.RequestType}}))(resp*{{$.Package.Name}}.{{.ReturnsType}},err error){
	ctrl := &controller.{{.Name}}Controller{}
	err = ctrl.CheckParams(ctx, r)
	if err != nil {
		return
	}
	resp, err = ctrl.Run(ctx, r)
	return
}
{{end}}
`
