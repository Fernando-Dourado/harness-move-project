package services

import (
	"encoding/json"
	"fmt"

	"github.com/Fernando-Dourado/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
)

type ServiceOverrideContext struct {
	source        *SourceRequest
	target        *TargetRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
}

func NewServiceOverrideOperation(sourceApi *SourceRequest, targetApi *TargetRequest, st *SourceTarget) ServiceOverrideContext {
	return ServiceOverrideContext{
		source:        sourceApi,
		target:        targetApi,
		sourceOrg:     st.SourceOrg,
		sourceProject: st.SourceProject,
		targetOrg:     st.TargetOrg,
		targetProject: st.TargetProject,
	}
}

func (c ServiceOverrideContext) Move() error {

	envs, err := c.source.listEnvironments(c.sourceOrg, c.sourceProject)
	if err != nil {
		return err
	}

	bar := progressbar.Default(int64(len(envs)), "Overrides V1")
	var failed []string

	for _, env := range envs {
		e := env.Environment
		overrides, err := listServiceOverrides(c.source, c.sourceOrg, c.sourceProject, e.Identifier)
		if err != nil {
			failed = append(failed, fmt.Sprintf("Unable to list service overrides for environment %s [%s]", env.Environment.Name, err))
			continue
		}

		bar.ChangeMax(bar.GetMax() + len(overrides))

		for _, o := range overrides {
			if len(o.YAML) == 0 {
				failed = append(failed, fmt.Sprintf("The YAML is empty [envId=%s,serviceRef=%s]", o.EnvironmentRef, o.ServiceRef))
			} else {
				err := createServiceOverride(c.target, &model.CreateServiceOverrideRequest{
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

	reportFailed(failed, "overrides v1:")
	return nil
}

func listServiceOverrides(s *SourceRequest, org, project, envId string) ([]*model.ServiceOverride, error) {

	resp, err := s.Client.R().
		SetHeader("x-api-key", s.Token).
		SetQueryParams(map[string]string{
			"accountIdentifier":     s.Account,
			"orgIdentifier":         org,
			"projectIdentifier":     project,
			"environmentIdentifier": envId,
			"size":                  "1000",
		}).
		Get(s.Url + "/ng/api/environmentsV2/serviceOverrides")
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

func createServiceOverride(t *TargetRequest, override *model.CreateServiceOverrideRequest) error {

	resp, err := t.Client.R().
		SetHeader("x-api-key", t.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(override).
		SetQueryParams(map[string]string{
			"accountIdentifier": t.Account,
		}).
		Post(t.Url + "/ng/api/environmentsV2/serviceOverrides")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}

	return nil
}
