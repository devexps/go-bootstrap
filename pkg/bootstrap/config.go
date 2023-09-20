package bootstrap

import (
	"os"
	"strings"

	"google.golang.org/grpc"

	"github.com/devexps/go-bootstrap/api/gen/go/common/conf"

	"github.com/devexps/go-micro/v2/config"
	"github.com/devexps/go-micro/v2/log"

	// file
	fileGoMicro "github.com/devexps/go-micro/v2/config/file"

	// consul
	consulGoMicro "github.com/devexps/go-micro/config/consul/v2"
	consulApi "github.com/hashicorp/consul/api"

	// etcd
	etcdGoMicro "github.com/devexps/go-micro/config/etcd/v2"
	etcdClient "go.etcd.io/etcd/client/v3"
)

type ConfigType string

const (
	ConfigTypeLocalFile ConfigType = "file"
	ConfigTypeConsul    ConfigType = "consul"
	ConfigTypeEtcd      ConfigType = "etcd"

	remoteConfigSourceConfigFile = "remote.yaml"
)

// LoadBootstrapConfig loader boot configuration
func LoadBootstrapConfig(configPath string) *conf.Bootstrap {
	cfg := NewConfigProvider(configPath)
	if err := cfg.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := cfg.Scan(&bc); err != nil {
		panic(err)
	}

	if bc.Server == nil {
		bc.Server = &conf.Server{}
		_ = cfg.Scan(&bc.Server)
	}

	if bc.Client == nil {
		bc.Client = &conf.Client{}
		_ = cfg.Scan(&bc.Client)
	}

	if bc.Data == nil {
		bc.Data = &conf.Data{}
		_ = cfg.Scan(&bc.Data)
	}

	if bc.Trace == nil {
		bc.Trace = &conf.Tracer{}
		_ = cfg.Scan(&bc.Trace)
	}

	if bc.Logger == nil {
		bc.Logger = &conf.Logger{}
		_ = cfg.Scan(&bc.Logger)
	}

	if bc.Registry == nil {
		bc.Registry = &conf.Registry{}
		_ = cfg.Scan(&bc.Registry)
	}

	return &bc
}

// NewConfigProvider creates a configuration
func NewConfigProvider(configPath string) config.Config {
	rc, err := LoadRemoteConfigSourceConfigs(configPath)
	if err != nil {
		log.Error("LoadRemoteConfigSourceConfigs: ", err.Error())
	}
	if rc != nil {
		return config.New(
			config.WithSource(
				NewFileConfigSource(configPath),
				NewRemoteConfigSource(rc),
			),
		)
	} else {
		return config.New(
			config.WithSource(
				NewFileConfigSource(configPath),
			),
		)
	}
}

// LoadRemoteConfigSourceConfigs loads the local configuration of the remote configuration source
func LoadRemoteConfigSourceConfigs(configPath string) (*conf.RemoteConfig, error) {
	configPath = configPath + "/" + remoteConfigSourceConfigFile
	if !pathExists(configPath) {
		return nil, nil
	}

	cfg := config.New(
		config.WithSource(
			NewFileConfigSource(configPath),
		),
	)
	defer func(cfg config.Config) {
		err := cfg.Close()
		if err != nil {
			panic(err)
		}
	}(cfg)

	var err error

	if err = cfg.Load(); err != nil {
		return nil, err
	}

	var rc conf.Bootstrap
	if err = cfg.Scan(&rc); err != nil {
		return nil, err
	}

	return rc.Config, nil
}

// NewFileConfigSource creates a local file configuration source
func NewFileConfigSource(filePath string) config.Source {
	return fileGoMicro.NewSource(filePath)
}

// NewRemoteConfigSource creates a remote configuration source
func NewRemoteConfigSource(c *conf.RemoteConfig) config.Source {
	switch ConfigType(c.Type) {
	default:
		fallthrough
	case ConfigTypeLocalFile:
		return nil
	case ConfigTypeConsul:
		return NewConsulConfigSource(c)
	case ConfigTypeEtcd:
		return NewEtcdConfigSource(c)
	}
}

// NewConsulConfigSource creates a remote config source - Consul
func NewConsulConfigSource(c *conf.RemoteConfig) config.Source {
	cfg := consulApi.DefaultConfig()
	cfg.Address = c.Consul.Address
	cfg.Scheme = c.Consul.Scheme

	cli, err := consulApi.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	source, err := consulGoMicro.New(cli,
		consulGoMicro.WithPath(getConfigKey(c.Consul.Key, true)),
	)
	if err != nil {
		log.Fatal(err)
	}

	return source
}

// NewEtcdConfigSource creates a remote config source - Etcd
func NewEtcdConfigSource(c *conf.RemoteConfig) config.Source {
	cfg := etcdClient.Config{
		Endpoints:   c.Etcd.Endpoints,
		DialTimeout: c.Etcd.Timeout.AsDuration(),
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	}

	cli, err := etcdClient.New(cfg)
	if err != nil {
		panic(err)
	}

	source, err := etcdGoMicro.New(cli, etcdGoMicro.WithPath(getConfigKey(c.Etcd.Key, true)))
	if err != nil {
		log.Fatal(err)
	}

	return source
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func getConfigKey(configKey string, useBackslash bool) string {
	if useBackslash {
		return strings.Replace(configKey, `.`, `/`, -1)
	} else {
		return configKey
	}
}
