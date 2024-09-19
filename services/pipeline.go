package services

import (
	"encoding/json"
	"fmt"

	"github.com/jf781/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
	"go.uber.org/zap"
)

const LIST_PIPELINES = "/pipeline/api/pipelines/list"
const GET_PIPELINE = "/pipeline/api/pipelines/%s"
const CREATE_PIPELINE = "/pipeline/api/pipelines/v2"

type PipelineContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
	logger        *zap.Logger
}

func NewPipelineOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string, logger *zap.Logger) PipelineContext {
	return PipelineContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
		logger:        logger,
	}
}

func (c PipelineContext) Copy() error {

	c.logger.Info("Copying pipelines",
		zap.String("project", c.sourceProject),
	)

	pipelines, err := c.api.listPipelines(c.sourceOrg, c.sourceProject, c.logger)
	if err != nil {
		c.logger.Error("Failed to retrive pipelines",
			zap.String("Project", c.sourceProject),
			zap.Error(err),
		)
		return err
	}

	bar := progressbar.Default(int64(len(pipelines)), "Pipelines   ")

	for _, pipe := range pipelines {

		IncrementPipelinesTotal()

		pipeData, err := c.api.getPipeline(c.sourceOrg, c.sourceProject, pipe.Identifier, c.logger)
		if err == nil {
			newYaml := createYaml(pipeData.YAMLPipeline, c.sourceOrg, c.sourceProject, c.targetOrg, c.targetProject)
			err = c.api.createPipeline(c.targetOrg, c.targetProject, newYaml, c.logger)
		}
		if err != nil {
			c.logger.Error("Failed to create pipeline",
				zap.String("pipeline", pipe.Name),
				zap.Error(err),
			)
		} else {
			IncrementConnectorsMoved()
		}
		bar.Add(1)
	}
	bar.Finish()

	return nil
}

func (api *ApiRequest) listPipelines(org, project string, logger *zap.Logger) ([]*model.PipelineListContent, error) {

	logger.Info("Fetching pipelines",
		zap.String("org", org),
		zap.String("project", project),
	)

	IncrementApiCalls()

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(`{"filterType": "PipelineSetup"}`).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     org,
			"projectIdentifier": project,
			"size":              "1000",
		}).
		Post(api.BaseURL + LIST_PIPELINES)
	if err != nil {
		logger.Error("Failed to request to list of pipelines",
			zap.Error(err),
		)
		return nil, err
	}
	if resp.IsError() {
		logger.Error("Error response from API when listing pipelines",
			zap.String("response",
				resp.String(),
			),
		)
		return nil, handleErrorResponse(resp)
	}

	result := model.PipelineListResult{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		logger.Error("Failed to parse response from API",
			zap.Error(err),
		)
		return nil, err
	}

	return result.Data.Content, nil
}

func (api *ApiRequest) getPipeline(org, project, pipeIdentifier string, logger *zap.Logger) (*model.PipelineGetData, error) {
	logger.Info("Fetching pipeline",
		zap.String("org", org),
		zap.String("project", project),
		zap.String("pipeline", pipeIdentifier),
	)

	IncrementApiCalls()

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Load-From-Cache", "false").
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     org,
			"projectIdentifier": project,
		}).
		Get(api.BaseURL + fmt.Sprintf(GET_PIPELINE, pipeIdentifier))
	if err != nil {
		logger.Error("Failed to request details of pipeline",
			zap.String("pipeline", pipeIdentifier),
			zap.Error(err),
		)
		return nil, err
	}
	if resp.IsError() {
		logger.Error("Error response from API when listing pipelines",
			zap.String("response",
				resp.String(),
			),
		)
		return nil, handleErrorResponse(resp)
	}

	result := model.PipelineGetResult{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		logger.Error("Failed to parse response from API",
			zap.Error(err),
		)
		return nil, err
	}

	return result.Data, nil
}

func (api *ApiRequest) createPipeline(org, project, yaml string, logger *zap.Logger) error {

	logger.Info("Creating pipeline",
		zap.String("project", project),
	)

	IncrementApiCalls()

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/yaml").
		SetBody(yaml).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     org,
			"projectIdentifier": project,
		}).
		Post(api.BaseURL + CREATE_PIPELINE)
	if err != nil {
		logger.Error("Failed to send request to create pipeline",
			zap.Error(err),
		)
		return err
	}
	if resp.IsError() {
		var errorResponse map[string]interface{}
		if err := json.Unmarshal(resp.Body(), &errorResponse); err == nil {
			if code, ok := errorResponse["code"].(string); ok && code == "DUPLICATE_FIELD" {
				// Log as a warning and skip the error
				logger.Info("Duplicate pipeline found, ignoring error")
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
