package services

import (
	"encoding/json"

	"github.com/jf781/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
	"go.uber.org/zap"
)

const TAGS = "/ng/api/serviceaccount"

type TagsContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
	logger        *zap.Logger
}

func NewTagOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string, logger *zap.Logger) TagsContext {
	return TagsContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
		logger:        logger,
	}
}

func (c TagsContext) Copy() error {

	c.logger.Info("Copying Tags",
		zap.String("project", c.sourceProject),
	)

	projectTags := []*model.Tag{}

	// Leveraging listEnvironments func from environment.go file
	envs, err := c.api.listEnvironments(c.sourceOrg, c.sourceProject, c.logger)
	if err != nil {
		c.logger.Error("Failed to retrive tags",
			zap.String("Project", c.sourceProject),
			zap.Error(err),
		)
		return nil
	}

	for _, env := range envs {
		e := env.Environment

		envTags, err := c.api.listTags(e.Name, c.sourceOrg, c.sourceProject, c.logger)

		if err != nil {
			c.logger.Error("Failed to retrive environments",
				zap.String("Project", c.sourceProject),
				zap.Error(err),
			)
		}

		projectTags = append(projectTags, envTags...)
	}

	bar := progressbar.Default(int64(len(projectTags)), "Project Tags    ")

	for _, t := range projectTags {

		IncrementTagsTotal()

		c.logger.Info("Processing tag",
			zap.String("tag", t.Name),
			zap.String("targetProject", c.targetProject),
		)

		newTag := &model.CreateTagRequest{
			OrgIdentifier:     c.targetOrg,
			ProjectIdentifier: c.targetProject,
			Name:              t.Name,
			Identifier:        t.Identifier,
		}

		err = c.api.createTags(newTag, c.logger)

		if err != nil {
			c.logger.Error("Failed to create tag",
				zap.String("tag", t.Name),
				zap.Error(err),
			)
		} else {
			IncrementTagsMoved()
		}
		bar.Add(1)
	}
	bar.Finish()

	return nil
}

func (api *ApiRequest) listTags(environment, org, project string, logger *zap.Logger) ([]*model.Tag, error) {

	logger.Info("Fetching infrastructure",
		zap.String("org", org),
		zap.String("project", project),
		zap.String("environment", environment),
	)

	IncrementApiCalls()

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"accountIdentifier":     api.Account,
			"orgIdentifier":         org,
			"projectIdentifier":     project,
			"environmentIdentifier": environment,
		}).
		Get(api.BaseURL + TAGS)
	if err != nil {
		logger.Error("Failed to request to list of tags",
			zap.Error(err),
		)
		return nil, err
	}
	if resp.IsError() {
		logger.Error("Error response from API when listing tags",
			zap.String("response",
				resp.String(),
			),
		)
		return nil, handleErrorResponse(resp)
	}

	result := model.TagListResult{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		logger.Error("Failed to parse response from API",
			zap.Error(err),
		)
		return nil, err
	}

	environmentTags := []*model.Tag{}
	for _, c := range result.Tags {
		tags := c
		environmentTags = append(environmentTags, &tags)
	}

	return environmentTags, nil
}

func (api *ApiRequest) createTags(tag *model.CreateTagRequest, logger *zap.Logger) error {

	logger.Info("Creating tag",
		zap.String("tag", tag.Name),
		zap.String("project", tag.ProjectIdentifier),
	)

	IncrementApiCalls()

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(tag).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     tag.OrgIdentifier,
			"projectIdentifier": tag.ProjectIdentifier,
		}).
		Post(api.BaseURL + TAGS)

	if err != nil {
		logger.Error("Failed to send request to create ",
			zap.String("tag", tag.Name),
			zap.Error(err),
		)
		return err
	}
	if resp.IsError() {
		var errorResponse map[string]interface{}
		if err := json.Unmarshal(resp.Body(), &errorResponse); err == nil {
			if code, ok := errorResponse["code"].(string); ok && code == "DUPLICATE_FIELD" {
				// Log as a warning and skip the error
				logger.Info("Duplicate tag found, ignoring error",
					zap.String("tag", tag.Name),
				)
				return nil
			}
		} else {
			logger.Error(
				"Error response from API when creating ",
				zap.String("tag", tag.Name),
				zap.String("response",
					resp.String(),
				),
			)
		}
		return handleErrorResponse(resp)
	}

	return nil
}
