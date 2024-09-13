package services

import (
	"encoding/json"

	"github.com/jf781/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
	"go.uber.org/zap"
)

type EnvironmentContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
	logger        *zap.Logger
}

func NewEnvironmentOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string, logger *zap.Logger) EnvironmentContext {
	return EnvironmentContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
		logger:        logger,
	}
}

func (c EnvironmentContext) Copy() error {

	c.logger.Info("Copying environments",
		zap.String("project", c.sourceProject),
	)

	envs, err := c.api.listEnvironments(c.sourceOrg, c.sourceProject, c.logger)
	if err != nil {
		c.logger.Error("Failed to retrive environments for ",
			zap.String("Project", c.sourceProject),
			zap.Error(err),
		)
		return nil
	}

	bar := progressbar.Default(int64(len(envs)), "Environments")

	for _, env := range envs {
		e := env.Environment

		c.logger.Info("Processing environments",
			zap.String("environemnt", e.Name),
			zap.String("targetProject", c.targetProject),
		)

		newYaml := createYaml(e.Yaml, c.sourceOrg, c.sourceProject, c.targetOrg, c.targetProject)
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
		if err := c.api.createEnvironment(req, c.logger); err != nil {
			c.logger.Error("Failed to create environment. ",
				zap.String("environment name", e.Name),
				zap.Error(err),
			)
		}
		bar.Add(1)
	}
	bar.Finish()

	return nil
}

func (api *ApiRequest) listEnvironments(org, project string, logger *zap.Logger) ([]*model.ListEnvironmentContent, error) {

	logger.Info("Fetching environments",
		zap.String("org", org),
		zap.String("project", project),
	)

	IncrementApiCalls()

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     org,
			"projectIdentifier": project,
			"size":              "1000",
		}).
		Get(api.BaseURL + "/ng/api/environmentsV2")
	if err != nil {
		logger.Error("Failed to request to list of environments",
			zap.Error(err),
		)
		return nil, err
	}
	if resp.IsError() {
		logger.Error(
			"Error response from API when listing environments",
			zap.String("response",
				resp.String(),
			),
		)
		return nil, handleErrorResponse(resp)
	}

	result := model.ListEnvironmentResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		logger.Error("Failed to parse response from API",
			zap.Error(err),
		)
		return nil, err
	}

	return result.Data.Content, nil
}

func (api *ApiRequest) createEnvironment(env *model.CreateEnvironmentRequest, logger *zap.Logger) error {

	logger.Info("Creating environment",
		zap.String("org", env.OrgIdentifier),
		zap.String("project", env.ProjectIdentifier),
	)

	IncrementApiCalls()

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(env).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
		}).
		Post(api.BaseURL + "/ng/api/environmentsV2")
	if err != nil {
		logger.Error("Failed to send request to create ",
			zap.String("Environment", env.Name),
			zap.Error(err),
		)
		return err
	}
	if resp.IsError() {
		var errorResponse map[string]interface{}
		if err := json.Unmarshal(resp.Body(), &errorResponse); err == nil {
			if code, ok := errorResponse["code"].(string); ok && code == "DUPLICATE_FIELD" {
				// Log as a warning and skip the error
				logger.Info("Duplicate environment found, ignoring error",
					zap.String("connectorName", env.Name),
				)
				return nil
			}
		} else {
			logger.Error(
				"Error response from API when creating ",
				zap.String("Environment", env.Name),
				zap.String("response",
					resp.String(),
				),
			)
		}
		return handleErrorResponse(resp)
	}

	return nil
}

// Commented out 9/11 - Was causing issues for creating environments without a description defined.
// func sanitizeEnvYaml(yaml string) string {
// 	return strings.ReplaceAll(yaml, "\"", "")
// }
