// AVS 的命令行主入口, 处理命令行参数
package main

import (
	"github.com/urfave/cli/v2"
	"goplus/avs/config"
	"log"
	"os"
)

var (
	ConfigFileFlag = cli.StringFlag{
		Name:    config.ConfigFileFlag,
		Aliases: []string{"c"},
		Usage:   "Config file path",
	}
)

func main() {
	flags := []cli.Flag{
		&ConfigFileFlag,
	}

	app := &cli.App{
		Name:  "goplus-operator",
		Usage: "Operator for GoPlus",
		Commands: []*cli.Command{
			{
				Name:   "start",
				Usage:  "Start the operator",
				Action: start,
				Flags:  flags,
			},
			{
				Name:   "register-with-avs",
				Usage:  "Register current operator to GoPlusAVS",
				Action: registerWithAVS,
				Flags:  flags,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln("Application failed.", "Message:", err)
	}
}
