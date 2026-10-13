package config

import (
	"fmt"
	"strconv"
)

type Config struct {
	Env    string
	Server Server
}

type Server struct {
	Port int
	Mode string
}

func (c Config) Addr() string { return ":" + strconv.Itoa(c.Server.Port) }

func (c Config) IsProduction() bool { return c.Env == "production" }

func (c Config) validate() []error {
	var errs []error

	if c.Server.Port < 1 || c.Server.Port > 65535 {
		errs = append(errs, fmt.Errorf("端口必须在 1-65535 之间，当前是 %d（检查 config.yaml 的 server.port）", c.Server.Port))
	}

	switch c.Env {
	case "development", "production":
	default:
		errs = append(errs, fmt.Errorf("env 只能是 development 或 production，当前是 %q", c.Env))
	}

	switch c.Server.Mode {
	case "debug", "release", "test":
	default:
		errs = append(errs, fmt.Errorf("server.mode 只能是 debug、release 或 test，当前是 %q", c.Server.Mode))
	}

	return errs
}
