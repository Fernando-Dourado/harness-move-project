package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/jf781/harness-move-project/operation"
	"github.com/urfave/cli/v2"
)

var Version = "development"

func main() {
	app := &cli.App{
		Name: "harness-copy-project",
		Version: Version,
		Usage: "Non-official Harness CLI to copy project between organizations",
		Action: run,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "csvPath",
				Usage:    "The path to the CSV file that contains the source and target project information.",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "apiToken",
				Usage: "The API token that will be used to authenticate with the Harness Account.",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "accountId",
				Usage: "The account ID that contains both the source and target orgnaizations.",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "baseUrl",
				Usage: "The URL of the harness instance that your projects reside in.",
				Required: true,
			},
		},
	}
	app.Run(os.Args)
}



func run(c *cli.Context) error {
	importCsv := operation.ImportCSV{
		CsvPath: c.String("csvPath"),
	}

	csvData, err := importCsv.Exec()
	if err != nil {
		fmt.Println(color.RedString(fmt.Sprint("Failed to pull CSV data:", err.Error())))
	}

	for i := 0; i < len(csvData.SourceOrg); i++ {
		mv := operation.Move{
			Config: operation.Config{
				Token:   c.String("apiToken"),
				Account: c.String("accountId"),
				BaseURL: c.String("baseUrl"),
			},
			Source: operation.NoName{
				Org:     csvData.SourceOrg[i],
				Project: csvData.SourceProject[i],
			},
			Target: operation.NoName{
				Org:     csvData.TargetOrg[i],
				Project: csvData.TargetProject[i],
			},
		}

		 // Check for missing or empty values
		 if mv.Source.Org == "" || mv.Source.Project == "" || mv.Target.Org == "" {
			fmt.Println("Invalid CSV data. Missing required fields.")
			continue // Skip this iteration if required data is missing
		}

		// Use source project name if target project name is missing
		if mv.Target.Project == "" {
			mv.Target.Project = mv.Source.Project
		}

		fmt.Println(color.GreenString("Moving project '%v' from org '%v' to org '%v'. The target project will be named '%v'", mv.Source.Project, mv.Source.Org, mv.Target.Org, mv.Target.Project)) 

		// Execute the copy operation from source to target operation
		if err := mv.Exec(); err != nil {
			fmt.Println(color.RedString(fmt.Sprint("Failed to Copy Project:", err.Error())))
			os.Exit(1)
		}
	}

	return nil

}