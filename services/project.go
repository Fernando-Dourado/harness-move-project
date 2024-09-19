package services

import (
	"encoding/json"
	"fmt"

	"github.com/jf781/harness-move-project/model"
	"go.uber.org/zap"
)

const GET_PROJECT = "/ng/api/projects/{identifier}"

// Validate if the project exists
func (api *ApiRequest) ValidateProject(org, project string, logger *zap.Logger) error {

	logger.Info("Validating if project exists",
		zap.String("project", project),
	)

	IncrementApiCalls()

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetPathParam("identifier", project).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     org,
		}).
		Get(api.BaseURL + GET_PROJECT)
	if err != nil {
		logger.Error("Failed to request project",
			zap.String("Project", project),
			zap.Error(err),
		)
		return err
	}
	if resp.IsError() {
		logger.Warn("Unable to find existing project in organization",
			zap.String("response",
				resp.String(),
			),
		)
		return handleErrorResponse(resp)
	}
	result := model.GetProjectResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		logger.Error("Failed to parse response from API",
			zap.Error(err),
		)
		return err
	}

	if result.Data == nil {
		logger.Warn("Project does not exist.  Will be creating it")
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
	logger        *zap.Logger
}

func NewProjectOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string, logger *zap.Logger) ProjectContext {
	return ProjectContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
		logger:        logger,
	}
}

func (c ProjectContext) Copy() error {
	c.logger.Info("Creating new project",
		zap.String("project", c.sourceProject),
	)

	sourceProject, err := c.api.getProject(c.sourceOrg, c.sourceProject, c.logger)
	if err != nil {
		c.logger.Error("Failed to retrive source project",
			zap.String("Project", c.sourceProject),
			zap.Error(err),
		)
		return err
	}

	newProject := &model.Project{
		Identifier:    c.targetProject,
		Name:          c.targetProject,
		Color:         sourceProject.Color,
		Modules:       sourceProject.Modules,
		Description:   sourceProject.Description,
		Tags:          sourceProject.Tags,
		OrgIdentifier: c.targetOrg,
	}

	err = c.api.CreateProject(newProject, c.logger)

	if err != nil {
		c.logger.Error("Failed to create target project ",
			zap.String("Project", c.targetProject),
			zap.Error(err),
		)
	}

	return nil
}

func (api *ApiRequest) getProject(org, project string, logger *zap.Logger) (model.Project, error) {

	logger.Info("Getting source project",
		zap.String("project", project),
	)

	IncrementApiCalls()

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetPathParam("identifier", project).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     org,
		}).
		Get(api.BaseURL + GET_PROJECT)

	if err != nil {
		logger.Error("Failed to request source project",
			zap.Error(err),
		)
		return model.Project{}, err
	}

	if resp.IsError() {
		logger.Error("Error response from API when listing source project",
			zap.String("response",
				resp.String(),
			),
		)
		return model.Project{}, handleErrorResponse(resp)
	}
	result := model.GetProjectResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		logger.Error("Failed to parse response from API",
			zap.Error(err),
		)
		return model.Project{}, err
	}

	existingProject := model.Project{}
	existingProject = *result.Data.Project

	return existingProject, nil
}

func (api *ApiRequest) CreateProject(project *model.Project, logger *zap.Logger) error {

	logger.Info("Creating target project",
		zap.String("project", project.Name),
	)

	IncrementApiCalls()

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
		Post(api.BaseURL + NEW_PROJECT)
	if err != nil {
		logger.Error("Failed to send request to create ",
			zap.String("project", project.Name),
			zap.Error(err),
		)
		return err
	}
	if resp.IsError() {
		var errorResponse map[string]interface{}
		if err := json.Unmarshal(resp.Body(), &errorResponse); err == nil {
			if code, ok := errorResponse["code"].(string); ok && code == "DUPLICATE_FIELD" {
				// Log as a warning and skip the error
				logger.Info("Duplicate project found, ignoring error",
					zap.String("project", project.Name),
				)
				return nil
			}
		} else {
			logger.Error(
				"Error response from API when creating ",
				zap.String("project", project.Name),
				zap.String("response",
					resp.String(),
				),
			)
		}
		return handleErrorResponse(resp)
	}
	result := model.GetProjectResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		logger.Error("Failed to parse response from API",
			zap.Error(err),
		)
		return err
	}

	if result.Data == nil {
		logger.Error("Failed to validate project exists",
			zap.Error(err),
		)
	}
	return nil
}
