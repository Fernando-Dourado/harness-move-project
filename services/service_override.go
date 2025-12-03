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
	exportDir     string
}

func NewServiceOverrideOperation(sourceApi *SourceRequest, targetApi *TargetRequest, st *SourceTarget) ServiceOverrideContext {
	return ServiceOverrideContext{
		source:        sourceApi,
		target:        targetApi,
		sourceOrg:     st.SourceOrg,
		sourceProject: st.SourceProject,
		targetOrg:     st.TargetOrg,
		targetProject: st.TargetProject,
		exportDir:     st.ExportDir,
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
				newYaml := createYaml(o.YAML, c.sourceOrg, c.sourceProject, c.targetOrg, c.targetProject, c.source.Account, c.target.Account)
				req := &model.CreateServiceOverrideRequest{
					OrgIdentifier:     c.targetOrg,
					ProjectIdentifier: c.targetProject,
					EnvironmentRef:    o.EnvironmentRef,
					ServiceRef:        o.ServiceRef,
					YAML:              newYaml,
				}

				var err error
				if isExportMode(c.exportDir) {
					// Save YAML string directly to file
					err = saveStringToYAML(c.exportDir, "service-overrides-v1", o.EnvironmentRef+"_"+o.ServiceRef, newYaml)
				} else {
					// Make API call to create service override
					err = createServiceOverride(c.target, req)
				}

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

	params := createQueryParams(s.Account, org, project)
	params["size"] = "1000"

	resp, err := s.Client.R().
		SetHeader("x-api-key", s.Token).
		SetQueryParams(params).
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
