package services

import (
	"encoding/json"
	"fmt"

	"github.com/jf781/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
)

const TRIGGER = "/pipeline/api/triggers"

type TriggerContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
}

func NewTriggerOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string) TriggerContext {
	return TriggerContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
	}
}

func (c TriggerContext) Move() error {

	triggers := []*model.TriggerContent{}

	// Leveraging listPipelines func from pipeline.go file
	pipelines, err := c.api.listPipelines(c.sourceOrg, c.sourceProject)
	if err != nil {
		return err
	}

	for _, p := range pipelines {
		triggerLists, err := c.api.listPipelineTriggers(p.Identifier, c.sourceOrg, c.sourceProject)
		if err != nil {
			return err
		}

		triggers = append(triggers, triggerLists...)
	}

	bar := progressbar.Default(int64(len(triggers)), "Triggers    ")
	var failed []string

	for _, t := range triggers {
		t.OrgIdentifier = c.targetOrg
		t.ProjectIdentifier = c.targetProject
		newYaml := createYaml(t.YAML, c.sourceOrg, c.sourceProject, c.targetOrg, c.targetProject)
		t.YAML = newYaml

		err = c.api.createPipelineTrigger(t)

		if err != nil {
			failed = append(failed, fmt.Sprintln(t.Name, "-", err.Error()))
		}
		bar.Add(1)
	}
	bar.Finish()

	reportFailed(failed, "Triggers:")
	return nil
}

func (api *ApiRequest) listPipelineTriggers(piplineId, org, project string) ([]*model.TriggerContent, error) {
	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     org,
			"projectIdentifier": project,
			"targetIdentifier":  piplineId,
			"size":              "100",
		}).
		Get(api.BaseURL + TRIGGER)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse(resp)
	}

	result := model.GetTriggerResposne{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	triggers := []*model.TriggerContent{}
	for _, c := range result.Data.Content {
		if c.TriggerStatus.Status == "SUCCESS" {
			newTrigger := c
			triggers = append(triggers, &newTrigger)
		} else {
			fmt.Println("Skipping trigger - Status is failed: ", c.Name)
		}
	}

	return triggers, nil
}

func (api *ApiRequest) createPipelineTrigger(trigger *model.TriggerContent) error {

	//api.Client.SetDebug(true)

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(trigger.YAML).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     trigger.OrgIdentifier,
			"projectIdentifier": trigger.ProjectIdentifier,
			"targetIdentifier":  trigger.Identifier,
		}).
		Post(api.BaseURL + TRIGGER)

	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}

	return nil
}
