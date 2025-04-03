package operation

import (
	"fmt"

	"github.com/Fernando-Dourado/harness-move-project/services"
	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
)

type (
	// NOT SURE WHICH NAME TO CHOSE TO THAT TYPE
	Config struct {
		Token   string
		Account string
		Org     string
		Project string
	}

	Move struct {
		Source Config
		Target Config
	}
)

func (o *Move) Exec() error {

	client := resty.New()

	sourceApi := services.SourceRequest{
		Client:  client,
		Token:   o.Source.Token,
		Account: o.Source.Account,
	}
	targetApi := services.TargetRequest{
		Client:  client,
		Token:   o.Target.Token,
		Account: o.Target.Account,
	}

	// SOURCE AND TARGET MUST EXIST
	if err := sourceApi.ValidateSource(o.Source.Org, o.Source.Project); err != nil {
		return err
	}
	if err := targetApi.ValidateTarget(o.Target.Org, o.Target.Project); err != nil {
		return err
	}

	var operations []services.Operation
	operations = append(operations, services.NewVariableOperation(&sourceApi, &targetApi, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))
	operations = append(operations, services.NewFileStoreOperation(&sourceApi, &targetApi, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))
	operations = append(operations, services.NewEnvironmentOperation(&sourceApi, &targetApi, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))
	operations = append(operations, services.NewInfrastructureOperation(&sourceApi, &targetApi, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))
	operations = append(operations, services.NewServiceOperation(&sourceApi, &targetApi, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))
	operations = append(operations, services.NewServiceOverrideOperation(&sourceApi, &targetApi, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))
	operations = append(operations, services.NewTemplateOperation(&sourceApi, &targetApi, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))
	operations = append(operations, services.NewPipelineOperation(&sourceApi, &targetApi, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))
	operations = append(operations, services.NewInputsetOperation(&sourceApi, &targetApi, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))

	for _, op := range operations {
		if err := op.Move(); err != nil {
			return err
		}
	}

	fmt.Println(color.GreenString("Done"))
	return nil
}
