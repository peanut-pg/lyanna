package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

func init() {
	grpc := &GrpcGenerator{}
	Register("grpc generator", grpc)
}

type GrpcGenerator struct {
}

func (g *GrpcGenerator) Run(opt *Option, metaData *ServiceMetaData) (err error) {
	dir := path.Join(opt.Output, "generate", metaData.Package.Name)
	os.MkdirAll(dir, 0755)
	outputParams := fmt.Sprintf("plugins=grpc:%s/generate/%s", opt.Output, metaData.Package.Name)

	cmd := exec.Command("protoc", "--go_out", outputParams, opt.Proto3FileName)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		fmt.Printf("grpc generator failed, err:%v\n", err)
		return
	}
	return
}
