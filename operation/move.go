package operation

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"github.com/jf781/harness-move-project/services"
	"go.uber.org/zap"
)

type (
	Config struct {
		Token   string
		Account string
		BaseURL string
		Logger  *zap.Logger
	}

	// NOT SURE WHICH NAME TO CHOSE TO THAT TYPE
	NoName struct {
		Org     string
		Project string
	}

	Copy struct {
		Config Config
		Source NoName
		Target NoName
	}
)

func (o *Copy) Exec() error {

	api := services.ApiRequest{
		Client:  resty.New(),
		Token:   o.Config.Token,
		Account: o.Config.Account,
		BaseURL: o.Config.BaseURL,
	}

	var operations []services.Operation

	// SOURCE PORJECT MUST EXIST.  RETURNS AN ERROR IF CAN'T BE FOUND/DOES NOT EXIST.
	if err := api.ValidateProject(o.Source.Org, o.Source.Project); err != nil {
		return err
	}
	if err := api.ValidateProject(o.Target.Org, o.Target.Project); err != nil {
		// CREATE NEW PROJECT IF IT DOES NOT EXIST IN THE TARGET ORG
		operations = append(operations, services.NewProjectOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))
	}

	// operations = append(operations, services.NewVariableOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))
	operations = append(operations, services.NewConnectorOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project, o.Config.Logger))
	// operations = append(operations, services.NewFileStoreOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))
	// operations = append(operations, services.NewEnvironmentOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))
	// operations = append(operations, services.NewEnvGroupOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))
	// operations = append(operations, services.NewInfrastructureOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))
	// operations = append(operations, services.NewServiceOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))
	// operations = append(operations, services.NewServiceOverrideOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))
	// operations = append(operations, services.NewTemplateOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))
	// operations = append(operations, services.NewPipelineOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))
	// operations = append(operations, services.NewInputsetOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))
	// // operations = append(operations, services.NewSecretOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))
	// operations = append(operations, services.NewTagOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))
	// operations = append(operations, services.NewUserScopeOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))
	// operations = append(operations, services.NewUserGroupOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))
	// operations = append(operations, services.NewServiceAccountOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))
	// operations = append(operations, services.NewRoleOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))
	// operations = append(operations, services.NewResourceGroupOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))
	// operations = append(operations, services.NewRoleAssignmentOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))
	// operations = append(operations, services.NewTriggerOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))
	// operations = append(operations, services.NewFeatureFlagOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))
	// operations = append(operations, services.NewTargets(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))
	// operations = append(operations, services.NewTargetGroups(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))

	for _, op := range operations {
		if err := op.Copy(); err != nil {
			return err
		}
	}

	fmt.Println(color.GreenString("Done"))
	return nil
}
