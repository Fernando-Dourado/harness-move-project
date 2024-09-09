package services

import (
	"encoding/json"
	"fmt"

	"github.com/jf781/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
)

const TARGETGROUPS = "/cf/admin/segments"

type TargetGroupContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
}

func NewTargetGroups(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string) TargetGroupContext {
	return TargetGroupContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
	}
}

func (c TargetGroupContext) Move() error {

	envs, err := c.api.listEnvironments(c.sourceOrg, c.sourceProject)
	if err != nil {
		return err
	}

	bar := progressbar.Default(int64(len(envs)), "Targets    ")
	var failed []string

	for _, env := range envs {
		e := env.Environment
		targetGroups, err := c.api.listTargetGroups(c.sourceOrg, c.sourceProject, e.Identifier)
		if err != nil {
			failed = append(failed, fmt.Sprintf("Unable to list targets for environment %s [%s]", env.Environment.Name, err))
			continue
		}

		bar.ChangeMax(bar.GetMax() + len(targetGroups))

		for _, targetGroup := range targetGroups {
			i := targetGroup
			// newYaml := createYaml(i.Yaml, c.sourceOrg, c.sourceProject, c.targetOrg, c.targetProject)

			includedNames := []string{}
			for _, included := range i.Included {
				iName := included.Name // Declare iName inside the loop
				includedNames = append(includedNames, iName)
			}

			err := c.api.createTargetGroups(&model.NewTargetGroup{
				Name:         i.Name,
				Identifier:   i.Identifier,
				Org:          c.targetOrg,
				Project:      c.targetProject,
				Account:      i.Account,
				Environment:  e.Identifier,
				Included:     includedNames,
				Excluded:     i.Excluded,
				Rules:        i.Rules,
				ServingRules: i.ServingRules,
			})
			if err != nil {
				failed = append(failed, fmt.Sprintln(e.Name, "/", i.Name, "-", err.Error()))
			}
			bar.Add(1)
		}
		bar.Add(1)
	}
	bar.Finish()

	reportFailed(failed, "target groups:")
	return nil
}

func (api *ApiRequest) listTargetGroups(org, project, envId string) ([]*model.TargetGroups, error) {

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
		Get(BaseURL + TARGETGROUPS)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse(resp)
	}

	result := model.TargetGroupListResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	targetGroups := []*model.TargetGroups{}
	for _, c := range result.TargetGroups {
		tg := c
		targetGroups = append(targetGroups, &tg)
	}

	return targetGroups, nil
}

func (api *ApiRequest) createTargetGroups(targetGroup *model.NewTargetGroup) error {

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(targetGroup).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     targetGroup.Org,
		}).
		Post(BaseURL + TARGETGROUPS)
	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}

	return nil
}
