package services

import (
	"encoding/json"
	"fmt"

	"github.com/Fernando-Dourado/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
)

const LIST_PIPELINES = "/pipeline/api/pipelines/list"
const GET_PIPELINE = "/pipeline/api/pipelines/%s"
const CREATE_PIPELINE = "/pipeline/api/pipelines/v2"

type PipelineContext struct {
	source        *SourceRequest
	target        *TargetRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
}

func NewPipelineOperation(sourceApi *SourceRequest, targetApi *TargetRequest, sourceOrg, sourceProject, targetOrg, targetProject string) PipelineContext {
	return PipelineContext{
		source:        sourceApi,
		target:        targetApi,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
	}
}

func (c PipelineContext) Move() error {

	pipelines, err := c.source.listPipelines(c.sourceOrg, c.sourceProject)
	if err != nil {
		return err
	}

	bar := progressbar.Default(int64(len(pipelines)), "Pipelines   ")
	var failed []string

	for _, pipe := range pipelines {
		pipeData, err := c.getPipeline(c.sourceOrg, c.sourceProject, pipe.Identifier)
		if err == nil {
			newYaml := createYaml(pipeData.YAMLPipeline, c.sourceOrg, c.sourceProject, c.targetOrg, c.targetProject)
			err = c.createPipeline(c.targetOrg, c.targetProject, newYaml)
		}
		if err != nil {
			failed = append(failed, fmt.Sprintln(pipe.Name, "-", err.Error()))
		}
		bar.Add(1)
	}
	bar.Finish()

	reportFailed(failed, "pipelines:")
	return nil
}

func (s *SourceRequest) listPipelines(org, project string) ([]*model.PipelineListContent, error) {

	resp, err := s.Client.R().
		SetHeader("x-api-key", s.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(`{"filterType": "PipelineSetup"}`).
		SetQueryParams(map[string]string{
			"accountIdentifier": s.Account,
			"orgIdentifier":     org,
			"projectIdentifier": project,
			"size":              "1000",
		}).
		Post(s.Url + LIST_PIPELINES)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse(resp)
	}

	result := model.PipelineListResult{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	return result.Data.Content, nil
}

func (c PipelineContext) getPipeline(org, project, pipeIdentifier string) (*model.PipelineGetData, error) {

	api := c.source
	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Load-From-Cache", "false").
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     org,
			"projectIdentifier": project,
		}).
		Get(api.Url + fmt.Sprintf(GET_PIPELINE, pipeIdentifier))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse(resp)
	}

	result := model.PipelineGetResult{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (c PipelineContext) createPipeline(org, project, yaml string) error {

	api := c.target
	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/yaml").
		SetBody(yaml).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     org,
			"projectIdentifier": project,
		}).
		Post(api.Url + CREATE_PIPELINE)
	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}

	return nil
}
