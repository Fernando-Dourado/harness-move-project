package services

import (
	"encoding/json"
	"fmt"

	"github.com/Fernando-Dourado/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
)

type ServiceOverrideContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
}

func NewServiceOverrideOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string) ServiceOverrideContext {
	return ServiceOverrideContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
	}
}

func (c ServiceOverrideContext) Move() error {

	envs, err := c.api.listEnvironments(c.sourceOrg, c.sourceProject)
	if err != nil {
		return err
	}

	bar := progressbar.Default(int64(len(envs)), "Service Override")
	var failed []string

	for _, env := range envs {
		e := env.Environment
		overrides, err := c.api.listServiceOverrides(c.sourceOrg, c.sourceProject, e.Identifier)
		if err != nil {
			failed = append(failed, fmt.Sprintf("Unable to list service overrides for environment %s [%s]", env.Environment.Name, err))
			continue
		}

		bar.ChangeMax(bar.GetMax() + len(overrides))

		for _, o := range overrides {
			if len(o.YAML) == 0 {
				failed = append(failed, fmt.Sprintf("The YAML is empty [envId=%s,serviceRef=%s]", o.EnvironmentRef, o.ServiceRef))
			} else {
				err := c.api.createServiceOverride(&model.CreateServiceOverrideRequest{
					OrgIdentifier:     c.targetOrg,
					ProjectIdentifier: c.targetProject,
					EnvironmentRef:    o.EnvironmentRef,
					ServiceRef:        o.ServiceRef,
					YAML:              o.YAML,
				})
				if err != nil {
					failed = append(failed, fmt.Sprintln(e.Name, "/", o.ServiceRef, "-", err.Error()))
				}
			}
			bar.Add(1)
		}
		bar.Add(1)
	}
	bar.Finish()

	reportFailed(failed, "overrides:")
	return nil
}

func (api *ApiRequest) listServiceOverrides(org, project, envId string) ([]*model.ServiceOverride, error) {

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetQueryParams(map[string]string{
			"accountIdentifier":     api.Account,
			"orgIdentifier":         org,
			"projectIdentifier":     project,
			"environmentIdentifier": envId,
			"size":                  "1000",
		}).
		Get(BaseURL + "/ng/api/environmentsV2/serviceOverrides")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse(resp)
	}

	result := model.ListServiceOverridesRequest{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	return result.Data.Content, nil
}

func (api *ApiRequest) createServiceOverride(override *model.CreateServiceOverrideRequest) error {

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(override).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
		}).
		Post(BaseURL + "/ng/api/environmentsV2/serviceOverrides")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}

	return nil
}
