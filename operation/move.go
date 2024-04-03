package operation

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"github.com/jf781/harness-move-project/services"
)

type (
	Config struct {
		Token   string
		Account string
	}

	// NOT SURE WHICH NAME TO CHOSE TO THAT TYPE
	NoName struct {
		Org     string
		Project string
	}

	Move struct {
		Config Config
		Source NoName
		Target NoName
	}
)

func (o *Move) Exec() error {

	api := services.ApiRequest{
		Client:  resty.New(),
		Token:   o.Config.Token,
		Account: o.Config.Account,
	}

	// SOURCE AND TARGET MUST EXIST
	if err := api.ValidateProject(o.Source.Org, o.Source.Project); err != nil {
		return err
	}
	if err := api.ValidateProject(o.Target.Org, o.Target.Project); err != nil {
		return err
	}

	var operations []services.Operation
	// operations = append(operations, services.NewVariableOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))
	// operations = append(operations, services.NewEnvironmentOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))
	// operations = append(operations, services.NewInfrastructureOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))
	// operations = append(operations, services.NewServiceOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))
	// operations = append(operations, services.NewTemplateOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))
	// operations = append(operations, services.NewPipelineOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))
	// operations = append(operations, services.NewInputsetOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))
	// operations = append(operations, services.NewSecretOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))
	operations = append(operations, services.NewAccessControlRoles(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))

	for _, op := range operations {
		if err := op.Move(); err != nil {
			return err
		}
	}

	fmt.Println(color.GreenString("Done"))
	return nil
}
