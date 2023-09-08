package bootstrap

import (
	"github.com/devexps/go-bootstrap/api/gen/go/common/conf"

	etcdGoMicro "github.com/devexps/go-micro/registry/etcd/v2"
	etcdClient "go.etcd.io/etcd/client/v3"

	consulGoMicro "github.com/devexps/go-micro/registry/consul/v2"
	consulClient "github.com/hashicorp/consul/api"

	zookeeperGoMicro "github.com/devexps/go-micro/registry/zookeeper/v2"
	"github.com/go-zookeeper/zk"

	"github.com/devexps/go-micro/v2/log"
	"github.com/devexps/go-micro/v2/registry"
)

type RegistryType string

const (
	RegistryTypeConsul  RegistryType = "consul"
	LoggerTypeEtcd      RegistryType = "etcd"
	LoggerTypeZooKeeper RegistryType = "zookeeper"
)

// NewRegistry creates a registry client
func NewRegistry(cfg *conf.Registry) registry.Registrar {
	if cfg == nil {
		return nil
	}

	switch RegistryType(cfg.Type) {
	case RegistryTypeConsul:
		return NewConsulRegistry(cfg)
	case LoggerTypeEtcd:
		return NewEtcdRegistry(cfg)
	case LoggerTypeZooKeeper:
		return NewZooKeeperRegistry(cfg)
	}

	return nil
}

// NewDiscovery creates a discovery client
func NewDiscovery(cfg *conf.Registry) registry.Discovery {
	if cfg == nil {
		return nil
	}

	switch RegistryType(cfg.Type) {
	case RegistryTypeConsul:
		return NewConsulRegistry(cfg)
	case LoggerTypeEtcd:
		return NewEtcdRegistry(cfg)
	case LoggerTypeZooKeeper:
		return NewZooKeeperRegistry(cfg)
	}

	return nil
}

// NewConsulRegistry creates a new registry client - Consul
func NewConsulRegistry(c *conf.Registry) *consulGoMicro.Registry {
	cfg := consulClient.DefaultConfig()
	cfg.Address = c.Consul.Address
	cfg.Scheme = c.Consul.Scheme

	var cli *consulClient.Client
	var err error
	if cli, err = consulClient.NewClient(cfg); err != nil {
		log.Fatal(err)
	}

	reg := consulGoMicro.New(cli, consulGoMicro.WithHealthCheck(c.Consul.HealthCheck))

	return reg
}

// NewEtcdRegistry creates a new registry client - Etcd
func NewEtcdRegistry(c *conf.Registry) *etcdGoMicro.Registry {
	cfg := etcdClient.Config{
		Endpoints: c.Etcd.Endpoints,
	}

	var err error
	var cli *etcdClient.Client
	if cli, err = etcdClient.New(cfg); err != nil {
		log.Fatal(err)
	}

	reg := etcdGoMicro.New(cli)

	return reg
}

// NewZooKeeperRegistry creates a new registry client - ZooKeeper
func NewZooKeeperRegistry(c *conf.Registry) *zookeeperGoMicro.Registry {
	conn, _, err := zk.Connect(c.Zookeeper.Endpoints, c.Zookeeper.Timeout.AsDuration())
	if err != nil {
		log.Fatal(err)
	}

	reg := zookeeperGoMicro.New(conn)

	return reg
}
