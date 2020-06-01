package main

import (
	"fmt"
	"os"
	"path"
)

func init() {
	main := &MainGenerator{}
	Register("controller generator", main)
}

type MainGenerator struct {
}

func (m *MainGenerator) Run(opt *Option) (err error) {
	filename := path.Join("./", opt.Output, "main/main.go")
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Printf("open file:%s failed, err:%v\n", filename, err)
		return
	}
	fmt.Fprintf(file, "package controller\n")
	fmt.Fprintf(file, "import(\n")
	fmt.Fprintf(file, `"net"`)
	fmt.Fprintln(file)
	fmt.Fprintf(file, `"google.golang.org/grpc"`)
	fmt.Fprintln(file)
	fmt.Fprintf(file, `"github.com/peanut-pg/lyanna/tools/lyanna/output/controller"`)
	fmt.Fprintln(file)

	fmt.Fprintf(file, ")\n")
	fmt.Fprintf(file, "var server = &controller.Server{}\n")
	fmt.Fprint(file, "\n\n")

	fmt.Fprintf(file, `
func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("failed to listen:%v",err)
	}
	s := grpc.NewServer()
	pb.RegisterHelloServiceServer(s, server)
	s.Serve(lis)
}
	`)
	defer file.Close()
	return
}
