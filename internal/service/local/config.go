package local

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type RemoteInfo struct {
	RemoteAddr string `json:"remote_addr"`
	Method     string `json:"method"`
	Key        string `json:"key"`
}

type ServerInfo struct {
	Port        int          `json:"port"`
	RemoteInfos []RemoteInfo `json:"remote_infos"`
}

type Config struct {
	Debug       bool         `json:"debug"`
	ServerInfos []ServerInfo `json:"server_infos"`
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

	if len(cfg.ServerInfos) == 0 {
		return nil, fmt.Errorf("need server_infos")
	}

	return &cfg, nil
}
