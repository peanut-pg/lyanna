package main

var controllerTemplate = `
package controller

import (
	"context"
	{{if not .Prefix}}
	"generate/{{.Package.Name}}"
	{{else}}
	""{{.Prefix}}/generate/{{.Package.Name}}"
	{{end}}
)
type {{.RPC.Name}}Controller struct{}

func (s *{{.RPC.Name}}Controller) CheckParams(ctx context.Context, r*{{.Package.Name}}.{{.RPC.RequestType}}) (err error) {
	return
}

func (s *{{.RPC.Name}}Controller) SayHello(ctx context.Context, r*{{.Package.Name}}.{{.RPC.RequestType}}) (resp*{{.Package.Name}}.{{.RPC.ReturnsType}}, err error) {
	return
}


`
