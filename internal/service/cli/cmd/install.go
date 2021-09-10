package cmd

import (
	CLI "github.com/kainhuck/shu-cli"
	"github.com/kainhuck/shu-cli/cmd"
)

var InstallCmd = &cmd.Command{
	Cmd:   "install",
	Usage: "install",
	Desc:  "install remote server",
	Handler: func(args ...string) {
		number := CLI.ReadOneInt("要开启几个服务器进程")
		for number == 0 {
			number = CLI.ReadOneInt("不可为零，请重新输入")
		}
		// todo
	},
}
