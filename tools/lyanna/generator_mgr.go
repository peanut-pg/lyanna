package main

import "fmt"

var genMgr *GeneratorMgr = &GeneratorMgr{
	genMap: make(map[string]Generator),
}

type GeneratorMgr struct {
	genMap map[string]Generator
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
	for _, gen := range g.genMap {
		err = gen.Run(opt)
		if err != nil {
			return
		}
	}
	return
}
