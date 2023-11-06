package bootstrap

import (
	"path/filepath"

	// etcd
	etcdGoMicro "github.com/devexps/go-micro/registry/etcd/v2"
	etcdClient "go.etcd.io/etcd/client/v3"

	// consul
	consulGoMicro "github.com/devexps/go-micro/registry/consul/v2"
	consulClient "github.com/hashicorp/consul/api"

	// zookeeper
	zookeeperGoMicro "github.com/devexps/go-micro/registry/zookeeper/v2"
	"github.com/go-zookeeper/zk"

	// kubernetes
	k8sRegistry "github.com/devexps/go-micro/registry/k8s/v2"
	k8s "k8s.io/client-go/kubernetes"
	k8sRest "k8s.io/client-go/rest"
	k8sTools "k8s.io/client-go/tools/clientcmd"
	k8sUtil "k8s.io/client-go/util/homedir"

	"github.com/devexps/go-micro/v2/log"
	"github.com/devexps/go-micro/v2/registry"

	conf "github.com/devexps/go-bootstrap/gen/api/go/conf/v1"
)

type RegistryType string

const (
	RegistryTypeConsul   RegistryType = "consul"
	LoggerTypeEtcd       RegistryType = "etcd"
	LoggerTypeZooKeeper  RegistryType = "zookeeper"
	LoggerTypeKubernetes RegistryType = "kubernetes"
)

// NewRegistry creates a registration client
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
	case LoggerTypeKubernetes:
		return NewKubernetesRegistry(cfg)
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
	case LoggerTypeKubernetes:
		return NewKubernetesRegistry(cfg)
	}
	return nil
}

// NewConsulRegistry creates a registration discovery client - Consul
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

// NewEtcdRegistry creates a registration discovery client - Etcd
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

// NewZooKeeperRegistry creates a registration discovery client - ZooKeeper
func NewZooKeeperRegistry(c *conf.Registry) *zookeeperGoMicro.Registry {
	conn, _, err := zk.Connect(c.Zookeeper.Endpoints, c.Zookeeper.Timeout.AsDuration())
	if err != nil {
		log.Fatal(err)
	}
	reg := zookeeperGoMicro.New(conn)
	if err != nil {
		log.Fatal(err)
	}
	return reg
}

// NewKubernetesRegistry creates a registration discovery client - Kubernetes
func NewKubernetesRegistry(_ *conf.Registry) *k8sRegistry.Registry {
	restConfig, err := k8sRest.InClusterConfig()
	if err != nil {
		home := k8sUtil.HomeDir()
		kubeConfig := filepath.Join(home, ".kube", "config")
		restConfig, err = k8sTools.BuildConfigFromFlags("", kubeConfig)
		if err != nil {
			log.Fatal(err)
			return nil
		}
	}
	clientSet, err := k8s.NewForConfig(restConfig)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	reg := k8sRegistry.NewRegistry(clientSet)

	return reg
}
