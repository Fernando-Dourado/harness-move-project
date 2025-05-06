package services

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Fernando-Dourado/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
)

type EnvironmentContext struct {
	source        *SourceRequest
	target        *TargetRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
}

func NewEnvironmentOperation(sourceApi *SourceRequest, targetApi *TargetRequest, sourceOrg, sourceProject, targetOrg, targetProject string) EnvironmentContext {
	return EnvironmentContext{
		source:        sourceApi,
		target:        targetApi,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
	}
}

func (c EnvironmentContext) Move() error {

	envs, err := c.source.listEnvironments(c.sourceOrg, c.sourceProject)
	if err != nil {
		return nil
	}

	bar := progressbar.Default(int64(len(envs)), "Environments")
	var failed []string

	for _, env := range envs {
		e := env.Environment

		newYaml := createYaml(sanitizeEnvYaml(e.Yaml), c.sourceOrg, c.sourceProject, c.targetOrg, c.targetProject)
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
		if err := createEnvironment(c.target, req); err != nil {
			failed = append(failed, fmt.Sprintln(e.Name, "-", err.Error()))
		}
		bar.Add(1)
	}
	bar.Finish()

	reportFailed(failed, "environments:")
	return nil
}

func (s *SourceRequest) listEnvironments(org, project string) ([]*model.ListEnvironmentContent, error) {

	resp, err := s.Client.R().
		SetHeader("x-api-key", s.Token).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"accountIdentifier": s.Account,
			"orgIdentifier":     org,
			"projectIdentifier": project,
			"size":              "1000",
		}).
		Get(s.Url + "/ng/api/environmentsV2")
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

func createEnvironment(t *TargetRequest, env *model.CreateEnvironmentRequest) error {

	resp, err := t.Client.R().
		SetHeader("x-api-key", t.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(env).
		SetQueryParams(map[string]string{
			"accountIdentifier": t.Account,
		}).
		Post(t.Url + "/ng/api/environmentsV2")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}

	return nil
}

func sanitizeEnvYaml(yaml string) string {
	return strings.ReplaceAll(yaml, "\"", "")
}
