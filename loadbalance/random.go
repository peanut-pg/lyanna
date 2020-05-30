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
	index := rand.Intn(len(nodes))
	node = nodes[index]
	return
}
