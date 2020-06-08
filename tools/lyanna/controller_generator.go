package main

import (
	"fmt"
	"html/template"
	"os"
	"path"

	"github.com/emicklei/proto"
)

func init() {
	ctrl := &CtrlGenerator{}
	Register("controller generator", ctrl)
}

type RpcMeta struct {
	RPC     *proto.RPC
	Package *proto.Package
	Prefix  string
}

type CtrlGenerator struct {
}

func (c *CtrlGenerator) Run(opt *Option, metaData *ServiceMetaData) (err error) {
	return c.generateGRPC(opt, metaData)
}

func (c *CtrlGenerator) generateGRPC(opt *Option, metaData *ServiceMetaData) (err error) {
	for _, rpc := range metaData.RPCs {

		filename := path.Join(opt.Output, "controller", fmt.Sprintf("%s.go", rpc.Name))
		fmt.Printf("------output is %v\n", opt.Output)
		fmt.Println(filename)
		var file *os.File
		file, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			fmt.Printf("open file:%s failed, err:%v\n", filename, err)
			return
		}
		defer file.Close()
		rpcMeta := &RpcMeta{}
		rpcMeta.Package = metaData.Package
		rpcMeta.RPC = rpc
		rpcMeta.Prefix = metaData.Prefix
		err = c.render(file, controllerTemplate, rpcMeta)
		if err != nil {
			fmt.Printf("render controller failed err:%v\n", err)
			return
		}
	}
	return
}

func (c *CtrlGenerator) render(file *os.File, data string, metaData *RpcMeta) (err error) {
	t := template.New("main")
	tml, err := t.Parse(data)
	if err != nil {
		fmt.Printf("render failed, err :%v\n", err)
		return
	}
	err = tml.Execute(file, metaData)
	return
}
