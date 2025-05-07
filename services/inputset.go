package services

import (
	"encoding/json"
	"fmt"

	"github.com/Fernando-Dourado/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
)

type InputsetContext struct {
	source        *SourceRequest
	target        *TargetRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
}

func NewInputsetOperation(sourceApi *SourceRequest, targetApi *TargetRequest, st *SourceTarget) InputsetContext {
	return InputsetContext{
		source:        sourceApi,
		target:        targetApi,
		sourceOrg:     st.SourceOrg,
		sourceProject: st.SourceProject,
		targetOrg:     st.TargetOrg,
		targetProject: st.TargetProject,
	}
}

func (c InputsetContext) Move() error {

	pipelines, err := c.source.listPipelines(c.sourceOrg, c.sourceProject)
	if err != nil {
		return err
	}

	bar := progressbar.Default(int64(len(pipelines)), "Inputsets")
	var failed []string

	for _, pipeline := range pipelines {
		inputsets, err := c.listInputsets(c.sourceOrg, c.sourceProject, pipeline.Identifier)
		if err != nil {
			failed = append(failed, fmt.Sprintf("Unable to list inputsets for pipeline %s [%s]", pipeline.Name, err))
			continue
		}

		bar.ChangeMax(bar.GetMax() + len(inputsets))

		for _, inputset := range inputsets {
			is, err := c.getInputset(c.sourceOrg, c.sourceProject, pipeline.Identifier, inputset.Identifier)
			if err == nil {
				newYaml := createYaml(is.Yaml, c.sourceOrg, c.sourceProject, c.targetOrg, c.targetProject)
				err = c.createInputset(c.targetOrg, c.targetProject, pipeline.Identifier, newYaml)
			}
			if err != nil {
				failed = append(failed, fmt.Sprintln(pipeline.Name, "/", err.Error()))
			}
			bar.Add(1)
		}
		bar.Add(1)
	}
	bar.Finish()

	reportFailed(failed, "inputsets:")
	return nil
}

func (c InputsetContext) listInputsets(org, project, pipelineIdentifier string) ([]*model.ListInputsetContent, error) {

	api := c.source
	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"accountIdentifier":  api.Account,
			"orgIdentifier":      org,
			"projectIdentifier":  project,
			"pipelineIdentifier": pipelineIdentifier,
			"inputSetType":       "ALL",
			"size":               "1000",
		}).
		Get(api.Url + "/pipeline/api/inputSets")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse(resp)
	}

	result := &model.ListInputsetResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	return result.Data.Content, nil
}

func (c InputsetContext) getInputset(org, project, pipelineIdentifier, isIdentifier string) (*model.GetInputsetData, error) {

	api := c.source
	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetHeader("Load-From-Cache", "false").
		SetPathParam("inputset", isIdentifier).
		SetQueryParams(map[string]string{
			"accountIdentifier":  api.Account,
			"orgIdentifier":      org,
			"projectIdentifier":  project,
			"pipelineIdentifier": pipelineIdentifier,
		}).
		Get(api.Url + "/pipeline/api/inputSets/{inputset}")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse(resp)
	}

	result := &model.GetInputsetResponse{}
	if err = json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (c InputsetContext) createInputset(org, project, pipelineIdentifier, yaml string) error {

	api := c.target
	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/yaml").
		SetBody(yaml).
		SetQueryParams(map[string]string{
			"accountIdentifier":  api.Account,
			"orgIdentifier":      org,
			"projectIdentifier":  project,
			"pipelineIdentifier": pipelineIdentifier,
		}).
		Post(api.Url + "/pipeline/api/inputSets")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}

	return nil
}
