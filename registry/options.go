package registry

import "time"

type Options struct {
	Addrs   []string
	Timeout time.Duration
	// syncd/agent_service/192.168.1.100:8001
	// syncd/agent_service/192.168.1.101:8001
	RegistryPath string
	HeartBeat    int64
}

type Option func(opts *Options)

func WithTimeout(timeout time.Duration) Option {
	return func(opts *Options) {
		opts.Timeout = timeout
	}
}

func WithAddrs(addrs []string) Option {
	return func(opts *Options) {
		opts.Addrs = addrs
	}
}

func WithRegistryPath(path string) Option {
	return func(opts *Options) {
		opts.RegistryPath = path
	}
}

func WithHeartBeat(heartBeat int64) Option {
	return func(opts *Options) {
		opts.HeartBeat = heartBeat
	}
}
