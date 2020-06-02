package main

import (
	"fmt"
	"os"
	"path"

	"github.com/emicklei/proto"
)

var genMgr *GeneratorMgr = &GeneratorMgr{
	genMap:   make(map[string]Generator),
	metaData: &ServiceMetaData{},
}

var AllDirList []string = []string{
	"controller",
	"idl",
	"main",
	"scripts",
	"conf",
	"app/router",
	"app/config",
	"model",
	"generate",
	"router",
}

type GeneratorMgr struct {
	genMap   map[string]Generator
	metaData *ServiceMetaData
}

func Register(name string, gen Generator) (err error) {
	_, ok := genMgr.genMap[name]
	if ok {
		err = fmt.Errorf("generator %s is exists\n", name)
		return
	}
	genMgr.genMap[name] = gen
	return
}

func (g *GeneratorMgr) Run(opt *Option) (err error) {

	err = g.parseService(opt)
	if err != nil {
		return
	}

	err = g.createAllDir(opt)
	if err != nil {
		return
	}

	g.metaData.Prefix = opt.Prefix
	for _, gen := range g.genMap {
		err = gen.Run(opt, g.metaData)
		if err != nil {
			return
		}
	}
	return
}

func (g *GeneratorMgr) parseService(opt *Option) (err error) {
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
		proto.WithService(g.handleService),
		proto.WithMessage(g.handleMessage),
		proto.WithRPC(g.handleRPC),
	)
	return
}

func (g *GeneratorMgr) handleService(s *proto.Service) {
	g.metaData.Service = s
}

func (g *GeneratorMgr) handleMessage(m *proto.Message) {
	g.metaData.Messages = append(g.metaData.Messages, m)
}

func (g *GeneratorMgr) handleRPC(r *proto.RPC) {
	g.metaData.RPCs = append(g.metaData.RPCs, r)
}

func (g *GeneratorMgr) createAllDir(opt *Option) (err error) {
	for _, dir := range AllDirList {
		fullDir := path.Join(opt.Output, dir)
		err = os.MkdirAll(fullDir, 0755)
		if err != nil {
			fmt.Printf("mkdir dir %s failed, err:%v\n", fullDir, err)
			return
		}
	}
	return
}
