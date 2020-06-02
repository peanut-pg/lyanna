package main

import (
	"fmt"
	"html/template"
	"os"
	"path"
)

func init() {
	mainGenerator := &MainGenerator{}
	Register("main generator", mainGenerator)
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
	defer file.Close()
	err = m.render(file, mainTemplate)
	if err != nil {
		fmt.Printf("render failed, err:%v\n", err)
		return
	}
	return
}

func (m *MainGenerator) render(file *os.File, data string) (err error) {
	t := template.New("main")
	tml, err := t.Parse(data)
	if err != nil {
		fmt.Printf("render failed, err :%v\n", err)
		return
	}
	err = tml.Execute(file, nil)
	return
}
