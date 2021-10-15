package middle

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// middle 作为中间节点角色，他扮演了`client`和`server`的部分角色
//// 所以middle需要有自己监听端口和加密方式，也需要知道下一个节点的地址和加密

// NextNodeInfo 下一个节点的信息
type NextNodeInfo struct {
	// Addr 下一个节点的地址
	Addr string `json:"addr"`
	// Method 下一额节点的加密方式
	Method string `json:"method"`
	// Key 下一个节点的秘钥
	Key string `json:"key"`
}

// ServerInfo 一个服务的信息
type ServerInfo struct {
	// Port 监听的端口
	Port int `json:"port"`
	// Method 加密方法
	Method string `json:"method"`
	// Key 解密秘钥
	Key string `json:"key"`
	// NextNodeInfos 下一跳地址
	NextNodeInfos []NextNodeInfo `json:"next_node_infos"`
}

// Config 中间服务的配置文件
type Config struct {
	// Debug 是否开启debug日志
	Debug bool `json:"debug"`
	// ServerInfos 服务信息
	ServerInfos []ServerInfo `json:"server_infos"`
}

// ReadConfig 解析配置文件
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
