package main

import (
	"fmt"
	"os"

	"github.com/Fernando-Dourado/harness-move-project/operation"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

var Version = "development"

func main() {
	app := cli.NewApp()
	app.Name = "harness-move-project"
	app.Version = Version
	app.Usage = "Non-official Harness CLI to move project between organizations or accounts."
	app.UsageText = "harness-move-project [options]"
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:     "api-token",
			Usage:    "API authentication token for accessing the source system.",
			Required: true,
		},
		cli.StringFlag{
			Name:     "account",
			Usage:    "The account identifier associated with the source system.",
			Required: true,
		},
		cli.StringFlag{
			Name:     "source-org",
			Usage:    "The organization identifier in the source account.",
			Required: true,
		},
		cli.StringFlag{
			Name:     "source-project",
			Usage:    "The project identifier in the source account.",
			Required: true,
		},
		cli.StringFlag{
			Name:     "target-org",
			Usage:    "The org identifier in the target account.",
			Required: true,
		},
		cli.StringFlag{
			Name:     "target-project",
			Usage:    "The project identifier in the target account.",
			Required: false,
		},
		cli.StringFlag{
			Name:     "target-token",
			Usage:    "API authentication token for accessing the target system. Not required if the target account is the same as the source account.",
			Required: false,
		},
		cli.StringFlag{
			Name:     "target-account",
			Usage:    "The account identifier associated with the target system. Not required if the target account is the same as the source account.",
			Required: false,
		},
		cli.StringFlag{
			Name:     "create-project",
			Usage:    "Creates the project in the target account and org if it does not exist.",
			Required: false,
		},
	}
	app.Run(os.Args)
}

func run(c *cli.Context) {
	mv := operation.NewMove(
		operation.CopyConfig{
			Org:     c.String("source-org"),
			Project: c.String("source-project"),
			Token:   c.String("api-token"),
			Account: c.String("account"),
		},
		operation.CopyConfig{
			Org:     c.String("target-org"),
			Project: c.String("target-project"),
			Token:   c.String("target-token"),
			Account: c.String("target-account"),
		},
		operation.OperationConfig{
			CreateProject: c.Bool("create-project"),
		},
	)

	applyArgumentRules(mv)

	if err := mv.Exec(); err != nil {
		fmt.Println(color.RedString(fmt.Sprint("Failed:", err.Error())))
		os.Exit(1)
	}
}

func applyArgumentRules(mv *operation.Move) {
	// USE SOURCE PROJECT AS TARGET, WHEN TARGET NOT SET
	if len(mv.Target.Project) == 0 {
		mv.Target.Project = mv.Source.Project
	}
	// USE TOKEN AND ACCOUNT FROM SOURCE, WHEN TARGET NOT SET
	if len(mv.Target.Token) == 0 {
		mv.Target.Token = mv.Source.Token
	}
	if len(mv.Target.Account) == 0 {
		mv.Target.Account = mv.Source.Account
	}
}
