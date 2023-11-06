package bootstrap

import (
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/grpc"

	"github.com/devexps/go-micro/v2/config"
	"github.com/devexps/go-micro/v2/log"

	// file
	fileGoMicro "github.com/devexps/go-micro/v2/config/file"

	// etcd
	etcdGoMicro "github.com/devexps/go-micro/config/etcd/v2"
	etcdClient "go.etcd.io/etcd/client/v3"

	// consul
	consulGoMicro "github.com/devexps/go-micro/config/consul/v2"
	consulApi "github.com/hashicorp/consul/api"

	// kubernetes
	k8sGoMicro "github.com/devexps/go-micro/config/k8s/v2"
	k8sUtil "k8s.io/client-go/util/homedir"

	conf "github.com/devexps/go-bootstrap/gen/api/go/conf/v1"
)

const (
	ConfigTypeLocalFile  ConfigType = "file"
	ConfigTypeConsul     ConfigType = "consul"
	ConfigTypeEtcd       ConfigType = "etcd"
	ConfigTypeKubernetes ConfigType = "kubernetes"

	remoteConfigSourceConfigFile = "remote.yaml"
)

var commonConfig = &conf.Bootstrap{}
var configList []interface{}

type ConfigType string

// RegisterConfig registration configuration
func RegisterConfig(c interface{}) {
	initBootstrapConfig()
	configList = append(configList, c)
}

func initBootstrapConfig() {
	if len(configList) > 0 {
		return
	}
	configList = append(configList, commonConfig)

	if commonConfig.Server == nil {
		commonConfig.Server = &conf.Server{}
		configList = append(configList, commonConfig.Server)
	}
	if commonConfig.Client == nil {
		commonConfig.Client = &conf.Client{}
		configList = append(configList, commonConfig.Client)
	}
	if commonConfig.Data == nil {
		commonConfig.Data = &conf.Data{}
		configList = append(configList, commonConfig.Data)
	}
	if commonConfig.Trace == nil {
		commonConfig.Trace = &conf.Tracer{}
		configList = append(configList, commonConfig.Trace)
	}
	if commonConfig.Logger == nil {
		commonConfig.Logger = &conf.Logger{}
		configList = append(configList, commonConfig.Logger)
	}
	if commonConfig.Registry == nil {
		commonConfig.Registry = &conf.Registry{}
		configList = append(configList, commonConfig.Registry)
	}
}

// NewConfigProvider creates a configuration
func NewConfigProvider(configPath string) config.Config {
	err, rc := LoadRemoteConfigSourceConfigs(configPath)
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

// LoadBootstrapConfig loader boot configuration
func LoadBootstrapConfig(configPath string) error {
	cfg := NewConfigProvider(configPath)

	var err error

	if err = cfg.Load(); err != nil {
		return err
	}
	initBootstrapConfig()

	if err = scanConfigs(cfg); err != nil {
		return err
	}
	return nil
}

func scanConfigs(cfg config.Config) error {
	initBootstrapConfig()

	for _, c := range configList {
		if err := cfg.Scan(c); err != nil {
			return err
		}
	}
	return nil
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

// LoadRemoteConfigSourceConfigs loads the local configuration of the remote configuration source
func LoadRemoteConfigSourceConfigs(configPath string) (error, *conf.RemoteConfig) {
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
		if err := cfg.Close(); err != nil {
			panic(err)
		}
	}(cfg)

	var err error

	if err = cfg.Load(); err != nil {
		return err, nil
	}
	if err = scanConfigs(cfg); err != nil {
		return err, nil
	}
	return nil, commonConfig.Config
}

// NewRemoteConfigSource 创建一个远程配置源
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
	case ConfigTypeKubernetes:
		return NewKubernetesConfigSource(c)
	}
}

// getConfigKey gets the legal configuration name
func getConfigKey(configKey string, useBackslash bool) string {
	if useBackslash {
		return strings.Replace(configKey, `.`, `/`, -1)
	} else {
		return configKey
	}
}

// NewFileConfigSource creates a local file configuration source
func NewFileConfigSource(filePath string) config.Source {
	return fileGoMicro.NewSource(filePath)
}

// NewEtcdConfigSource creates a remote configuration source - Etcd
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

// NewConsulConfigSource creates a remote configuration source - Consul
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

// NewKubernetesConfigSource creates a remote configuration source - Kubernetes
func NewKubernetesConfigSource(c *conf.RemoteConfig) config.Source {
	source := k8sGoMicro.NewSource(
		k8sGoMicro.Namespace(c.Kubernetes.Namespace),
		k8sGoMicro.LabelSelector(""),
		k8sGoMicro.KubeConfig(filepath.Join(k8sUtil.HomeDir(), ".kube", "config")),
	)
	return source
}
