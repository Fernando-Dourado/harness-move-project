package services

import (
	"encoding/json"
	"fmt"

	"github.com/jf781/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
)

const LIST_TEMPLATES_ENDPOINT = "/v1/orgs/{org}/projects/{project}/templates"
const GET_TEMPLATE_ENDPOINT = "/template/api/templates/{templateIdentifier}"
const CREATE_TEMPLATE_ENDPOINT = "/template/api/templates"

type TemplateContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
}

func NewTemplateOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string) TemplateContext {
	return TemplateContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
	}
}

func (c TemplateContext) Move() error {

	templates, err := c.listTemplates(c.sourceOrg, c.sourceProject)
	if err != nil {
		return err
	}

	bar := progressbar.Default(int64(len(templates)), "Templates   ")
	var failed []string

	for _, template := range templates {
		t, err := c.getTemplate(c.sourceOrg, c.sourceProject, template.Identifier, template.VersionLabel)
		if err == nil {
			newYaml := createYaml(t.Yaml, c.sourceOrg, c.sourceProject, c.targetOrg, c.targetProject)
			err = c.createTemplate(c.targetOrg, c.targetProject, newYaml)
		}
		if err != nil {
			failed = append(failed, fmt.Sprintln(template.Name, "-", err.Error()))
		}
		bar.Add(1)
	}
	bar.Finish()

	reportFailed(failed, "templates:")
	return nil
}

func (c TemplateContext) listTemplates(org, project string) (model.TemplateListResult, error) {

	api := c.api
	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetHeader("Harness-Account", api.Account).
		SetPathParam("org", org).
		SetPathParam("project", project).
		SetQueryParams(map[string]string{
			"limit": "1000",
		}).
		Get(BaseURL + LIST_TEMPLATES_ENDPOINT)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse(resp)
	}

	result := model.TemplateListResult{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c TemplateContext) getTemplate(org, project, templateIdentifier, versionLabel string) (*model.TemplateGetData, error) {

	api := c.api
	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetHeader("Load-From-Cache", "false").
		SetPathParam("templateIdentifier", templateIdentifier).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     org,
			"projectIdentifier": project,
			"versionLabel":      versionLabel,
		}).
		Get(BaseURL + GET_TEMPLATE_ENDPOINT)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse(resp)
	}

	result := model.TemplateGetResult{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (c TemplateContext) createTemplate(org, project, yaml string) error {

	api := c.api
	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(yaml).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     org,
			"projectIdentifier": project,
		}).
		Post(BaseURL + CREATE_TEMPLATE_ENDPOINT)
	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}

	return nil
}
