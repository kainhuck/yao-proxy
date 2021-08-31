package local

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	Port       int    `json:"port"`
	Key        string `json:"key"`
	RemoteHost string `json:"remote_host"`
	RemotePort int    `json:"remote_port"`
	Debug      bool   `json:"debug"`
	Method     string `json:"method"`
}

func ReadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	bts, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var cfg = Config{}
	err = json.Unmarshal(bts, &cfg)
	if err != nil {
		return nil, err
	}

	if cfg.RemoteHost == "" {
		return nil, fmt.Errorf("need remote host")
	}
	if cfg.RemotePort == 0 {
		cfg.RemotePort = 20807
	}
	if cfg.Port == 0 {
		cfg.Port = 20808
	}
	if cfg.Key == "" {
		cfg.Key = "Atu@&^_^&Ak1314$$"
	}
	return &cfg, nil
}
