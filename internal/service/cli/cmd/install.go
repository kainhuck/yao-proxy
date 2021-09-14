package cmd

import (
	"encoding/json"
	"fmt"
	CLI "github.com/kainhuck/shu-cli"
	"github.com/kainhuck/shu-cli/cmd"
	"github.com/kainhuck/yao-proxy/internal/cipher"
	"github.com/kainhuck/yao-proxy/internal/service/remote"
)

/*
todo:
	1. 获取服务端的IP地址
	2. 生成local端的配置文件
	3. 完成systemd部署或二进制文件部署（需要获取服务器的信息，需要提前准备好对应的二进制文件）
*/

var InstallCmd = &cmd.Command{
	Cmd:   "install",
	Usage: "install",
	Desc:  "install remote server",
	Handler: func(args ...string) {
		remoteCfg := new(remote.Config)
		debugStr := CLI.ReadOne("远程服务是否需要开启Debug模式(输入Y/N)")
		switch debugStr {
		case "Y":
			remoteCfg.Debug = true
		default:
			remoteCfg.Debug = false
		}

		number := CLI.ReadOneInt("要开启几个服务器进程")
		for number <= 0 {
			number = CLI.ReadOneInt("不可小于1，请重新输入")
		}

		remoteCfg.ServerInfos = make([]remote.ServerInfo, number)

		for i := 0; i < number; i++ {
			remoteCfg.ServerInfos[i].Port = CLI.ReadOneInt(fmt.Sprintf("请输入第 %d 个进程的端口", i+1))
			method := CLI.ReadOne(fmt.Sprintf("请输入第 %d 个进程的加密方法", i+1))

			for {
				if _, ok := cipher.Methods[method]; !ok {
					method = CLI.ReadOne("不支持的加密方法，请重新输入")
				} else {
					break
				}
			}
			remoteCfg.ServerInfos[i].Method = method
			key := CLI.ReadOne(fmt.Sprintf("请输入第 %d 个进程的秘钥", i+1))
			for len(key) == 0 {
				key = CLI.ReadOne("秘钥不可为空，请重新输入")
			}
			remoteCfg.ServerInfos[i].Key = key
		}

		bts, _ := json.MarshalIndent(remoteCfg, "", "  ")

		CLI.Store("remote", string(bts))
		CLI.Printf("安装成功, 通过`%s`命令查看配置文件\n", CLI.Cyan("print"))
	},
}
