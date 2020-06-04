package main

var routerTemplate = `
package router
import (
	"context"
	{{if not .Prefix}}
	"{{.Package.Name}}"
	{{else}}
	"{{.Prefix}}/{{.Package.Name}}"
	{{end}}
	{{if not .Prefix}}
	"controller"
	{{ else }}
	"{{.Prefix}}/controller"
	{{end}}
)
type RouterServer struct{}

{{range .RPCs}}
func (r *RouterServer) {{.Name}}(ctx context.Context, r*{{$.Package.Name}}.{{.RequestType}}))(resp*{{$.Package.Name}}.{{.ReturnsType}}){
	ctrl := &{{.Name}}Controller{}
	err = ctrl.CheckParams(ctx, r)
	if err != nil {
		return
	}
	resp, err = ctrl.Run(ctx, r)
	return
}
{{end}}
`
