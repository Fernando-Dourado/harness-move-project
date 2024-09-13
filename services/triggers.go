package services

import (
	"encoding/json"

	"github.com/jf781/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
	"go.uber.org/zap"
)

const TRIGGER = "/pipeline/api/triggers"

type TriggerContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
	logger        *zap.Logger
}

func NewTriggerOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string, logger *zap.Logger) TriggerContext {
	return TriggerContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
		logger:        logger,
	}
}

func (c TriggerContext) Copy() error {

	c.logger.Info("Copying triggers",
		zap.String("project", c.sourceProject),
	)

	triggers := []*model.TriggerContent{}

	// Leveraging listPipelines func from pipeline.go file
	pipelines, err := c.api.listPipelines(c.sourceOrg, c.sourceProject, c.logger)
	if err != nil {
		c.logger.Error("Failed to retrive pipelines for project",
			zap.String("Project", c.sourceProject),
			zap.Error(err),
		)
		return err
	}

	for _, p := range pipelines {
		triggerLists, err := c.api.listPipelineTriggers(p.Identifier, c.sourceOrg, c.sourceProject, c.logger)
		if err != nil {
			c.logger.Error("Getting pipeline details",
				zap.String("pipeline", p.Identifier),
				zap.String("Project", c.sourceProject),
				zap.Error(err),
			)
			return err
		}

		triggers = append(triggers, triggerLists...)
	}

	bar := progressbar.Default(int64(len(triggers)), "Triggers    ")

	for _, t := range triggers {

		c.logger.Info("Processing trigger",
			zap.String("trigger", t.Name),
			zap.String("targetProject", c.targetProject),
		)

		t.OrgIdentifier = c.targetOrg
		t.ProjectIdentifier = c.targetProject
		newYaml := createYaml(t.YAML, c.sourceOrg, c.sourceProject, c.targetOrg, c.targetProject)
		t.YAML = newYaml

		err = c.api.createPipelineTrigger(t, c.logger)

		if err != nil {
			c.logger.Error("Failed to create trigger",
				zap.String("trigger", t.Name),
				zap.Error(err),
			)
		}
		bar.Add(1)
	}
	bar.Finish()

	return nil
}

func (api *ApiRequest) listPipelineTriggers(piplineId, org, project string, logger *zap.Logger) ([]*model.TriggerContent, error) {

	logger.Info("Fetching triggers",
		zap.String("org", org),
		zap.String("project", project),
	)

	IncrementApiCalls()

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
		logger.Error("Failed to request to list of triggers",
			zap.Error(err),
		)
		return nil, err
	}
	if resp.IsError() {
		logger.Error("Error response from API when listing triggers",
			zap.String("response",
				resp.String(),
			),
		)
		return nil, handleErrorResponse(resp)
	}

	result := model.GetTriggerResposne{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		logger.Error("Failed to parse response from API",
			zap.Error(err),
		)
		return nil, err
	}

	triggers := []*model.TriggerContent{}
	for _, c := range result.Data.Content {
		if c.TriggerStatus.Status == "SUCCESS" {
			newTrigger := c
			triggers = append(triggers, &newTrigger)
		} else {
			logger.Warn("Skipping trigger because the status is not active",
				zap.String("trigger", c.Name),
				zap.String("t", c.TriggerStatus.Status),
			)
		}
	}

	return triggers, nil
}

func (api *ApiRequest) createPipelineTrigger(trigger *model.TriggerContent, logger *zap.Logger) error {

	logger.Info("Creating trigger",
		zap.String("trigger", trigger.Name),
		zap.String("project", trigger.ProjectIdentifier),
	)

	IncrementApiCalls()

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
		logger.Error("Failed to send request to create ",
			zap.String("trigger", trigger.Name),
			zap.Error(err),
		)
		return err
	}
	if resp.IsError() {
		var errorResponse map[string]interface{}
		if err := json.Unmarshal(resp.Body(), &errorResponse); err == nil {
			if code, ok := errorResponse["code"].(string); ok && code == "DUPLICATE_FIELD" {
				// Log as a warning and skip the error
				logger.Info("Duplicate trigger found, ignoring error",
					zap.String("trigger", trigger.Name),
				)
				return nil
			}
		} else {
			logger.Error(
				"Error response from API when creating ",
				zap.String("trigger", trigger.Name),
				zap.String("response",
					resp.String(),
				),
			)
		}
		return handleErrorResponse(resp)
	}

	return nil
}
