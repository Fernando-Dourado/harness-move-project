package operation

import (
	"errors"
	"fmt"

	"github.com/Fernando-Dourado/harness-move-project/services"
	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
)

type (
	OperationConfig struct {
		CreateProject bool
	}

	CopyConfig struct {
		Token   string
		Account string
		Org     string
		Project string
	}

	Move struct {
		Source CopyConfig
		Target CopyConfig
		Config OperationConfig
	}
)

func NewMove(s, t CopyConfig, c OperationConfig) *Move {
	return &Move{
		Source: s,
		Target: t,
		Config: c,
	}
}

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
		if err = o.createProjectWhenRequired(&sourceApi, &targetApi, err); err != nil {
			return err
		}
	}

	st := &services.SourceTarget{
		SourceOrg:     o.Source.Org,
		SourceProject: o.Source.Project,
		TargetOrg:     o.Target.Org,
		TargetProject: o.Target.Project,
	}

	var operations []services.Operation
	operations = append(operations, services.NewVariableOperation(&sourceApi, &targetApi, st))
	operations = append(operations, services.NewSecretOperation(&sourceApi, &targetApi, st))
	operations = append(operations, services.NewConnectorOperation(&sourceApi, &targetApi, st))
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

func (o *Move) createProjectWhenRequired(sourceApi *services.SourceRequest, targetApi *services.TargetRequest, err error) error {

	if errors.Is(err, services.ErrEntityNotFound) {
		if o.Config.CreateProject {
			fmt.Println("Creating project in target...")

			err = services.NewProjectOperation(sourceApi, targetApi, &services.SourceTarget{
				SourceOrg:     o.Source.Org,
				SourceProject: o.Source.Project,
				TargetOrg:     o.Target.Org,
				TargetProject: o.Target.Project,
			}).Move()

			if err == nil {
				fmt.Println(color.GreenString("Project %s created in target org %s", o.Target.Project, o.Target.Org))
				return nil
			}
		} else {
			err = fmt.Errorf("project %s not found in target org %s; create project not set", o.Target.Project, o.Target.Org)
		}
	}
	return err
}
