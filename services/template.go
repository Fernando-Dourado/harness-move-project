package services

import (
	"encoding/json"

	"github.com/jf781/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
	"go.uber.org/zap"
)

const LIST_TEMPLATES_ENDPOINT = "/v1/orgs/{org}/projects/{project}/templates"
const GET_TEMPLATE_ENDPOINT = "/template/api/templates/{templateIdentifier}"
const CREATE_TEMPLATE_ENDPOINT = "/template/api/templates"

type TemplateContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
	logger        *zap.Logger
}

func NewTemplateOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string, logger *zap.Logger) TemplateContext {
	return TemplateContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
		logger:        logger,
	}
}

func (c TemplateContext) Copy() error {

	c.logger.Info("Copying template",
		zap.String("project", c.sourceProject),
	)

	templates, err := c.listTemplates(c.sourceOrg, c.sourceProject, c.logger)
	if err != nil {
		c.logger.Error("Failed to retrive templates",
			zap.String("Project", c.sourceProject),
			zap.Error(err),
		)
		return err
	}

	bar := progressbar.Default(int64(len(templates)), "Templates   ")
	var failed []string

	for _, template := range templates {

		c.logger.Info("Processing template",
			zap.String("template", template.Name),
			zap.String("targetProject", c.targetProject),
		)
		t, err := c.getTemplate(c.sourceOrg, c.sourceProject, template.Identifier, template.VersionLabel, c.logger)
		if err == nil {
			newYaml := createYaml(t.Yaml, c.sourceOrg, c.sourceProject, c.targetOrg, c.targetProject)
			err = c.createTemplate(c.targetOrg, c.targetProject, newYaml, c.logger)
		}
		if err != nil {
			c.logger.Error("Failed to create template",
				zap.String("template", template.Name),
				zap.Error(err),
			)
		}
		bar.Add(1)
	}
	bar.Finish()

	reportFailed(failed, "templates:")
	return nil
}

func (c TemplateContext) listTemplates(org, project string, logger *zap.Logger) (model.TemplateListResult, error) {

	logger.Info("Fetching templates",
		zap.String("org", org),
		zap.String("project", project),
	)

	IncrementApiCalls()

	api := c.api
	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetHeader("Harness-Account", api.Account).
		SetPathParam("org", org).
		SetPathParam("project", project).
		SetQueryParams(map[string]string{
			"limit": "1000",
		}).
		Get(api.BaseURL + LIST_TEMPLATES_ENDPOINT)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		logger.Error("Failed to request to list of templates",
			zap.Error(err),
		)
		return nil, handleErrorResponse(resp)
	}

	result := model.TemplateListResult{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		logger.Error("Failed to parse response from API",
			zap.Error(err),
		)
		return nil, err
	}

	return result, nil
}

func (c TemplateContext) getTemplate(org, project, templateIdentifier, versionLabel string, logger *zap.Logger) (*model.TemplateGetData, error) {

	logger.Info("Getting template",
		zap.String("template", templateIdentifier),
		zap.String("project", project),
	)

	IncrementApiCalls()

	api := c.api
	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetHeader("Load-From-Cache", "false").
		SetPathParam("templateIdentifier", templateIdentifier).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     org,
			"projectIdentifier": project,
			"versionLabel":      versionLabel,
		}).
		Get(api.BaseURL + GET_TEMPLATE_ENDPOINT)
	if err != nil {
		logger.Error("Failed to request to template",
			zap.String("template", templateIdentifier),
			zap.Error(err),
		)
		return nil, err
	}
	if resp.IsError() {
		logger.Error("Error response from API when listing templates",
			zap.String("response",
				resp.String(),
			),
		)
		return nil, handleErrorResponse(resp)
	}

	result := model.TemplateGetResult{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		logger.Error("Failed to parse response from API",
			zap.Error(err),
		)
		return nil, err
	}

	return result.Data, nil
}

func (c TemplateContext) createTemplate(org, project, yaml string, logger *zap.Logger) error {

	api := c.api
	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(yaml).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     org,
			"projectIdentifier": project,
		}).
		Post(api.BaseURL + CREATE_TEMPLATE_ENDPOINT)
	if err != nil {
		logger.Error("Failed to send request to create ",
			zap.Error(err),
		)
		return err
	}
	if resp.IsError() {
		var errorResponse map[string]interface{}
		if err := json.Unmarshal(resp.Body(), &errorResponse); err == nil {
			if code, ok := errorResponse["code"].(string); ok && code == "DUPLICATE_FIELD" {
				// Log as a warning and skip the error
				logger.Info("Duplicate template found, ignoring error")
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
