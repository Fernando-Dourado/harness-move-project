package services

import (
	"encoding/json"
	"fmt"

	"github.com/Fernando-Dourado/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
)

type InfrastructureContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
}

func NewInfrastructureOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string) InfrastructureContext {
	return InfrastructureContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
	}
}

func (c InfrastructureContext) Move() error {

	envs, err := c.api.listEnvironments(c.sourceOrg, c.sourceProject)
	if err != nil {
		return err
	}

	bar := progressbar.Default(int64(len(envs)), "Infrastructure")
	var failed []string

	for _, env := range envs {
		e := env.Environment
		infras, err := c.api.listInfraDef(c.sourceOrg, c.sourceProject, e.Identifier)
		if err != nil {
			failed = append(failed, fmt.Sprintf("Unable to list infrastructures for environment %s [%s]", env.Environment.Name, err))
			continue
		}

		bar.ChangeMax(bar.GetMax() + len(infras))

		for _, infra := range infras {
			i := infra.Infrastructure
			newYaml := createYaml(i.Yaml, c.sourceOrg, c.sourceProject, c.targetOrg, c.targetProject)

			err := c.api.createInfrastructure(&model.CreateInfrastructureRequest{
				Name:              i.Name,
				Identifier:        i.Identifier,
				OrgIdentifier:     c.targetOrg,
				ProjectIdentifier: c.targetProject,
				Description:       i.Description,
				EnvironmentRef:    e.Identifier,
				DeploymentType:    i.DeploymentType,
				Type:              i.Type,
				Yaml:              newYaml,
			})
			if err != nil {
				failed = append(failed, fmt.Sprintln(e.Name, "/", i.Name, "-", err.Error()))
			}
			bar.Add(1)
		}
		bar.Add(1)
	}
	bar.Finish()

	reportFailed(failed, "infrastructures:")
	return nil
}

func (api *ApiRequest) listInfraDef(org, project, envId string) ([]*model.InfraDefListContent, error) {

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"accountIdentifier":     api.Account,
			"orgIdentifier":         org,
			"projectIdentifier":     project,
			"environmentIdentifier": envId,
			"size":                  "1000",
		}).
		Get(BaseURL + "/ng/api/infrastructures")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse(resp)
	}

	result := model.InfraDefListResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	return result.Data.Content, nil
}

func (api *ApiRequest) createInfrastructure(infra *model.CreateInfrastructureRequest) error {

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(infra).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
		}).
		Post(BaseURL + "/ng/api/infrastructures")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}

	return nil
}
