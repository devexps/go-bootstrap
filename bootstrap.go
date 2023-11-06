package bootstrap

import (
	"fmt"

	"github.com/devexps/go-micro/v2/log"
	"github.com/devexps/go-micro/v2/registry"

	conf "github.com/devexps/go-bootstrap/gen/api/go/conf/v1"
)

func Bootstrap(serviceInfo *ServiceInfo) (*conf.Bootstrap, log.Logger, registry.Registrar) {
	// inject command flags
	Flags := NewCommandFlags()
	Flags.Init()

	var err error

	// load configs
	if err = LoadBootstrapConfig(Flags.Conf); err != nil {
		panic(fmt.Sprintf("load config failed: %v", err))
	}
	// init logger
	ll := NewLoggerProvider(commonConfig.Logger, serviceInfo)

	// init registrar
	reg := NewRegistry(commonConfig.Registry)

	// init tracer
	if err = NewTracerProvider(commonConfig.Trace, serviceInfo); err != nil {
		panic(fmt.Sprintf("init tracer failed: %v", err))
	}
	return commonConfig, ll, reg
}
