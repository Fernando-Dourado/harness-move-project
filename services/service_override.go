package services

import (
	"encoding/json"

	"github.com/jf781/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
	"go.uber.org/zap"
)

type ServiceOverrideContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
	logger        *zap.Logger
}

func NewServiceOverrideOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string, logger *zap.Logger) ServiceOverrideContext {
	return ServiceOverrideContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
		logger:        logger,
	}
}

func (c ServiceOverrideContext) Copy() error {

	envs, err := c.api.listEnvironments(c.sourceOrg, c.sourceProject, c.logger)
	if err != nil {
		c.logger.Error("Failed to retrive service overrides",
			zap.String("Project", c.sourceProject),
			zap.Error(err),
		)
		return err
	}

	bar := progressbar.Default(int64(len(envs)), "Service Override")

	for _, env := range envs {
		e := env.Environment
		overrides, err := c.api.listServiceOverrides(c.sourceOrg, c.sourceProject, e.Identifier, c.logger)
		if err != nil {
			c.logger.Error("Failed to retrive environments",
				zap.String("Project", c.sourceProject),
				zap.Error(err),
			)
			continue
		}

		bar.ChangeMax(bar.GetMax() + len(overrides))

		for _, o := range overrides {

			c.logger.Info("Processing service override",
				zap.String("service override", o.ServiceRef),
				zap.String("targetProject", c.targetProject),
			)

			if len(o.YAML) == 0 {
				c.logger.Error("YAML file is empty",
					zap.String("Project", c.sourceProject),
					zap.Error(err),
				)
			} else {
				err := c.api.createServiceOverride(&model.CreateServiceOverrideRequest{
					OrgIdentifier:     c.targetOrg,
					ProjectIdentifier: c.targetProject,
					EnvironmentRef:    o.EnvironmentRef,
					ServiceRef:        o.ServiceRef,
					YAML:              o.YAML,
				}, c.logger)
				if err != nil {
					c.logger.Error("Failed to create service override",
						zap.String("service override", o.ServiceRef),
						zap.Error(err),
					)
				}
			}
			bar.Add(1)
		}
		bar.Add(1)
	}
	bar.Finish()

	return nil
}

func (api *ApiRequest) listServiceOverrides(org, project, envId string, logger *zap.Logger) ([]*model.ServiceOverride, error) {

	logger.Info("Fetching service overrides",
		zap.String("org", org),
		zap.String("project", project),
	)

	IncrementApiCalls()

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetQueryParams(map[string]string{
			"accountIdentifier":     api.Account,
			"orgIdentifier":         org,
			"projectIdentifier":     project,
			"environmentIdentifier": envId,
			"size":                  "1000",
		}).
		Get(api.BaseURL + "/ng/api/environmentsV2/serviceOverrides")
	if err != nil {
		logger.Error("Failed to request to list of service overrides",
			zap.Error(err),
		)
		return nil, err
	}
	if resp.IsError() {
		logger.Error("Error response from API when listing service overrides",
			zap.String("response",
				resp.String(),
			),
		)
		return nil, handleErrorResponse(resp)
	}

	result := model.ListServiceOverridesRequest{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		logger.Error("Failed to parse response from API",
			zap.Error(err),
		)
		return nil, err
	}

	return result.Data.Content, nil
}

func (api *ApiRequest) createServiceOverride(override *model.CreateServiceOverrideRequest, logger *zap.Logger) error {

	logger.Info("Creating service override",
		zap.String("project", override.ProjectIdentifier),
	)

	IncrementApiCalls()

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(override).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
		}).
		Post(api.BaseURL + "/ng/api/environmentsV2/serviceOverrides")
	if err != nil {
		logger.Error("Failed to send request to create ")
		return err
	}
	if resp.IsError() {
		var errorResponse map[string]interface{}
		if err := json.Unmarshal(resp.Body(), &errorResponse); err == nil {
			if code, ok := errorResponse["code"].(string); ok && code == "DUPLICATE_FIELD" {
				// Log as a warning and skip the error
				logger.Info("Duplicate connector found, ignoring error")
				return nil
			}
		} else {
			logger.Error(
				"Error response from API when creating ",
				zap.String("response",
					resp.String(),
				),
			)
		}
		return handleErrorResponse(resp)
	}

	return nil
}
