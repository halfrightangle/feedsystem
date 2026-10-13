package config

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

const defaultConfigPath = "configs/config.yaml"

type fileConfig struct {
	Env    string     `yaml:"env"`
	Server fileServer `yaml:"server"`
}

type fileServer struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

func Load() (Config, error) {
	_ = godotenv.Load()

	path, explicit := configPath()
	fc, err := loadFile(path, explicit)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{Env: env("ENVIRONMENT", fc.Env)}
	if cfg.Env == "" {
		cfg.Env = "development"
	}

	cfg.Server.Port = 8080
	if cfg.IsProduction() {
		cfg.Server.Mode = "release"
	} else {
		cfg.Server.Mode = "debug"
	}

	var errs []error

	if fc.Server.Port != 0 {
		cfg.Server.Port = fc.Server.Port
	}
	if fc.Server.Mode != "" {
		cfg.Server.Mode = fc.Server.Mode
	}

	if v := os.Getenv("GIN_MODE"); v != "" {
		cfg.Server.Mode = v
	}

	errs = append(errs, cfg.validate()...)

	if err := errors.Join(errs...); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func configPath() (path string, explicit bool) {
	if v := os.Getenv("CONFIG_FILE"); v != "" {
		return v, true
	}
	return defaultConfigPath, false
}

func loadFile(path string, required bool) (fileConfig, error) {
	var fc fileConfig

	f, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		if required {
			return fc, fmt.Errorf("配置文件 %s 不存在：CONFIG_FILE 显式指定了它，不能当成「没有配置文件」", path)
		}
		return fc, nil
	}
	if err != nil {
		return fc, fmt.Errorf("打开配置文件 %s 失败: %w", path, err)
	}
	defer f.Close()

	dec := yaml.NewDecoder(f)

	dec.KnownFields(true)

	if err := dec.Decode(&fc); err != nil {
		if errors.Is(err, io.EOF) {
			return fc, nil
		}
		return fc, fmt.Errorf("解析配置文件 %s 失败: %w", path, err)
	}
	return fc, nil
}

func env(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}
