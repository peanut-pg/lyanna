package etcd

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/peanut-pg/lyanna/registry"
)

func TestRegistry(t *testing.T) {
	registryInst, err := registry.InitRegistry(
		context.TODO(),
		"etcd",
		registry.WithAddrs([]string{"192.168.31.182:2379"}),
		registry.WithTimeout(time.Second),
		registry.WithRegistryPath("syncd/lyanna/"),
		registry.WithHeartBeat(5),
	)
	if err != nil {
		t.Fatalf("init registry failed, err :%v", err)
	}
	service := &registry.Service{
		Name: "agent_service",
	}
	service.Nodes = append(service.Nodes,
		&registry.Node{
			IP:   "192.168.1.100",
			Port: 8001,
		},
		&registry.Node{
			IP:   "192.168.1.101",
			Port: 8001,
		},
	)
	err = registryInst.Register(context.TODO(), service)
	if err != nil {
		log.Fatal(err)
	}
	for {
		time.Sleep(time.Second)
	}

}
