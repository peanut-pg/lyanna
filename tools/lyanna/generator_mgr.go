package main

import "fmt"

var genMgr *GeneratorMgr

type GeneratorMgr struct {
	genMap map[string]Generator
}

func init() {
	genMgr = &GeneratorMgr{
		genMap: make(map[string]Generator),
	}
}

func Register(name string, gen Generator) (err error) {
	_, ok := genMgr.genMap[name]
	if ok {
		err = fmt.Errorf("generator %s is exists", name)
		return
	}
	genMgr.genMap[name] = gen
	return
}
