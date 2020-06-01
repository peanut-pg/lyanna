package main

import (
	"fmt"
	"os"
	"path"

	"github.com/emicklei/proto"
)

func init() {
	ctrl := &CtrlGenerator{}
	Register("controller generator", ctrl)
}

type CtrlGenerator struct {
	service  *proto.Service
	messages []*proto.Message
	rpcs     []*proto.RPC
}

func (c *CtrlGenerator) Run(opt *Option) (err error) {
	render, err := os.Open(opt.Proto3FileName)
	if err != nil {
		fmt.Printf("open file:%s failed,err:%v\n", opt.Proto3FileName)
		return
	}
	defer render.Close()
	parser := proto.NewParser(render)
	definition, err := parser.Parse()
	if err != nil {
		fmt.Printf("parse file:%s failed, err:%v\n", opt.Proto3FileName, err)
		return
	}

	proto.Walk(
		definition,
		proto.WithService(c.handleService),
		proto.WithMessage(c.handleMessage),
		proto.WithRPC(c.handleRPC),
	)
	//fmt.Printf("parse protoc success, rpc:%#v\n", c.rpc)
	return c.generateGRPC(opt)
}

func (c *CtrlGenerator) handleService(s *proto.Service) {
	//fmt.Printf(s.Name)
	c.service = s
}

func (c *CtrlGenerator) handleMessage(m *proto.Message) {
	//fmt.Println(m.Name)
	c.messages = append(c.messages, m)
}

func (c *CtrlGenerator) handleRPC(r *proto.RPC) {
	/*
		fmt.Println(r.Name)
		fmt.Println(r.RequestType)
		fmt.Println(r.ReturnsType)
		fmt.Printf("rpc:%#v, comment:%#v\n", r, r.Comment)
	*/
	c.rpcs = append(c.rpcs, r)
}

func (c *CtrlGenerator) generateGRPC(opt *Option) (err error) {
	filename := path.Join(opt.Output, "controller", fmt.Sprintf("%s.go", c.service.Name))
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Printf("open file:%s failed, err:%v\n", filename, err)
		return
	}
	defer file.Close()
	fmt.Fprintf(file, "package controller\n")
	fmt.Fprintf(file, "import(\n")
	fmt.Fprintf(file, `"context"`)
	fmt.Fprintln(file)
	fmt.Fprintf(file, `hello "github.com/peanut-pg/lyanna/tools/lyanna/output/generate"`)
	fmt.Fprintln(file)
	fmt.Fprintf(file, ")\n")
	fmt.Fprintf(file, "type server struct{}\n")
	fmt.Fprint(file, "\n\n")

	for _, rpc := range c.rpcs {
		fmt.Fprintf(file, "func(s *server) %s(ctx context.Context, r*hello.%s)(resp*hello.%s, err error){\nreturn\n}\n\n",
			rpc.Name, rpc.RequestType, rpc.ReturnsType)
	}
	return
}
