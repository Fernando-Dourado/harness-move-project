package services

import (
	"encoding/json"
	"fmt"

	"github.com/Fernando-Dourado/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
)

type InfrastructureContext struct {
	source        *SourceRequest
	target        *TargetRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
}

func NewInfrastructureOperation(sourceApi *SourceRequest, targetApi *TargetRequest, sourceOrg, sourceProject, targetOrg, targetProject string) InfrastructureContext {
	return InfrastructureContext{
		source:        sourceApi,
		target:        targetApi,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
	}
}

func (c InfrastructureContext) Move() error {

	envs, err := c.source.listEnvironments(c.sourceOrg, c.sourceProject)
	if err != nil {
		return err
	}

	bar := progressbar.Default(int64(len(envs)), "Infrastructure")
	var failed []string

	for _, env := range envs {
		e := env.Environment
		infras, err := listInfraDef(c.source, c.sourceOrg, c.sourceProject, e.Identifier)
		if err != nil {
			failed = append(failed, fmt.Sprintf("Unable to list infrastructures for environment %s [%s]", env.Environment.Name, err))
			continue
		}

		bar.ChangeMax(bar.GetMax() + len(infras))

		for _, infra := range infras {
			i := infra.Infrastructure
			newYaml := createYaml(i.Yaml, c.sourceOrg, c.sourceProject, c.targetOrg, c.targetProject)

			err := createInfrastructure(c.target, &model.CreateInfrastructureRequest{
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

func listInfraDef(s *SourceRequest, org, project, envId string) ([]*model.InfraDefListContent, error) {

	resp, err := s.Client.R().
		SetHeader("x-api-key", s.Token).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"accountIdentifier":     s.Account,
			"orgIdentifier":         org,
			"projectIdentifier":     project,
			"environmentIdentifier": envId,
			"size":                  "1000",
		}).
		Get(s.Url + "/ng/api/infrastructures")
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

func createInfrastructure(t *TargetRequest, infra *model.CreateInfrastructureRequest) error {

	resp, err := t.Client.R().
		SetHeader("x-api-key", t.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(infra).
		SetQueryParams(map[string]string{
			"accountIdentifier": t.Account,
		}).
		Post(t.Url + "/ng/api/infrastructures")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}

	return nil
}
