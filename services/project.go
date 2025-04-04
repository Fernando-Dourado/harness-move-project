package services

import (
	"encoding/json"
	"fmt"

	"github.com/Fernando-Dourado/harness-move-project/model"
	"github.com/go-resty/resty/v2"
)

const (
	GET_PROJECT    = "/ng/api/projects/{identifier}"
	CREATE_PROJECT = "/v1/orgs/{org}/projects"
)

func (s *SourceRequest) ValidateSource(org, project string) error {
	return validateOrgProject(s.Client, s.Token, s.Account, org, project)
}

func (t *TargetRequest) ValidateTarget(org, project string) error {
	return validateOrgProject(t.Client, t.Token, t.Account, org, project)
}

func validateOrgProject(c *resty.Client, token, account, org, project string) error {
	result, err := getProject(c, token, account, org, project)
	if err != nil {
		return err
	}
	if result.Data == nil {
		return fmt.Errorf("org %s or project %s not exist", org, project)
	}
	return nil
}

func getProject(c *resty.Client, token, account, org, project string) (*model.GetProjectResponse, error) {
	resp, err := c.R().
		SetHeader("x-api-key", token).
		SetHeader("Content-Type", "application/json").
		SetPathParam("identifier", project).
		SetQueryParams(map[string]string{
			"accountIdentifier": account,
			"orgIdentifier":     org,
		}).
		Get(BaseURL + GET_PROJECT)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse(resp)
	}
	result := model.GetProjectResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

type ProjectContext struct {
	source        *SourceRequest
	target        *TargetRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
}

func NewProjectOperation(sourceApi *SourceRequest, targetApi *TargetRequest, st *SourceTarget) ProjectContext {
	return ProjectContext{
		source:        sourceApi,
		target:        targetApi,
		sourceOrg:     st.SourceOrg,
		sourceProject: st.SourceProject,
		targetOrg:     st.TargetOrg,
		targetProject: st.TargetProject,
	}
}

func (c ProjectContext) Move() error {
	response, err := getProject(c.source.Client, c.source.Token, c.source.Account, c.sourceOrg, c.sourceProject)
	if err != nil {
		return err
	}
	if response.Data == nil || response.Data.Project == nil {
		return fmt.Errorf("invalid response data for project %s and org %s", c.sourceProject, c.sourceOrg)
	}

	newProject := model.Project{
		OrgIdentifier: c.targetOrg,
		Identifier:    c.targetProject,
		Name:          response.Data.Project.Name,
		Description:   response.Data.Project.Description,
		Color:         response.Data.Project.Color,
	}

	err = c.target.createProject(newProject)
	if err != nil {
		return err
	}

	return nil
}

func (t *TargetRequest) createProject(project model.Project) error {
	request := model.CreateProjectRequest{
		Project: &project,
	}
	resp, err := t.Client.R().
		SetHeader("x-api-key", t.Token).
		SetHeader("Content-Type", "application/json").
		SetHeader("Harness-Account", t.Account).
		SetPathParam("org", project.OrgIdentifier).
		SetBody(request).
		Post(BaseURL + CREATE_PROJECT)
	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}
	return nil
}
