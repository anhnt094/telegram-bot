package cmd

import (
	"github.com/urfave/cli/v2"
	"os"
	"sort"
	"time"
)

func Execute() {
	app := cli.App{
		Name:      "anhnt094-bot",
		Usage:     "Personal Bot of Nguyen Tuan Anh (anh.nt094@gmail.com)",
		ArgsUsage: "",
		Version:   "v1.0.0",
		Compiled:  time.Time{},
		Authors: []*cli.Author{{
			Name:  "Nguyen Tuan Anh",
			Email: "anh.nt094@gmail.com",
		}},
	}

	app.Commands = []*cli.Command{
		{
			Name:  "listen",
			Usage: "Listen updates from Telegram",
			Action: func(c *cli.Context) error {
				if err := listen(); err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name:  "analyze-k8s",
			Usage: "Run kubernetes analysis",
			Action: func(c *cli.Context) error {
				if err := analyzeK8s(); err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name:  "get-usdt-price",
			Usage: "Get highest usdt price on P2P Binance",
			Action: func(c *cli.Context) error {
				if err := getUsdtPrice(); err != nil {
					return err
				}
				return nil
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
