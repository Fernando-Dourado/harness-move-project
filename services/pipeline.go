package services

import (
	"encoding/json"
	"fmt"

	"github.com/Fernando-Dourado/harness-move-project/model"
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
}

func NewPipelineOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string) PipelineContext {
	return PipelineContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
	}
}

func (c PipelineContext) Move() error {

	pipelines, err := c.api.listPipelines(c.sourceOrg, c.sourceProject)
	if err != nil {
		return err
	}

	for _, pipe := range pipelines {
		pipeData, err := c.api.getPipeline(c.sourceOrg, c.sourceProject, pipe.Identifier)
		if err != nil {
			return err
		}

		if !pipeData.EntityValidityDetails.Valid {
			// log warn as pipeline YAML is invalid
			continue
		}

		newYaml := createYaml(pipeData.YAMLPipeline, c.sourceOrg, c.sourceProject, c.targetOrg, c.targetProject)
		if err := c.api.createPipeline(c.targetOrg, c.targetProject, newYaml); err != nil {
			return err
		}
	}

	return nil
}

func (api *ApiRequest) listPipelines(org, project string) ([]*model.PipelineListContent, error) {

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
		Post(BaseURL + LIST_PIPELINES)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf(resp.Status())
	}

	result := model.PipelineListResult{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	return result.Data.Content, nil
}

func (api *ApiRequest) getPipeline(org, project, pipeIdentifier string) (*model.PipelineGetData, error) {

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Load-From-Cache", "false").
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     org,
			"projectIdentifier": project,
		}).
		Get(BaseURL + fmt.Sprintf(GET_PIPELINE, pipeIdentifier))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf(resp.Status())
	}

	result := model.PipelineGetResult{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (api *ApiRequest) createPipeline(org, project, yaml string) error {

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/yaml").
		SetBody(yaml).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     org,
			"projectIdentifier": project,
		}).
		Post(BaseURL + CREATE_PIPELINE)
	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleCreateErrorResponse(resp)
	}

	return nil
}
