package services

import (
	"encoding/json"
	"fmt"

	"github.com/jf781/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
	"go.uber.org/zap"
)

const FEATFLAGS = "/cf/admin/features"

type FeatureContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
	logger        *zap.Logger
}

func NewFeatureFlagOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string, logger *zap.Logger) FeatureContext {
	return FeatureContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
		logger:        logger,
	}
}

func (c FeatureContext) Copy() error {

	c.logger.Info("Copying Feature Flags",
		zap.String("project", c.sourceProject),
	)

	featureFlags, err := c.api.listFeatureFlags(c.sourceOrg, c.sourceProject, c.logger)
	if err != nil {
		c.logger.Error("Failed to retrive feature flags",
			zap.String("Project", c.sourceProject),
			zap.Error(err),
		)
		return err
	}

	bar := progressbar.Default(int64(len(featureFlags)), "Feature Flags    ")

	for _, f := range featureFlags {
		if f.Tags == nil {
			f.Tags = []string{}
		}

		IncrementFeatureFlagsTotal()

		c.logger.Info("Processing feature flag",
			zap.String("feature flag", f.Name),
			zap.String("targetProject", c.targetProject),
		)

		err = c.api.createFeatureFlags(&model.CreateFeatureFlag{
			OrgIdentifier:       c.targetOrg,
			ProjectIdentifier:   c.targetProject,
			Archived:            f.Archived,
			CreatedAt:           f.CreatedAt,
			DefaultOffVariation: f.DefaultOffVariation,
			DefaultOnVariation:  f.DefaultOnVariation,
			Description:         f.Description,
			EnvProperties:       f.EnvProperties,
			Evaluation:          f.Evaluation,
			Identifier:          f.Identifier,
			Kind:                f.Kind,
			ModifiedAt:          f.ModifiedAt,
			Name:                f.Name,
			Owner:               fmt.Sprint(f.Owner),
			Permanent:           f.Permanent,
			Prerequisites:       f.Prerequisites,
			Project:             c.targetProject,
			Services:            f.Services,
			Tags:                f.Tags,
			Variations:          f.Variations,
		}, c.logger)

		if err != nil {
			c.logger.Error("Failed to create feature flag",
				zap.String("feature flag", f.Name),
				zap.Error(err),
			)
			return err
		} else {
			IncrementFeatureFlagsMoved()
		}
		bar.Add(1)
	}
	bar.Finish()

	return nil
}

func (api *ApiRequest) listFeatureFlags(org, project string, logger *zap.Logger) ([]*model.FeatureFlag, error) {

	logger.Info("Fetching feature flags",
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
		}).
		Get(api.BaseURL + FEATFLAGS)
	if err != nil {
		logger.Error("Failed to request to list of feature flags",
			zap.Error(err),
		)
		return nil, err
	}
	if resp.IsError() {
		logger.Error("Error response from API when listing feature flags",
			zap.String("response",
				resp.String(),
			),
		)
		return nil, handleErrorResponse(resp)
	}

	result := model.FeatureFlagListResult{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		logger.Error("Failed to parse response from API",
			zap.Error(err),
		)
		return nil, err
	}

	featureFlags := []*model.FeatureFlag{}
	for _, c := range result.Features {
		feature := c
		featureFlags = append(featureFlags, &feature)
	}

	return featureFlags, nil
}

func (api *ApiRequest) createFeatureFlags(featureFlag *model.CreateFeatureFlag, logger *zap.Logger) error {

	logger.Info("Creating feature flag",
		zap.String("feature flag", featureFlag.Name),
		zap.String("project", featureFlag.ProjectIdentifier),
	)

	IncrementApiCalls()

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(featureFlag).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     featureFlag.OrgIdentifier,
			"projectIdentifier": featureFlag.ProjectIdentifier,
		}).
		Post(api.BaseURL + FEATFLAGS)

	if err != nil {
		logger.Error("Failed to send request to create ",
			zap.String("feature flag", featureFlag.Name),
			zap.Error(err),
		)
		return err
	}
	if resp.IsError() {
		var errorResponse map[string]interface{}
		if err := json.Unmarshal(resp.Body(), &errorResponse); err == nil {
			if code, ok := errorResponse["code"].(string); ok && code == "409" {
				// Log as a warning and skip the error
				logger.Info("Duplicate feature flag found, ignoring error",
					zap.String("feature flag", featureFlag.Name),
				)
				return nil
			}
		} else {
			logger.Error(
				"Error response from API when creating ",
				zap.String("feature flag", featureFlag.Name),
				zap.String("response",
					resp.String(),
				),
			)
		}
		return handleErrorResponse(resp)
	}

	return nil
}
