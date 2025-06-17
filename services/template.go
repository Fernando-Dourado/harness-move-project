package services

import (
	"encoding/json"
	"fmt"

	"github.com/Fernando-Dourado/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
)

const LIST_TEMPLATES_ENDPOINT = "/v1/orgs/{org}/projects/{project}/templates"
const GET_TEMPLATE_ENDPOINT = "/template/api/templates/{templateIdentifier}"
const CREATE_TEMPLATE_ENDPOINT = "/template/api/templates"

type TemplateContext struct {
	source        *SourceRequest
	target        *TargetRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
}

func NewTemplateOperation(sourceApi *SourceRequest, targetApi *TargetRequest, st *SourceTarget) TemplateContext {
	return TemplateContext{
		source:        sourceApi,
		target:        targetApi,
		sourceOrg:     st.SourceOrg,
		sourceProject: st.SourceProject,
		targetOrg:     st.TargetOrg,
		targetProject: st.TargetProject,
	}
}

func (c TemplateContext) Move() error {

	templates, err := listTemplates(c.source, c.sourceOrg, c.sourceProject)
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

func listTemplates(s *SourceRequest, org, project string) (model.TemplateListResult, error) {

	params := map[string]string{
		"org": org,
	}
	if len(project) > 0 {
		params["project"] = project
	}

	endpoint := LIST_TEMPLATES_ENDPOINT
	if len(project) == 0 {
		endpoint = "/v1/orgs/{org}/templates"
	}

	resp, err := s.Client.R().
		SetHeader("x-api-key", s.Token).
		SetHeader("Content-Type", "application/json").
		SetHeader("Harness-Account", s.Account).
		SetPathParams(params).
		SetQueryParams(map[string]string{
			"limit": "1000",
		}).
		Get(s.Url + endpoint)
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

	params := createQueryParams(c.source.Account, org, project)
	params["versionLabel"] = versionLabel

	api := c.source
	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetHeader("Load-From-Cache", "false").
		SetPathParam("templateIdentifier", templateIdentifier).
		SetQueryParams(params).
		Get(api.Url + GET_TEMPLATE_ENDPOINT)
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

	params := createQueryParams(c.target.Account, org, project)

	api := c.target
	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(yaml).
		SetQueryParams(params).
		Post(api.Url + CREATE_TEMPLATE_ENDPOINT)
	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}

	return nil
}
