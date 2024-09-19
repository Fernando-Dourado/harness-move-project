package services

import (
	"encoding/json"

	"github.com/jf781/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
	"go.uber.org/zap"
)

const TARGETS = "/cf/admin/targets"

type TargetContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
	logger        *zap.Logger
}

func NewTargets(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string, logger *zap.Logger) TargetContext {
	return TargetContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
		logger:        logger,
	}
}

func (c TargetContext) Copy() error {

	c.logger.Info("Copying targets",
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

	bar := progressbar.Default(int64(len(envs)), "Target Groups    ")
	var failed []string

	for _, env := range envs {
		e := env.Environment
		targets, err := c.api.listTargets(c.sourceOrg, c.sourceProject, e.Identifier, c.logger)
		if err != nil {
			c.logger.Error("Failed to retrive targets",
				zap.String("Project", c.sourceProject),
				zap.Error(err),
			)
			continue
		}

		bar.ChangeMax(bar.GetMax() + len(targets))

		for _, target := range targets {

			IncrementTargetsTotal()

			i := target

			c.logger.Info("Processing target",
				zap.String("target", i.Name),
				zap.String("targetProject", c.targetProject),
			)

			err := c.api.createTarget(&model.Target{
				Name:        i.Name,
				Identifier:  i.Identifier,
				Org:         c.targetOrg,
				Project:     c.targetProject,
				Environment: i.Environment,
				Attributes:  i.Attributes,
				Segments:    i.Segments,
			}, c.logger)

			if err != nil {
				c.logger.Error("Failed to create target",
					zap.String("target", i.Name),
					zap.Error(err),
				)
			} else {
				IncrementTargetsMoved()
			}
			bar.Add(1)
		}
		bar.Add(1)
	}
	bar.Finish()

	reportFailed(failed, "targets:")
	return nil
}

func (api *ApiRequest) listTargets(org, project, envId string, logger *zap.Logger) ([]*model.Target, error) {

	logger.Info("Fetching target",
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
		Get(api.BaseURL + TARGETS)
	if err != nil {
		logger.Error("Failed to request to list of targets",
			zap.Error(err),
		)
		return nil, err
	}
	if resp.IsError() {
		logger.Error("Error response from API when listing targets",
			zap.String("response",
				resp.String(),
			),
		)
		return nil, handleErrorResponse(resp)
	}

	result := model.TargetListResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		logger.Error("Failed to parse response from API",
			zap.Error(err),
		)
		return nil, err
	}

	targets := []*model.Target{}
	for _, c := range result.Targets {
		target := c
		targets = append(targets, &target)
	}

	return targets, nil
}

func (api *ApiRequest) createTarget(target *model.Target, logger *zap.Logger) error {

	logger.Info("Creating target",
		zap.String("target", target.Name),
		zap.String("project", target.Project),
	)

	IncrementApiCalls()

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(target).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     target.Org,
		}).
		Post(api.BaseURL + TARGETS)
	if err != nil {
		return err
	}
	if resp.IsError() {
		var errorResponse map[string]interface{}
		if err := json.Unmarshal(resp.Body(), &errorResponse); err == nil {
			if code, ok := errorResponse["code"].(string); ok && code == "409" {
				// Log as a warning and skip the error
				logger.Info("Duplicate target found, ignoring error",
					zap.String("target", target.Name),
				)
				return nil
			}
		} else {
			logger.Error(
				"Error response from API when creating ",
				zap.String("target", target.Name),
				zap.String("response",
					resp.String(),
				),
			)
		}
		return handleErrorResponse(resp)
	}

	return nil
}
