package cmd

import (
	CLI "github.com/kainhuck/shu-cli"
	"github.com/kainhuck/shu-cli/cmd"
	"github.com/kainhuck/yao-proxy/internal/cipher"
)

var ListCmd = &cmd.Command{
	Cmd:   "list",
	Usage: "list",
	Desc:  "list all supported encrypt method",
	Handler: func(args ...string) {
		index := 1
		for k := range cipher.Methods {
			CLI.Printf(" %d. %s\n", index, CLI.Blue(k))
			index++
		}
	},
}
