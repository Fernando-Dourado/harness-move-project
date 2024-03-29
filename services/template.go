package services

import (
	"encoding/json"
	"fmt"

	"github.com/Fernando-Dourado/harness-move-project/model"
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

	templates, err := c.api.listTemplates(c.sourceOrg, c.sourceProject)
	if err != nil {
		return err
	}

	for _, template := range templates {
		t, err := c.api.getTemplate(c.sourceOrg, c.sourceProject, template.Identifier, template.VersionLabel)
		if err != nil {
			return err
		}

		newYaml := createYaml(t.Yaml, c.sourceOrg, c.sourceProject, c.targetOrg, c.targetProject)
		if err := c.api.createTemplate(c.targetOrg, c.targetProject, newYaml); err != nil {
			return err
		}
	}

	return nil
}

func (api *ApiRequest) listTemplates(org, project string) (model.TemplateListResult, error) {

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
		return nil, fmt.Errorf(resp.Status())
	}

	result := model.TemplateListResult{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (api *ApiRequest) getTemplate(org, project, templateIdentifier, versionLabel string) (*model.TemplateGetData, error) {

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
		return nil, fmt.Errorf(resp.Status())
	}

	result := model.TemplateGetResult{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (api *ApiRequest) createTemplate(org, project, yaml string) error {

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
		return handleCreateErrorResponse(resp)
	}

	return nil
}
