package bootstrap

import "flag"

type CommandFlags struct {
	Conf       string // Boot configuration file path, default: ../../configs
	Env        string // Development environment: dev, debug...
	ConfigHost string // Remote configuration server address
	ConfigType string // Remote configuration server type
	Daemon     bool   // Whether to convert to daemon process
}

func NewCommandFlags() *CommandFlags {
	return &CommandFlags{
		Conf:       "",
		Env:        "",
		ConfigHost: "",
		ConfigType: "",
		Daemon:     false,
	}
}

func (f *CommandFlags) Init() {
	flag.StringVar(&f.Conf, "conf", "../../configs", "config path, eg: -conf ../../configs")
	flag.StringVar(&f.Env, "env", "dev", "runtime environment, eg: -env dev")
	flag.StringVar(&f.ConfigHost, "chost", "127.0.0.1:8500", "config server host, eg: -chost 127.0.0.1:8500")
	flag.StringVar(&f.ConfigType, "ctype", "consul", "config server host, eg: -ctype consul")
	flag.BoolVar(&f.Daemon, "d", false, "run app as a daemon with -d=true.")

	if f.Daemon {
		BeDaemon("-d")
	}
	flag.Parse()
}
