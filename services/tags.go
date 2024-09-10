package services

import (
	"encoding/json"
	"fmt"

	"github.com/jf781/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
)

const TAGS = "/ng/api/serviceaccount"

type TagsContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
}

func NewTagOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string) TagsContext {
	return TagsContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
	}
}

func (c TagsContext) Move() error {

	projectTags := []*model.Tag{}

	// Leveraging listEnvironments func from environment.go file
	envs, err := c.api.listEnvironments(c.sourceOrg, c.sourceProject)
	if err != nil {
		return nil
	}

	var failed []string

	for _, env := range envs {
		e := env.Environment

		envTags, err := c.api.listTags(e.Name, c.sourceOrg, c.sourceProject)

		if err != nil {
			failed = append(failed, fmt.Sprintln(e.Name, "-", err.Error()))
		}

		projectTags = append(projectTags, envTags...)
	}

	bar := progressbar.Default(int64(len(projectTags)), "Project Tags    ")

	for _, t := range projectTags {

		newTag := &model.CreateTagRequest{
			OrgIdentifier:     c.targetOrg,
			ProjectIdentifier: c.targetProject,
			Name:              t.Name,
			Identifier:        t.Identifier,
		}

		err = c.api.createTags(newTag)

		if err != nil {
			failed = append(failed, fmt.Sprintln(t.Name, "-", err.Error()))
		}
		bar.Add(1)
	}
	bar.Finish()

	reportFailed(failed, "Project Tags:")
	return nil
}

func (api *ApiRequest) listTags(environment, org, project string) ([]*model.Tag, error) {

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
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse(resp)
	}

	result := model.TagListResult{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	environmentTags := []*model.Tag{}
	for _, c := range result.Tags {
		tags := c
		environmentTags = append(environmentTags, &tags)
	}

	return environmentTags, nil
}

func (api *ApiRequest) createTags(tag *model.CreateTagRequest) error {

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
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}

	return nil
}
