package main

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/jf781/harness-move-project/operation"
	"github.com/jf781/harness-move-project/services"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Version = "development"

var logger *zap.Logger

func main() {

	// Initlize the logger
	// logger = zap.Must(zap.NewProduction())
	config := zap.NewProductionConfig()

	// Set the log level to WARN (or ERROR, if you prefer)
	config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	logger, _ = config.Build()

	startTime := time.Now()

	// Defer the logger and print the final log
	defer logger.Sync()
	defer func() {
		stopTime := time.Now()
		apiCalls := services.GetApiCalls()
		projects := services.GetProjects()
		duration := stopTime.Sub(startTime)

		var avgApiCallDuration time.Duration
		if apiCalls > 0 {
			avgApiCallDuration = duration / time.Duration(apiCalls)
		} else {
			avgApiCallDuration = 0
		}

		logger.Info("Harness Copy Project has completed.",
			zap.Int("Number of API Calls: ", apiCalls),
			zap.String("Run Duration: ", duration.String()),
			zap.Duration("Average API Call Duration: ", avgApiCallDuration),
			zap.Int("Number of projects moved: ", projects),
			zap.String("Stop Time: ", stopTime.Format("12:00:00")),
		)
	}()

	// Start logging
	logger.Info("Harness Copy Project has started.")
	logger.Info("Start time: " + startTime.Format("12:00:00"))

	// Create a new CLI app
	app := &cli.App{
		Name:    "harness-copy-project",
		Version: Version,
		Usage:   "Non-official Harness CLI to copy project between organizations",
		Action:  run,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "csvPath",
				Usage:    "The path to the CSV file that contains the source and target project information.",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "apiToken",
				Usage:    "The API token that will be used to authenticate with the Harness Account.",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "accountId",
				Usage:    "The account ID that contains both the source and target orgnaizations.",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "baseUrl",
				Usage:    "The URL of the harness instance that your projects reside in.",
				Required: true,
			},
		},
	}

	// Run the CLI app
	app.Run(os.Args)
}

func run(c *cli.Context) error {
	importCsv := operation.ImportCSV{
		CsvPath: c.String("csvPath"),
	}

	csvData, err := importCsv.Exec()
	if err != nil {
		logger.Error("Failed to pull CSV data",
			zap.String("csvPath", c.String("csvPath")),
			zap.Error(err),
		)
		return err
	}

	for i := 0; i < len(csvData.SourceOrg); i++ {
		// Increment the number of projects moved
		services.IncrementProjects()

		// Create a new copy operation
		cp := operation.Copy{
			Config: operation.Config{
				Token:   c.String("apiToken"),
				Account: c.String("accountId"),
				BaseURL: c.String("baseUrl"),
				Logger:  logger,
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
		if cp.Source.Org == "" || cp.Source.Project == "" || cp.Target.Org == "" {
			logger.Warn("Invalid CSV data. Missing required fields.",
				zap.String("Source Org", cp.Source.Org),
				zap.String("Source Project", cp.Source.Project),
				zap.String("Target Org", cp.Target.Org),
			)
			continue // Skip this iteration if required data is missing
		}

		// Use source project name if target project name is missing
		if cp.Target.Project == "" {
			cp.Target.Project = cp.Source.Project
		}

		fmt.Println(color.GreenString("Moving project '%v' from org '%v' to org '%v'. The target project will be named '%v'", cp.Source.Project, cp.Source.Org, cp.Target.Org, cp.Target.Project))

		// Execute the copy operation from source to target operation
		if err := cp.Exec(); err != nil {
			logger.Error("Failed to Copy Project",
				zap.String("Source Project", cp.Source.Project),
				zap.String("Target Project", cp.Target.Project),
				zap.Error(err),
			)
			return err
		}

		logger.Info(fmt.Sprintf("Project '%v' has been copied to '%v'", cp.Source.Project, cp.Target.Project))
	}
	return nil
}
