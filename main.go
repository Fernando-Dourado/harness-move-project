package main

import (
	"fmt"
	"os"

	"github.com/Fernando-Dourado/harness-move-project/operation"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Harness Move"
	app.Usage = "Non-official Harness CLI to move project between organizations"
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:     "api-token",
			Usage:    "TODO: Explain the api-token usage",
			Required: true,
		},
		cli.StringFlag{
			Name:     "account",
			Usage:    "TODO: Explain the arg usage",
			Required: true,
		},
		cli.StringFlag{
			Name:     "source-org",
			Usage:    "TODO: Explain the arg usage",
			Required: true,
		},
		cli.StringFlag{
			Name:     "source-project",
			Usage:    "TODO: Explain the arg usage",
			Required: true,
		},
		cli.StringFlag{
			Name:     "target-org",
			Usage:    "TODO: Explain the arg usage",
			Required: true,
		},
		cli.StringFlag{
			Name:     "target-project",
			Usage:    "TODO: Explain the arg usage",
			Required: true,
		},
	}
	app.Run(os.Args)
}

func run(c *cli.Context) {
	mv := operation.Move{
		Config: operation.Config{
			Token:   c.String("api-token"),
			Account: c.String("account"),
		},
		Source: operation.NoName{
			Org:     c.String("source-org"),
			Project: c.String("source-project"),
		},
		Target: operation.NoName{
			Org:     c.String("target-org"),
			Project: c.String("target-project"),
		},
	}
	if err := mv.Exec(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
