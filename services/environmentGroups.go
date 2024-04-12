package services

import (
	"encoding/json"
	"fmt"

	"github.com/jf781/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
)

const ENVGROUPLIST = "/ng/api/environmentGroup/list"
const ENVGROUP = "/ng/api/environmentGroup"

type EnvGroupContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
}

func NewEnvGroupOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string) EnvGroupContext {
	return EnvGroupContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
	}
}

func (c EnvGroupContext) Move() error {

	// Leveraging listPipelines func from pipeline.go file
	envGroups, err := c.api.listEnvGroups(c.sourceOrg, c.sourceProject)
	if err != nil {
		return err
	}

	bar := progressbar.Default(int64(len(envGroups)), "Triggers    ")
	var failed []string

	for _, eg := range envGroups {

		e := model.CreateEnvGroup{}

		newYaml := createYaml(eg.EnvGroup.YAML, c.sourceOrg, c.sourceProject, c.targetOrg, c.targetProject)
		e.OrgIdentifier = c.targetOrg
		e.ProjectIdentifier = c.targetProject
		e.Color = eg.EnvGroup.Color
		e.Identifier = eg.EnvGroup.Identifier
		e.YAML = newYaml

		err = c.api.createEnvGroup(e)

		if err != nil {
			failed = append(failed, fmt.Sprintln(eg.EnvGroup.Name, "-", err.Error()))
		}
		bar.Add(1)
	}
	bar.Finish()

	reportFailed(failed, "Triggers:")
	return nil
}

func (api *ApiRequest) listEnvGroups(org, project string) ([]model.EnvGroupContent, error) {

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     org,
			"projectIdentifier": project,
			"size":              "100",
		}).
		Post(BaseURL + ENVGROUPLIST)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse(resp)
	}

	result := model.GetEnvGroupResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	return result.Data.Content, nil
}

func (api *ApiRequest) createEnvGroup(envGroup model.CreateEnvGroup) error {

	//api.Client.SetDebug(true)

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(envGroup).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
		}).
		Post(BaseURL + ENVGROUP)

	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}

	return nil
}
