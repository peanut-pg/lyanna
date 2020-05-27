package etcd

import (
	"context"
	"testing"
	"time"

	registry "github.com/peanut-pg/lyanna"
)

func TestRegistry(t *testing.T) {
	registryInst, err := registry.InitRegistry(
		context.TODO(),
		"etcd",
		registry.WithAddrs([]string{"10.211.55.10:2379"}),
		registry.WithTimeout(time.Second),
		registry.WithRegistryPath("syncd/lyanna/"),
		registry.WithHeartBeat(5),
	)
	if err != nil {
		t.Fatalf("init registry failed, err :%v", err)
	}

}
