package remote

import (
	"encoding/json"
	"io"
	"os"
)

type ServerInfo struct {
	Port   int    `json:"port"`
	Key    string `json:"key"`
	Method string `json:"method"`
}

type Config struct {
	Debug       bool         `json:"debug"`
	ServerInfos []ServerInfo `json:"server_infos"`
}

var defaultCfg = &Config{
	Debug: false,
	ServerInfos: []ServerInfo{{
		Port:   20807,
		Key:    "Atu@&^_^&Ak1314$$",
		Method: "aes-128-cfb",
	}},
}

func ReadConfig(path string) *Config {
	f, err := os.Open(path)
	if err != nil {
		return defaultCfg
	}
	defer f.Close()

	bts, err := io.ReadAll(f)
	if err != nil {
		return defaultCfg
	}

	var cfg = Config{}
	err = json.Unmarshal(bts, &cfg)
	if err != nil {
		return defaultCfg
	}

	if len(cfg.ServerInfos) == 0 {
		cfg.ServerInfos = defaultCfg.ServerInfos
	}

	return &cfg
}
