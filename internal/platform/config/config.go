package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name string `mapstructure:"name"`
		Env  string `mapstructure:"env"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"app"`
	Logging struct {
		Level  string `mapstructure:"level"`
		Format string `mapstructure:"format"`
	} `mapstructure:"logging"`
	External struct {
		Pharmacy struct {
			BaseURL string `mapstructure:"base_url"`
			Path    string `mapstructure:"path"`
			UseMock bool   `mapstructure:"use_mock"`
			APIKey  string `mapstructure:"api_key"`
			Timeout string `mapstructure:"timeout"`
		} `mapstructure:"pharmacy"`
		Billing struct {
			BaseURL string `mapstructure:"base_url"`
			Path    string `mapstructure:"path"`
			UseMock bool   `mapstructure:"use_mock"`
			Timeout string `mapstructure:"timeout"`
		} `mapstructure:"billing"`
	} `mapstructure:"external"`
}

func Load() *Config {
	v := viper.New()
	v.SetConfigName("app")
	v.SetConfigType("yaml")
	v.AddConfigPath("./internal/configs")
	_ = v.ReadInConfig()

	// optional env-specific file: RX_APP_ENV=prod -> app.prod.yaml
	if env := os.Getenv("RX_APP_ENV"); env != "" {
		v2 := viper.New()
		v2.SetConfigName("app." + env)
		v2.AddConfigPath("./internal/configs")
		v2.SetConfigType("yaml")
		if err := v2.ReadInConfig(); err == nil {
			v.MergeConfigMap(v2.AllSettings())
		}
	}

	v.SetEnvPrefix("RX")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	cfg := &Config{}
	_ = v.Unmarshal(cfg)
	if cfg.App.Port == 0 {
		cfg.App.Port = 8080
	}
	if cfg.App.Name == "" {
		cfg.App.Name = "rxintake"
	}
	if cfg.App.Env == "" {
		cfg.App.Env = "dev"
	}
	if cfg.Logging.Level == "" {
		cfg.Logging.Level = "debug"
	}
	if cfg.Logging.Format == "" {
		cfg.Logging.Format = "console"
	}
	return cfg
}
