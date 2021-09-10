package cmd

import (
	CLI "github.com/kainhuck/shu-cli"
	"github.com/kainhuck/shu-cli/cmd"
)

var PrintCmd = &cmd.Command{
	Cmd:   "print",
	Usage: "print <local>/<remote>",
	Desc:  "print config file for local or remote",
	Handler: func(args ...string) {
		if len(args) == 0 {
			CLI.Println(CLI.Red("必须指定 local 或者 remote"))
			return
		}
		switch args[0] {
		case "local", "remote":
		default:
			CLI.Printf("unknown field `%s`\n", args[0])
			return
		}

		file, ok := CLI.Load(args[0])
		if !ok {
			CLI.Println("请先生成`%s`的配置文件!", args[0])
			return
		}

		CLI.Println(file)
	},
}
