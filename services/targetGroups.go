package services

import (
	"encoding/json"

	"github.com/jf781/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
	"go.uber.org/zap"
)

const TARGETGROUPS = "/cf/admin/segments"

type TargetGroupContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
	logger        *zap.Logger
}

func NewTargetGroups(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string, logger *zap.Logger) TargetGroupContext {
	return TargetGroupContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
		logger:        logger,
	}
}

func (c TargetGroupContext) Copy() error {

	c.logger.Info("Copying target groups",
		zap.String("project", c.sourceProject),
	)

	envs, err := c.api.listEnvironments(c.sourceOrg, c.sourceProject, c.logger)
	if err != nil {
		c.logger.Error("Failed to retrive environments",
			zap.String("Project", c.sourceProject),
			zap.Error(err),
		)
		return err
	}

	bar := progressbar.Default(int64(len(envs)), "Targets    ")

	for _, env := range envs {
		e := env.Environment
		targetGroups, err := c.api.listTargetGroups(c.sourceOrg, c.sourceProject, e.Identifier, c.logger)
		if err != nil {
			c.logger.Error("Failed to retrive target group",
				zap.String("Project", c.sourceProject),
				zap.String("Environment", e.Identifier),
				zap.Error(err),
			)
			continue
		}

		bar.ChangeMax(bar.GetMax() + len(targetGroups))

		for _, targetGroup := range targetGroups {
			i := targetGroup
			c.logger.Info("Processing target group",
				zap.String("target group", i.Name),
				zap.String("targetProject", c.targetProject),
			)

			includedNames := []string{}
			for _, included := range i.Included {
				iName := included.Name // Declare iName inside the loop
				includedNames = append(includedNames, iName)
			}

			err := c.api.createTargetGroups(&model.NewTargetGroup{
				Name:         i.Name,
				Identifier:   i.Identifier,
				Org:          c.targetOrg,
				Project:      c.targetProject,
				Account:      i.Account,
				Environment:  e.Identifier,
				Included:     includedNames,
				Excluded:     i.Excluded,
				Rules:        i.Rules,
				ServingRules: i.ServingRules,
			}, c.logger)
			if err != nil {
				c.logger.Error("Failed to create target group",
					zap.String("target group", i.Name),
					zap.Error(err),
				)
			}
			bar.Add(1)
		}
		bar.Add(1)
	}
	bar.Finish()

	return nil
}

func (api *ApiRequest) listTargetGroups(org, project, envId string, logger *zap.Logger) ([]*model.TargetGroups, error) {

	logger.Info("Fetching target groups",
		zap.String("org", org),
		zap.String("project", project),
		zap.String("environment", envId),
	)

	IncrementApiCalls()

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"accountIdentifier":     api.Account,
			"orgIdentifier":         org,
			"projectIdentifier":     project,
			"environmentIdentifier": envId,
			"size":                  "1000",
		}).
		Get(api.BaseURL + TARGETGROUPS)
	if err != nil {
		logger.Error("Failed to request to list of target groups",
			zap.Error(err),
		)
		return nil, err
	}
	if resp.IsError() {
		logger.Error("Error response from API when listing target groups",
			zap.String("response",
				resp.String(),
			),
		)
		return nil, handleErrorResponse(resp)
	}

	result := model.TargetGroupListResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		logger.Error("Failed to parse response from API",
			zap.Error(err),
		)
		return nil, err
	}

	targetGroups := []*model.TargetGroups{}
	for _, c := range result.TargetGroups {
		tg := c
		targetGroups = append(targetGroups, &tg)
	}

	return targetGroups, nil
}

func (api *ApiRequest) createTargetGroups(targetGroup *model.NewTargetGroup, logger *zap.Logger) error {

	logger.Info("Creating target group",
		zap.String("target group", targetGroup.Name),
		zap.String("project", targetGroup.Project),
	)

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(targetGroup).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     targetGroup.Org,
		}).
		Post(api.BaseURL + TARGETGROUPS)
	if err != nil {
		logger.Error("Failed to send request to create ",
			zap.String("target group", targetGroup.Name),
			zap.Error(err),
		)
		return err
	}
	if resp.IsError() {
		var errorResponse map[string]interface{}
		if err := json.Unmarshal(resp.Body(), &errorResponse); err == nil {
			if code, ok := errorResponse["code"].(string); ok && code == "409" {
				// Log as a warning and skip the error
				logger.Info("Duplicate target group found, ignoring error",
					zap.String("target group", targetGroup.Name),
				)
				return nil
			}
		} else {
			logger.Error(
				"Error response from API when creating ",
				zap.String("target group", targetGroup.Name),
				zap.String("response",
					resp.String(),
				),
			)
		}
		return handleErrorResponse(resp)
	}

	return nil
}
