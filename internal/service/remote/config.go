package remote

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	Port  int    `json:"port"`
	Key   string `json:"key"`
	Debug bool   `json:"debug"`
}

var defaultCfg = &Config{
	Port: 20807,
	Key:  "Atu@&^_^&Ak1314$$",
}

func ReadConfig(path string) *Config {
	f, err := os.Open(path)
	if err != nil {
		return defaultCfg
	}

	bts, err := io.ReadAll(f)
	if err != nil {
		return defaultCfg
	}

	var cfg = Config{}
	err = json.Unmarshal(bts, &cfg)
	if err != nil {
		return defaultCfg
	}

	if cfg.Port == 0 {
		cfg.Port = 20807
	}
	if cfg.Key == "" {
		cfg.Key = "Atu@&^_^&Ak1314$$"
	}
	return &cfg
}
