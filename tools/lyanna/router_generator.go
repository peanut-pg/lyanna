package main

import (
	"fmt"
	"html/template"
	"os"
	"path"
)

func init() {
	r := &RouterGenerator{}
	Register("router generator", r)
}

type RouterGenerator struct {
}

func (r *RouterGenerator) Run(opt *Option, metaData *ServiceMetaData) (err error) {
	filename := path.Join("./", opt.Output, "router/main.go")
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Printf("open file:%s failed, err:%v\n", filename, err)
		return
	}
	defer file.Close()
	err = r.render(file, routerTemplate, metaData)
	if err != nil {
		fmt.Printf("render failed, err:%v\n", err)
		return
	}
	return
}

func (r *RouterGenerator) render(file *os.File, data string, metaData *ServiceMetaData) (err error) {
	t := template.New("main")
	tml, err := t.Parse(data)
	if err != nil {
		fmt.Printf("render failed, err :%v\n", err)
		return
	}
	err = tml.Execute(file, metaData)
	return
}
