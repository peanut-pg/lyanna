package main

import (
	"fmt"
	"os"
	"os/exec"
)

func init() {
	grpc := &GrpcGenerator{}
	Register("grpc generator", grpc)
}

type GrpcGenerator struct {
}

func (g *GrpcGenerator) Run(opt *Option) (err error) {
	outputParams := fmt.Sprintf("plugins=grpc:%s/generate/", opt.Output)
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
