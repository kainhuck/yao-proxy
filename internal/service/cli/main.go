package cli

import (
	CLI "github.com/kainhuck/shu-cli"
	"github.com/kainhuck/yao-proxy/internal/service/cli/cmd"
)

func Main() {
	cli := CLI.DefaultCli()
	cli.SetWelcomeMsg(logo)
	cli.Register(cmd.ListCmd)
	cli.Run()
}


var logo = `██╗   ██╗ █████╗  ██████╗ 
╚██╗ ██╔╝██╔══██╗██╔═══██╗
 ╚████╔╝ ███████║██║   ██║
  ╚██╔╝  ██╔══██║██║   ██║
   ██║   ██║  ██║╚██████╔╝
   ╚═╝   ╚═╝  ╚═╝ ╚═════╝`