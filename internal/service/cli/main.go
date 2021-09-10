package cli

import (
	CLI "github.com/kainhuck/shu-cli"
	"github.com/kainhuck/yao-proxy/internal/service/cli/cmd"
)

func Main() {
	cli := CLI.DefaultCli()
	cli.SetWelcomeMsg(logo)

	cli.Register(cmd.ListCmd)
	cli.Register(cmd.PrintCmd)
	cli.Register(cmd.InstallCmd)

	cli.Run()
}

var logo = `██╗   ██╗ █████╗  ██████╗ 
╚██╗ ██╔╝██╔══██╗██╔═══██╗
 ╚████╔╝ ███████║██║   ██║
  ╚██╔╝  ██╔══██║██║   ██║
   ██║   ██║  ██║╚██████╔╝
   ╚═╝   ╚═╝  ╚═╝ ╚═════╝`
