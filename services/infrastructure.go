package services

import (
	"encoding/json"

	"github.com/jf781/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
	"go.uber.org/zap"
)

const INFRASTRUCTURE = "/ng/api/infrastructures"

type InfrastructureContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
	logger        *zap.Logger
}

func NewInfrastructureOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string, logger *zap.Logger) InfrastructureContext {
	return InfrastructureContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
		logger:        logger,
	}
}

func (c InfrastructureContext) Copy() error {

	c.logger.Info("Copying infrastructure",
		zap.String("project", c.sourceProject),
	)

	envs, err := c.api.listEnvironments(c.sourceOrg, c.sourceProject, c.logger)
	if err != nil {
		c.logger.Error("Failed to retrive environments",
			zap.String("Project", c.sourceProject),
			zap.Error(err),
		)
		return err
	}

	bar := progressbar.Default(int64(len(envs)), "Infrastructure")

	for _, env := range envs {
		e := env.Environment
		infras, err := c.api.listInfraDef(c.sourceOrg, c.sourceProject, e.Identifier, c.logger)
		if err != nil {
			c.logger.Error("Failed to retrive infrastructure",
				zap.String("Project", c.sourceProject),
				zap.Error(err),
			)
			continue
		}

		bar.ChangeMax(bar.GetMax() + len(infras))

		for _, infra := range infras {
			i := infra.Infrastructure

			c.logger.Info("Processing infrastructure",
				zap.String("infrastructure", i.Name),
				zap.String("targetProject", c.targetProject),
			)
			newYaml := createYaml(i.Yaml, c.sourceOrg, c.sourceProject, c.targetOrg, c.targetProject)

			err := c.api.createInfrastructure(&model.CreateInfrastructureRequest{
				Name:              i.Name,
				Identifier:        i.Identifier,
				OrgIdentifier:     c.targetOrg,
				ProjectIdentifier: c.targetProject,
				Description:       i.Description,
				EnvironmentRef:    e.Identifier,
				DeploymentType:    i.DeploymentType,
				Type:              i.Type,
				Yaml:              newYaml,
			}, c.logger)
			if err != nil {
				c.logger.Error("Failed to create infrastructure",
					zap.String("infrastructure", i.Name),
					zap.Error(err),
				)
			}
			bar.Add(1)
		}
		bar.Add(1)
	}
	bar.Finish()

	return nil
}

func (api *ApiRequest) listInfraDef(org, project, envId string, logger *zap.Logger) ([]*model.InfraDefListContent, error) {

	logger.Info("Fetching infrastructure",
		zap.String("org", org),
		zap.String("project", project),
		zap.String("environment", envId),
	)

	IncrementApiCalls()

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
		Get(api.BaseURL + INFRASTRUCTURE)
	if err != nil {
		logger.Error("Failed to request to list of infrastructure",
			zap.Error(err),
		)
		return nil, err
	}
	if resp.IsError() {
		logger.Error("Error response from API when listing infrastructure",
			zap.String("response",
				resp.String(),
			),
		)
		return nil, handleErrorResponse(resp)
	}

	result := model.InfraDefListResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		logger.Error("Failed to parse response from API",
			zap.Error(err),
		)
		return nil, err
	}

	return result.Data.Content, nil
}

func (api *ApiRequest) createInfrastructure(infra *model.CreateInfrastructureRequest, logger *zap.Logger) error {
	logger.Info("Creating infrastructure",
		zap.String("infrastructure", infra.Name),
		zap.String("project", infra.ProjectIdentifier),
	)

	IncrementApiCalls()

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(infra).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
		}).
		Post(api.BaseURL + INFRASTRUCTURE)
	if err != nil {
		logger.Error("Failed to send request to create ",
			zap.String("Infrastructure", infra.Name),
			zap.Error(err),
		)
		return err
	}
	if resp.IsError() {
		var errorResponse map[string]interface{}
		if err := json.Unmarshal(resp.Body(), &errorResponse); err == nil {
			if code, ok := errorResponse["code"].(string); ok && code == "DUPLICATE_FIELD" {
				// Log as a warning and skip the error
				logger.Info("Duplicate infrastructure found, ignoring error",
					zap.String("infrastructure", infra.Name),
				)
				return nil
			}
		} else {
			logger.Error(
				"Error response from API when creating ",
				zap.String("infrastructure", infra.Name),
				zap.String("response",
					resp.String(),
				),
			)
		}
		return handleErrorResponse(resp)
	}

	return nil
}
