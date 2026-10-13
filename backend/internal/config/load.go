package config

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

const defaultConfigPath = "configs/config.yaml"

type fileConfig struct {
	Env    string     `yaml:"env"`
	Server fileServer `yaml:"server"`
	Log    fileLog    `yaml:"log"`
}

type fileServer struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

type fileLog struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
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
		cfg.Server.Mode, cfg.Log.Format, cfg.Log.Level = "release", "json", slog.LevelInfo
	} else {
		cfg.Server.Mode, cfg.Log.Format, cfg.Log.Level = "debug", "text", slog.LevelDebug
	}

	var errs []error

	if fc.Server.Port != 0 {
		cfg.Server.Port = fc.Server.Port
	}
	if fc.Server.Mode != "" {
		cfg.Server.Mode = fc.Server.Mode
	}
	if fc.Log.Format != "" {
		cfg.Log.Format = fc.Log.Format
	}
	if fc.Log.Level != "" {
		if lv, err := parseLevel(path+" 的 log.level", fc.Log.Level); err != nil {
			errs = append(errs, err)
		} else {
			cfg.Log.Level = lv
		}
	}

	if v := os.Getenv("GIN_MODE"); v != "" {
		cfg.Server.Mode = v
	}
	if v := os.Getenv("LOG_FORMAT"); v != "" {
		cfg.Log.Format = v
	}
	if raw := os.Getenv("LOG_LEVEL"); raw != "" {
		if lv, err := parseLevel("LOG_LEVEL", raw); err != nil {
			errs = append(errs, err)
		} else {
			cfg.Log.Level = lv
		}
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

func parseLevel(where, raw string) (slog.Level, error) {
	switch strings.ToLower(raw) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	}
	return 0, fmt.Errorf("%s 只能是 debug/info/warn/error，当前值是 %q", where, raw)
}
