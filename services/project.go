package services

import (
	"encoding/json"
	"fmt"

	"github.com/jf781/harness-move-project/model"
)

const GET_PROJECT = "/ng/api/projects/{identifier}"

// Validate if the project exists
func (api *ApiRequest) ValidateProject(org, project string) error {

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetPathParam("identifier", project).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     org,
		}).
		Get(BaseURL + GET_PROJECT)
	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}
	result := model.GetProjectResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return err
	}

	if result.Data == nil {
		return fmt.Errorf("org %s or project %s not exist", org, project)
	}
	return nil
}

// Create new project if it does not exist

const NEW_PROJECT = "/ng/api/projects"

type ProjectContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
}

func NewProjectOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string) ProjectContext {
	return ProjectContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
	}
}

func (c ProjectContext) Move() error {

	sourceProject, err := c.api.getProject(c.sourceOrg, c.sourceProject)
	if err != nil {
		return err
	}
	var failed []string

	newProject := &model.Project{
		Identifier:    sourceProject.Identifier,
		Name:          sourceProject.Name,
		Color:         sourceProject.Color,
		Modules:       sourceProject.Modules,
		Description:   sourceProject.Description,
		Tags:          sourceProject.Tags,
		OrgIdentifier: c.targetOrg,
	}

	err = c.api.CreateProject(newProject)

	if err != nil {
		failed = append(failed, fmt.Sprintln(sourceProject.Identifier, "-", err.Error()))
	}

	reportFailed(failed, "Create Project:")
	return nil
}

func (api *ApiRequest) getProject(org, project string) (model.Project, error) {

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetPathParam("identifier", project).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     org,
		}).
		Get(BaseURL + GET_PROJECT)

	if err != nil {
		return model.Project{}, err
	}

	if resp.IsError() {
		return model.Project{}, handleErrorResponse(resp)
	}
	result := model.GetProjectResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return model.Project{}, err
	}

	existingProject := model.Project{}
	existingProject = *result.Data.Project

	return existingProject, nil
}

func (api *ApiRequest) CreateProject(project *model.Project) error {

	wrappedProject := model.ProjectWrapper{
		Project: project,
	}

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(wrappedProject).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     project.OrgIdentifier,
		}).
		Post(BaseURL + NEW_PROJECT)
	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}
	result := model.GetProjectResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return err
	}

	if result.Data == nil {
		return fmt.Errorf(" project %s not exist", project)
	}
	return nil
}
