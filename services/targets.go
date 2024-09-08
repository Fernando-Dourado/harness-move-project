package services

import (
	"encoding/json"
	"fmt"

	"github.com/jf781/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
)

const TARGETS = "/cf/admin/targets"

type TargetContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
}

func NewTargets(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string) TargetContext {
	return TargetContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
	}
}

func (c TargetContext) Move() error {

	envs, err := c.api.listEnvironments(c.sourceOrg, c.sourceProject)
	if err != nil {
		return err
	}

	bar := progressbar.Default(int64(len(envs)), "Target Groups    ")
	var failed []string

	for _, env := range envs {
		e := env.Environment
		targets, err := c.api.listTargets(c.sourceOrg, c.sourceProject, e.Identifier)
		if err != nil {
			failed = append(failed, fmt.Sprintf("Unable to list targets for environment %s [%s]", env.Environment.Name, err))
			continue
		}

		bar.ChangeMax(bar.GetMax() + len(targets))

		for _, target := range targets {
			i := target
			// newYaml := createYaml(i.Yaml, c.sourceOrg, c.sourceProject, c.targetOrg, c.targetProject)

			err := c.api.createTarget(&model.Target{
				Name:        i.Name,
				Identifier:  i.Identifier,
				Org:         c.targetOrg,
				Project:     c.targetProject,
				Environment: i.Environment,
				Attributes:  i.Attributes,
				Segments:    i.Segments,
			})
			if err != nil {
				failed = append(failed, fmt.Sprintln(e.Name, "/", i.Name, "-", err.Error()))
			}
			bar.Add(1)
		}
		bar.Add(1)
	}
	bar.Finish()

	reportFailed(failed, "targets:")
	return nil
}

func (api *ApiRequest) listTargets(org, project, envId string) ([]*model.Target, error) {

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
		Get(BaseURL + TARGETS)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse(resp)
	}

	result := model.TargetListResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	targets := []*model.Target{}
	for _, c := range result.Targets {
		target := c
		targets = append(targets, &target)
	}

	return targets, nil
}

func (api *ApiRequest) createTarget(target *model.Target) error {

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(target).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     target.Org,
		}).
		Post(BaseURL + TARGETS)
	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}

	return nil
}
