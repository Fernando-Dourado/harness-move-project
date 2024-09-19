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

	// Initlize and configure the logger
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel) // Set to Info, Debug, or Error for more verbose logging
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
			&cli.BoolFlag{
				Name:     "copyCDComponents",
				Usage:    "The if set to 'true', then it will copy the Continuous Delivery components.",
				Required: false,
				Value:    false,
			},
			&cli.BoolFlag{
				Name:     "copyFFComponents",
				Usage:    "The if set to 'true, then it will copy the Feature Flag components.",
				Required: false,
				Value:    false,
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
				CopyCD:  c.Bool("copyCDComponents"),
				CopyFF:  c.Bool("copyFFComponents"),
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
		} else {
			fmt.Println(color.GreenString("Project '%v' has been copied to '%v'", cp.Source.Project, cp.Target.Project))
			fmt.Println(color.GreenString("Connectors Total: %v ", services.GetConnectorsTotal()))
			fmt.Println(color.GreenString("Connectors Moved: %v ", services.GetConnectorsMoved()))
			fmt.Println(color.GreenString("Environments Total: %v ", services.GetEnvironmentsTotal()))
			fmt.Println(color.GreenString("Environments Moved: %v ", services.GetEnvironmentsMoved()))

			fmt.Println(color.GreenString("EnvironmentGroups Total: %v ", services.GetEnvironmentGroupsTotal()))
			fmt.Println(color.GreenString("EnvironmentGroups Moved: %v ", services.GetEnvironmentGroupsMoved()))

			fmt.Println(color.GreenString("FeatureFlags Total: %v ", services.GetFeatureFlagsTotal()))
			fmt.Println(color.GreenString("FeatureFlags Moved: %v ", services.GetFeatureFlagsMoved()))

			fmt.Println(color.GreenString("FileStores Total: %v ", services.GetFileStoresTotal()))
			fmt.Println(color.GreenString("FileStores Moved: %v ", services.GetFileStoresMoved()))

			fmt.Println(color.GreenString("Infrastructure Total: %v ", services.GetInfrastructureTotal()))
			fmt.Println(color.GreenString("Infrastructure Moved: %v ", services.GetInfrastructureMoved()))

			fmt.Println(color.GreenString("InputSets Total: %v ", services.GetInputSetsTotal()))
			fmt.Println(color.GreenString("InputSets Moved: %v ", services.GetInputSetsMoved()))

			fmt.Println(color.GreenString("Pipelines Total: %v ", services.GetPipelinesTotal()))
			fmt.Println(color.GreenString("Pipelines Moved: %v ", services.GetPipelinesMoved()))

			fmt.Println(color.GreenString("ResourceGroups Total: %v ", services.GetResourceGroupsTotal()))
			fmt.Println(color.GreenString("ResourceGroups Moved: %v ", services.GetResourceGroupsMoved()))

			fmt.Println(color.GreenString("RoleAssignments Total: %v ", services.GetRoleAssignmentsTotal()))
			fmt.Println(color.GreenString("RoleAssignments Moved: %v ", services.GetRoleAssignmentsMoved()))

			fmt.Println(color.GreenString("Roles Total: %v ", services.GetRolesTotal()))
			fmt.Println(color.GreenString("Roles Moved: %v ", services.GetRolesMoved()))

			fmt.Println(color.GreenString("Service Overrides Total: %v ", services.GetOverridesTotal()))
			fmt.Println(color.GreenString("Service Overrides Moved: %v ", services.GetOverridesMoved()))

			fmt.Println(color.GreenString("Services Total: %v ", services.GetServicesTotal()))
			fmt.Println(color.GreenString("Services Moved: %v ", services.GetServicesMoved()))

			fmt.Println(color.GreenString("Tags Total: %v ", services.GetTagsTotal()))
			fmt.Println(color.GreenString("Tags Moved: %v ", services.GetTagsMoved()))

			fmt.Println(color.GreenString("TargetGroups Total: %v ", services.GetTargetGroupsTotal()))
			fmt.Println(color.GreenString("TargetGroups Moved: %v ", services.GetTargetGroupsMoved()))

			fmt.Println(color.GreenString("Targets Total: %v ", services.GetTargetsTotal()))
			fmt.Println(color.GreenString("Targets Moved: %v ", services.GetTargetsMoved()))

			fmt.Println(color.GreenString("Templates Total: %v ", services.GetTemplatesTotal()))
			fmt.Println(color.GreenString("Templates Moved: %v ", services.GetTemplatesMoved()))

			fmt.Println(color.GreenString("UserGroups Total: %v ", services.GetUserGroupsTotal()))
			fmt.Println(color.GreenString("UserGroups Moved: %v ", services.GetUserGroupsMoved()))

			fmt.Println(color.GreenString("Users Total: %v ", services.GetUsersTotal()))
			fmt.Println(color.GreenString("Users Moved: %v ", services.GetUsersMoved()))

			fmt.Println(color.GreenString("Variables Total: %v ", services.GetVariablesTotal()))
			fmt.Println(color.GreenString("Variables Moved: %v ", services.GetVariablesMoved()))

			// Assuming you have a zap logger instance initialized as 'logger'

			logger.Info("Project Migration Status:",
				zap.Int("ConnectorsTotal", services.GetConnectorsTotal()),
				zap.Int("ConnectorsMoved", services.GetConnectorsMoved()),
				zap.Int("EnvironmentsTotal", services.GetEnvironmentsTotal()),
				zap.Int("EnvironmentsMoved", services.GetEnvironmentsMoved()),
				zap.Int("EnvironmentGroupsTotal", services.GetEnvironmentGroupsTotal()),
				zap.Int("EnvironmentGroupsMoved", services.GetEnvironmentGroupsMoved()),
				zap.Int("FeatureFlagsTotal", services.GetFeatureFlagsTotal()),
				zap.Int("FeatureFlagsMoved", services.GetFeatureFlagsMoved()),
				zap.Int("FileStoresTotal", services.GetFileStoresTotal()),
				zap.Int("FileStoresMoved", services.GetFileStoresMoved()),
				zap.Int("InfrastructureTotal", services.GetInfrastructureTotal()),
				zap.Int("InfrastructureMoved", services.GetInfrastructureMoved()),
				zap.Int("InputSetsTotal", services.GetInputSetsTotal()),
				zap.Int("InputSetsMoved", services.GetInputSetsMoved()),
				zap.Int("PipelinesTotal", services.GetPipelinesTotal()),
				zap.Int("PipelinesMoved", services.GetPipelinesMoved()),
				zap.Int("ResourceGroupsTotal", services.GetResourceGroupsTotal()),
				zap.Int("ResourceGroupsMoved", services.GetResourceGroupsMoved()),
				zap.Int("RoleAssignmentsTotal", services.GetRoleAssignmentsTotal()),
				zap.Int("RoleAssignmentsMoved", services.GetRoleAssignmentsMoved()),
				zap.Int("RolesTotal", services.GetRolesTotal()),
				zap.Int("RolesMoved", services.GetRolesMoved()),
				zap.Int("OverridesTotal", services.GetOverridesTotal()),
				zap.Int("OverridesMoved", services.GetOverridesMoved()),
				zap.Int("ServicesTotal", services.GetServicesTotal()),
				zap.Int("ServicesMoved", services.GetServicesMoved()),
				zap.Int("TagsTotal", services.GetTagsTotal()),
				zap.Int("TagsMoved", services.GetTagsMoved()),
				zap.Int("TargetGroupsTotal", services.GetTargetGroupsTotal()),
				zap.Int("TargetGroupsMoved", services.GetTargetGroupsMoved()),
				zap.Int("TargetsTotal", services.GetTargetsTotal()),
				zap.Int("TargetsMoved", services.GetTargetsMoved()),
				zap.Int("TemplatesTotal", services.GetTemplatesTotal()),
				zap.Int("TemplatesMoved", services.GetTemplatesMoved()),
				zap.Int("UserGroupsTotal", services.GetUserGroupsTotal()),
				zap.Int("UserGroupsMoved", services.GetUserGroupsMoved()),
				zap.Int("UsersTotal", services.GetUsersTotal()),
				zap.Int("UsersMoved", services.GetUsersMoved()),
				zap.Int("VariablesTotal", services.GetVariablesTotal()),
				zap.Int("VariablesMoved", services.GetVariablesMoved()),
			)

		}

		logger.Info(fmt.Sprintf("Project '%v' has been copied to '%v'", cp.Source.Project, cp.Target.Project))

		services.ResetAllCounters()

	}
	return nil
}
