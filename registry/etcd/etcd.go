package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"path"
	"sync"
	"sync/atomic"
	"time"

	"github.com/peanut-pg/lyanna/registry"
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

type AllServiceInfo struct {
	serviceMap map[string]*registry.Service
}

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
	value              atomic.Value // 用于实现缓存服务的原子操作
	lock               sync.Mutex
}

func init() {
	allServiceInfo := &AllServiceInfo{
		serviceMap: make(map[string]*registry.Service, MaxServiceNum),
	}
	etcdRegistry.value.Store(allServiceInfo)
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

func (e *EtcdRegistry) Init(ctx context.Context, opts ...registry.Option) (err error) {
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

func (e *EtcdRegistry) Register(ctx context.Context, service *registry.Service) (err error) {
	select {
	case e.serviceCh <- service:
	default:
		err = fmt.Errorf("register chan is full")
		return
	}
	return
}

func (e *EtcdRegistry) Unregister(ctx context.Context, service *registry.Service) (err error) {
	return nil
}

func (e *EtcdRegistry) GetService(ctx context.Context, name string) (service *registry.Service, err error) {

	// 先从缓存中获取
	service, ok := e.getServiceFromCache(ctx, name)
	if ok {
		return
	}
	// 缓存没有 则从etcd中获取
	e.lock.Lock()
	defer e.lock.Unlock()
	service, ok = e.getServiceFromCache(ctx, name)
	if ok {
		return
	}
	// 这里才会真正从etcd中获取
	key := e.getServicePath(name)
	resp, err := e.client.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return
	}
	service = &registry.Service{
		Name: name,
	}

	for _, kv := range resp.Kvs {
		//fmt.Printf("index:%v key:%v val:%v\n", index, string(kv.Key), string(kv.Value))
		value := kv.Value
		var tmpService registry.Service
		err = json.Unmarshal(value, &tmpService)
		if err != nil {
			return
		}
		for _, node := range tmpService.Nodes {
			service.Nodes = append(service.Nodes, node)
		}

	}
	allServiceInfoOld := e.value.Load().(*AllServiceInfo)
	var allServiceInfoNew = &AllServiceInfo{
		serviceMap: make(map[string]*registry.Service, MaxServiceNum),
	}
	for key, val := range allServiceInfoOld.serviceMap {
		allServiceInfoNew.serviceMap[key] = val
	}
	allServiceInfoNew.serviceMap[name] = service
	e.value.Store(allServiceInfoNew)
	return
}

// 后台进行服务的注册
func (e *EtcdRegistry) Run() {
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
		/*
			fmt.Printf("service:%s node:%s port:%v\n", registryService.service.Name,
				registryService.service.Nodes[0].IP, registryService.service.Nodes[0].Port)
			fmt.Printf("service:%s node:%s port:%v\n", registryService.service.Name,
				registryService.service.Nodes[1].IP, registryService.service.Nodes[1].Port)
		*/
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
		log.Println("service key:", key)
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
	return path.Join(e.options.RegistryPath, service.Name, nodeIP)
}

// 获取服务的前缀 如：/syncd/lyanna/agent_service
func (e *EtcdRegistry) getServicePath(name string) string {
	return path.Join(e.options.RegistryPath, name)
}

func (e *EtcdRegistry) getServiceFromCache(ctx context.Context, name string) (service *registry.Service, ok bool) {
	allServiceInfo := e.value.Load().(*AllServiceInfo)
	service, ok = allServiceInfo.serviceMap[name]
	return
}
