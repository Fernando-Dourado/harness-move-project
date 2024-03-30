package services

import (
	"encoding/json"
	"fmt"

	"github.com/Fernando-Dourado/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
)

type EnvironmentContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
}

func NewEnvironmentOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string) EnvironmentContext {
	return EnvironmentContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
	}
}

func (c EnvironmentContext) Move() error {

	envs, err := c.api.listEnvironments(c.sourceOrg, c.sourceProject)
	if err != nil {
		return nil
	}

	bar := progressbar.Default(int64(len(envs)), "Environments")
	var failed []string

	for _, env := range envs {
		e := env.Environment

		newYaml := createYaml(e.Yaml, c.sourceOrg, c.sourceProject, c.targetOrg, c.targetProject)
		req := &model.CreateEnvironmentRequest{
			OrgIdentifier:     c.targetOrg,
			ProjectIdentifier: c.targetProject,
			Identifier:        e.Identifier,
			Name:              e.Name,
			Description:       e.Description,
			Color:             e.Color,
			Type:              e.Type,
			Yaml:              newYaml,
		}
		if err := c.api.createEnvironment(req); err != nil {
			failed = append(failed, fmt.Sprintln(e.Name, "-", err.Error()))
		}
		bar.Add(1)
	}
	bar.Finish()

	reportFailed(failed, "environments:")
	return nil
}

func (api *ApiRequest) listEnvironments(org, project string) ([]*model.ListEnvironmentContent, error) {

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     org,
			"projectIdentifier": project,
			"size":              "1000",
		}).
		Get(BaseURL + "/ng/api/environmentsV2")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse(resp)
	}

	result := model.ListEnvironmentResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	return result.Data.Content, nil
}

func (api *ApiRequest) createEnvironment(env *model.CreateEnvironmentRequest) error {

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(env).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
		}).
		Post(BaseURL + "/ng/api/environmentsV2")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}

	return nil
}
