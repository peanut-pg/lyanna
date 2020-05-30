package loadbalance

import (
	"context"
	"math/rand"

	"github.com/peanut-pg/lyanna/registry"
)

type RandomeBalance struct {
}

func (r *RandomeBalance) Name() string {
	return "random"
}

func (r *RandomeBalance) Select(ctx context.Context, nodes []*registry.Node) (node *registry.Node, err error) {
	if len(nodes) == 0 {
		err = ErrNotHaveNodes
		return
	}

	var totalWeight int
	for _, val := range nodes {
		if val.Weight == 0 {
			val.Weight = DefaultNodeWeight
		}
		totalWeight += val.Weight
	}
	curWeight := rand.Intn(totalWeight)
	curIndex := -1
	for index, node := range nodes {
		curWeight -= node.Weight
		if curWeight < 0 {
			curIndex = index
			break
		}
	}
	if curIndex == -1 {
		err = ErrNotHaveNodes
		return
	}

	node = nodes[curIndex]
	return
}
