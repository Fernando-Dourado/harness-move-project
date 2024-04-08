package services

import (
	"fmt"

	"github.com/Fernando-Dourado/harness-move-project/rest"
	"github.com/schollz/progressbar/v3"
)

type InputsetContext struct {
	rest          *rest.InputsetRestContext
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
	p             *rest.PipelineRestContext
}

func NewInputsetOperation(rest *rest.InputsetRestContext, p *rest.PipelineRestContext, sourceOrg, sourceProject, targetOrg, targetProject string) *InputsetContext {
	return &InputsetContext{
		rest:          rest,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
		p:             p,
	}
}

func (c InputsetContext) Move() error {

	pipelines, err := c.p.ListPipelines(c.sourceOrg, c.sourceProject)
	if err != nil {
		return err
	}

	bar := progressbar.Default(int64(len(pipelines)), "Inputsets")
	var failed []string

	for _, pipeline := range pipelines {
		inputsets, err := c.rest.ListInputsets(c.sourceOrg, c.sourceProject, pipeline.Identifier)
		if err != nil {
			failed = append(failed, fmt.Sprintf("Unable to list inputsets for pipeline %s [%s]", pipeline.Name, err))
			continue
		}

		bar.ChangeMax(bar.GetMax() + len(inputsets))

		for _, inputset := range inputsets {
			is, err := c.rest.GetInputset(c.sourceOrg, c.sourceProject, pipeline.Identifier, inputset.Identifier)
			if err == nil {
				newYaml := createYaml(is.Yaml, c.sourceOrg, c.sourceProject, c.targetOrg, c.targetProject)
				err = c.rest.CreateInputset(c.targetOrg, c.targetProject, pipeline.Identifier, newYaml)
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

// func (api *ApiRequest) listInputsets(org, project, pipelineIdentifier string) ([]*model.ListInputsetContent, error) {

// 	resp, err := api.Client.R().
// 		SetHeader("x-api-key", api.Token).
// 		SetHeader("Content-Type", "application/json").
// 		SetQueryParams(map[string]string{
// 			"accountIdentifier":  api.Account,
// 			"orgIdentifier":      org,
// 			"projectIdentifier":  project,
// 			"pipelineIdentifier": pipelineIdentifier,
// 			"inputSetType":       "ALL",
// 			"size":               "1000",
// 		}).
// 		Get(BaseURL + "/pipeline/api/inputSets")
// 	if err != nil {
// 		return nil, err
// 	}
// 	if resp.IsError() {
// 		return nil, handleErrorResponse(resp)
// 	}

// 	result := &model.ListInputsetResponse{}
// 	err = json.Unmarshal(resp.Body(), &result)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return result.Data.Content, nil
// }

// func (api *ApiRequest) getInputset(org, project, pipelineIdentifier, isIdentifier string) (*model.GetInputsetData, error) {

// 	resp, err := api.Client.R().
// 		SetHeader("x-api-key", api.Token).
// 		SetHeader("Content-Type", "application/json").
// 		SetHeader("Load-From-Cache", "false").
// 		SetPathParam("inputset", isIdentifier).
// 		SetQueryParams(map[string]string{
// 			"accountIdentifier":  api.Account,
// 			"orgIdentifier":      org,
// 			"projectIdentifier":  project,
// 			"pipelineIdentifier": pipelineIdentifier,
// 		}).
// 		Get(BaseURL + "/pipeline/api/inputSets/{inputset}")
// 	if err != nil {
// 		return nil, err
// 	}
// 	if resp.IsError() {
// 		return nil, handleErrorResponse(resp)
// 	}

// 	result := &model.GetInputsetResponse{}
// 	if err = json.Unmarshal(resp.Body(), &result); err != nil {
// 		return nil, err
// 	}

// 	return result.Data, nil
// }

// func (api *ApiRequest) createInputset(org, project, pipelineIdentifier, yaml string) error {

// 	resp, err := api.Client.R().
// 		SetHeader("x-api-key", api.Token).
// 		SetHeader("Content-Type", "application/yaml").
// 		SetBody(yaml).
// 		SetQueryParams(map[string]string{
// 			"accountIdentifier":  api.Account,
// 			"orgIdentifier":      org,
// 			"projectIdentifier":  project,
// 			"pipelineIdentifier": pipelineIdentifier,
// 		}).
// 		Post(BaseURL + "/pipeline/api/inputSets")
// 	if err != nil {
// 		return err
// 	}
// 	if resp.IsError() {
// 		return handleErrorResponse(resp)
// 	}

// 	return nil
// }
