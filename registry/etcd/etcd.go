package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"time"


	"go.etcd.io/etcd/clientv3"
)

const (
	MaxServiceNum = 8
)

var (
	etcdRegistry *EtcdRegistry = &EtcdRegistry{
		serviceCh:          make(chan *registry.Service, MaxServiceNum),
		registryServiceMap: make(map[string]*RegisterService, MaxServiceNum),
	}
)

type RegisterService struct {
	id          clientv3.LeaseID
	service     *registry.Service
	registered  bool // 服务是否注册
	keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse
}

type EtcdRegistry struct {
	options   *registry.Options
	client    *clientv3.Client
	serviceCh chan *registry.Service

	registryServiceMap map[string]*RegisterService
}

func init() {
	err := registry.RegisterPlugin(etcdRegistry)
	if err != nil {
		fmt.Printf("register plugin [%v] errror [%v]\n", etcdRegistry.Name(), err)
	}
	// 后台将服务注册到etcd中
	go etcdRegistry.Run()
}

func (e EtcdRegistry) Name() string {
	return "etcd"
}

func (e EtcdRegistry) Init(ctx context.Context, opts ...registry.Option) (err error) {
	e.options = &registry.Options{}
	for _, opt := range opts {
		opt(e.options)
	}
	e.client, err = clientv3.New(clientv3.Config{
		Endpoints:   e.options.Addrs,
		DialTimeout: e.options.Timeout,
	})
	if err != nil {
		return err
	}
	return err
}

func (e EtcdRegistry) Register(ctx context.Context, service *registry.Service) (err error) {
	select {
	case e.serviceCh <- service:
	default:
		err = fmt.Errorf("register chan is full")
		return err
	}
	return err
}

func (e EtcdRegistry) Unregister(ctx context.Context, service *registry.Service) (err error) {
	panic("implement me")
}

// 后台进行服务的注册
func (e EtcdRegistry) Run() {
	for {
		select {
		case service := <-e.serviceCh:
			_, ok := e.registryServiceMap[service.Name]
			if ok {
				// 已经注册过的服务不再重复注册
				break
			}
			registryService := &RegisterService{
				service: service,
			}
			e.registryServiceMap[service.Name] = registryService
		default:
			e.registerOrKeepAlive()
			time.Sleep(time.Millisecond * 500)
		}
	}
}

// 注册或续约
func (e *EtcdRegistry) registerOrKeepAlive() {
	for _, registryService := range e.registryServiceMap {
		if registryService.registered {
			// 已经注册的进行续约
			e.keepAlive(registryService)
			continue
		}
		// 进行注册
		e.registerService(registryService)
	}
}

// 进行续约
func (e *EtcdRegistry) keepAlive(registryService *RegisterService) {
	select {
	case resp := <-registryService.keepAliveCh:
		if resp == nil {
			registryService.registered = false
			return
		}
		fmt.Printf("service:%s node:%s port:%v\n", registryService.service.Name,
			registryService.service.Nodes[0].IP, registryService.service.Nodes[0].Port)
	}
	return
}

// 进行服务注册
func (e *EtcdRegistry) registerService(registryService *RegisterService) (err error) {
	resp, err := e.client.Grant(context.TODO(), e.options.HeartBeat)
	if err != nil {
		return
	}
	registryService.id = resp.ID
	for _, node := range registryService.service.Nodes {
		tmp := &registry.Service{
			Name: registryService.service.Name,
			Nodes: []*registry.Node{
				node,
			},
		}
		data, err := json.Marshal(tmp)
		if err != nil {
			continue
		}
		key := e.serviceNodePath(tmp)
		_, err = e.client.Put(context.TODO(), key, string(data), clientv3.WithLease(resp.ID))
		if err != nil {
			continue
		}
		// 如果注册成功，则进行keepalive
		ch, kaerr := e.client.KeepAlive(context.TODO(), resp.ID)
		if kaerr != nil {
			continue
		}
		registryService.keepAliveCh = ch
		registryService.registered = true
	}
	return
}

func (e *EtcdRegistry) serviceNodePath(service *registry.Service) string {
	nodeIP := fmt.Sprintf("%s:%d", service.Nodes[0].IP, service.Nodes[0].Port)
	return path.Join(e.options.RegistryPath, nodeIP, service.Name)
}
