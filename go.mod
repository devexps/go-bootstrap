module github.com/devexps/go-bootstrap

go 1.18

replace (
	google.golang.org/grpc => google.golang.org/grpc v1.46.2
	google.golang.org/protobuf => google.golang.org/protobuf v1.28.0
)

require (
	github.com/devexps/go-micro/config/consul/v2 v2.0.0-20231106083309-0e1ceee74915
	github.com/devexps/go-micro/config/etcd/v2 v2.0.0-20231106083309-0e1ceee74915
	github.com/devexps/go-micro/config/k8s/v2 v2.0.0-20231106083309-0e1ceee74915
	github.com/devexps/go-micro/log/fluent/v2 v2.0.0-20231106083309-0e1ceee74915
	github.com/devexps/go-micro/log/logrus/v2 v2.0.0-20231106083309-0e1ceee74915
	github.com/devexps/go-micro/log/zap/v2 v2.0.0-20231106083309-0e1ceee74915
	github.com/devexps/go-micro/registry/consul/v2 v2.0.0-20231106083309-0e1ceee74915
	github.com/devexps/go-micro/registry/etcd/v2 v2.0.0-20231106083309-0e1ceee74915
	github.com/devexps/go-micro/registry/k8s/v2 v2.0.0-20231106083309-0e1ceee74915
	github.com/devexps/go-micro/registry/zookeeper/v2 v2.0.0-20231106083309-0e1ceee74915
	github.com/devexps/go-micro/v2 v2.0.0
	github.com/devexps/go-pkg/v2 v2.0.0
	github.com/go-redis/redis/extra/redisotel/v8 v8.11.5
	github.com/go-redis/redis/v8 v8.11.5
	github.com/go-zookeeper/zk v1.0.3
	github.com/google/subcommands v1.2.0
	github.com/gorilla/handlers v1.5.2
	github.com/hashicorp/consul/api v1.20.0
	github.com/olekukonko/tablewriter v0.0.5
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.8.0
	go.etcd.io/etcd/client/v3 v3.5.8
	go.opentelemetry.io/otel v1.14.0
	go.opentelemetry.io/otel/exporters/jaeger v1.14.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.14.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.14.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.14.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.14.0
	go.opentelemetry.io/otel/exporters/zipkin v1.14.0
	go.opentelemetry.io/otel/sdk v1.14.0
	go.uber.org/zap v1.23.0
	golang.org/x/tools v0.4.0
	google.golang.org/grpc v1.53.0
	google.golang.org/protobuf v1.28.1
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
	k8s.io/client-go v0.26.3
)

require (
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/armon/go-metrics v0.0.0-20180917152333-f0300d1749da // indirect
	github.com/cenkalti/backoff/v4 v4.2.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/emicklei/go-restful/v3 v3.9.0 // indirect
	github.com/fatih/color v1.9.0 // indirect
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/fluent/fluent-logger-golang v1.9.0 // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.20.0 // indirect
	github.com/go-openapi/swag v0.19.14 // indirect
	github.com/go-playground/form/v4 v4.2.0 // indirect
	github.com/go-redis/redis/extra/rediscmd/v8 v8.11.5 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/gnostic v0.5.7-v3refs // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.7.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.1 // indirect
	github.com/hashicorp/go-hclog v0.12.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.0.0 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hashicorp/serf v0.10.1 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.6 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/openzipkin/zipkin-go v0.4.1 // indirect
	github.com/philhofer/fwd v1.1.2 // indirect
	github.com/shirou/gopsutil/v3 v3.21.8 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/tinylib/msgp v1.1.8 // indirect
	github.com/tklauser/go-sysconf v0.3.11 // indirect
	github.com/tklauser/numcpus v0.6.0 // indirect
	go.etcd.io/etcd/api/v3 v3.5.8 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.8 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.14.0 // indirect
	go.opentelemetry.io/otel/trace v1.14.0 // indirect
	go.opentelemetry.io/proto/otlp v0.19.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/mod v0.7.0 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/oauth2 v0.0.0-20220223155221-ee480838109b // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/term v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	golang.org/x/time v0.0.0-20220210224613-90d013bbcef8 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/api v0.26.3 // indirect
	k8s.io/apimachinery v0.26.3 // indirect
	k8s.io/klog/v2 v2.80.1 // indirect
	k8s.io/kube-openapi v0.0.0-20221012153701-172d655c2280 // indirect
	k8s.io/utils v0.0.0-20221107191617-1a15be271d1d // indirect
	sigs.k8s.io/json v0.0.0-20220713155537-f223a00ba0e2 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)
