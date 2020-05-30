package loadbalance

import (
	"context"
	"fmt"
	"testing"

	"github.com/peanut-pg/lyanna/registry"
)

func TestRandomeBalance_Select(t *testing.T) {
	balance := &RandomeBalance{}
	var nodes []*registry.Node
	var weights = []int{50, 100, 150}
	for i := 0; i < 4; i++ {
		node := &registry.Node{
			IP:     fmt.Sprintf("192.168.1.%d", i),
			Port:   8001,
			Weight: weights[i%3],
		}
		fmt.Printf("node:%#v\n", node)
		nodes = append(nodes, node)
	}

	countStat := make(map[string]int)
	for i := 0; i < 1000; i++ {
		node, err := balance.Select(context.TODO(), nodes)
		if err != nil {
			t.Fatalf("select failed, err :%v", err)
		}
		countStat[node.IP]++
	}
	for key, val := range countStat {
		fmt.Printf("ip:%s, count:%d\n", key, val)
	}
}
