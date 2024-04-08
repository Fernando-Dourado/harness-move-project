package services

import (
	"fmt"

	"github.com/Fernando-Dourado/harness-move-project/rest"
	"github.com/schollz/progressbar/v3"
)

const LIST_PIPELINES = ""
const GET_PIPELINE = ""
const CREATE_PIPELINE = "/pipeline/api/pipelines/v2"

type PipelineContext struct {
	rest          *rest.PipelineRestContext
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
}

func NewPipelineOperation(rest *rest.PipelineRestContext, sourceOrg, sourceProject, targetOrg, targetProject string) *PipelineContext {
	return &PipelineContext{
		rest:          rest,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
	}
}

func (c PipelineContext) Move() error {

	pipelines, err := c.rest.ListPipelines(c.sourceOrg, c.sourceProject)
	if err != nil {
		return err
	}

	bar := progressbar.Default(int64(len(pipelines)), "Pipelines   ")
	var failed []string

	for _, pipe := range pipelines {
		pipeData, err := c.rest.GetPipeline(c.sourceOrg, c.sourceProject, pipe.Identifier)
		if err == nil {
			newYaml := createYaml(pipeData.YAMLPipeline, c.sourceOrg, c.sourceProject, c.targetOrg, c.targetProject)
			err = c.rest.CreatePipeline(c.targetOrg, c.targetProject, newYaml)
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

// func (c PipelineContext) createPipeline(org, project, yaml string) error {

// 	resp, err := c.api.Client.R().
// 		SetHeader("x-api-key", c.api.Token).
// 		SetHeader("Content-Type", "application/yaml").
// 		SetBody(yaml).
// 		SetQueryParams(map[string]string{
// 			"accountIdentifier": c.api.Account,
// 			"orgIdentifier":     org,
// 			"projectIdentifier": project,
// 		}).
// 		Post(BaseURL + CREATE_PIPELINE)
// 	if err != nil {
// 		return err
// 	}
// 	if resp.IsError() {
// 		return handleErrorResponse(resp)
// 	}

// 	return nil
// }
