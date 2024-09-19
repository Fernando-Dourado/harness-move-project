package services

import (
	"encoding/json"

	"github.com/jf781/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
	"go.uber.org/zap"
)

const LIST_SERVICES = "/ng/api/servicesV2"
const CREATE_SERVICES = "/ng/api/servicesV2"

type ServiceContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
	logger        *zap.Logger
}

func NewServiceOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string, logger *zap.Logger) ServiceContext {
	return ServiceContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
		logger:        logger,
	}
}

func (c ServiceContext) Copy() error {

	c.logger.Info("Copying Services",
		zap.String("project", c.sourceProject),
	)

	services, err := c.listServices(c.sourceOrg, c.sourceProject, c.logger)
	if err != nil {
		c.logger.Error("Failed to retrive serivces",
			zap.String("Project", c.sourceProject),
			zap.Error(err),
		)
		return err
	}

	bar := progressbar.Default(int64(len(services)), "Services    ")

	for _, s := range services {

		IncrementServicesTotal()

		c.logger.Info("Processing service",
			zap.String("service", s.Service.Name),
			zap.String("targetProject", c.targetProject),
		)
		newYaml := createYaml(s.Service.Yaml, c.sourceOrg, c.sourceProject, c.targetOrg, c.targetProject)
		service := &model.CreateServiceRequest{
			OrgIdentifier:     c.targetOrg,
			ProjectIdentifier: c.targetProject,
			Identifier:        s.Service.Identifier,
			Name:              s.Service.Name,
			Description:       s.Service.Description,
			Yaml:              newYaml,
		}
		if err := c.createService(service, c.logger); err != nil {
			c.logger.Error("Failed to create service",
				zap.String("service", s.Service.Name),
				zap.Error(err),
			)
		} else {
			IncrementServicesMoved()
		}
		bar.Add(1)
	}
	bar.Finish()

	return nil
}

func (c ServiceContext) listServices(org, project string, logger *zap.Logger) ([]*model.ServiceListContent, error) {

	logger.Info("Fetching service overrides",
		zap.String("org", org),
		zap.String("project", project),
	)

	IncrementApiCalls()

	api := c.api
	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     org,
			"projectIdentifier": project,
			"size":              "1000",
		}).
		Get(api.BaseURL + LIST_SERVICES)
	if err != nil {
		logger.Error("Failed to request to list of services",
			zap.Error(err),
		)
		return nil, err
	}
	if resp.IsError() {
		logger.Error("Error response from API when listing services",
			zap.String("response",
				resp.String(),
			),
		)
		return nil, handleErrorResponse(resp)
	}

	result := model.ServiceListResult{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		logger.Error("Failed to parse response from API",
			zap.Error(err),
		)
		return nil, err
	}

	return result.Data.Content, nil
}

func (c ServiceContext) createService(service *model.CreateServiceRequest, logger *zap.Logger) error {

	logger.Info("Creating service",
		zap.String("service", service.Name),
		zap.String("project", service.ProjectIdentifier),
	)

	IncrementApiCalls()

	api := c.api
	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(service).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
		}).
		Post(api.BaseURL + CREATE_SERVICES)
	if err != nil {
		logger.Error("Failed to send request to create ",
			zap.String("Service", service.Name),
			zap.Error(err),
		)
		return err
	}
	if resp.IsError() {
		var errorResponse map[string]interface{}
		if err := json.Unmarshal(resp.Body(), &errorResponse); err == nil {
			if code, ok := errorResponse["code"].(string); ok && code == "DUPLICATE_FIELD" {
				// Log as a warning and skip the error
				logger.Info("Duplicate service found, ignoring error",
					zap.String("service", service.Name),
				)
				return nil
			}
		} else {
			logger.Error(
				"Error response from API when creating ",
				zap.String("Service", service.Name),
				zap.String("response",
					resp.String(),
				),
			)
		}
		return handleErrorResponse(resp)
	}

	return nil
}
