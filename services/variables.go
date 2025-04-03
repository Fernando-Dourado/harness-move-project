package services

import (
	"encoding/json"
	"fmt"

	"github.com/Fernando-Dourado/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
)

type VariableContext struct {
	source        *SourceRequest
	target        *TargetRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
}

func NewVariableOperation(sourceApi *SourceRequest, targetApi *TargetRequest, sourceOrg, sourceProject, targetOrg, targetProject string) VariableContext {
	return VariableContext{
		source:        sourceApi,
		target:        targetApi,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
	}
}

func (c VariableContext) Move() error {

	variables, err := c.listVariables(c.sourceOrg, c.sourceProject)
	if err != nil {
		return err
	}

	bar := progressbar.Default(int64(len(variables)), "Variables")
	var failed []string

	for _, v := range variables {
		v.OrgIdentifier = c.targetOrg
		v.ProjectIdentifier = c.targetProject

		err = c.createVariable(&model.CreateVariableRequest{
			Variable: v,
		})
		if err != nil {
			failed = append(failed, fmt.Sprintln(v.Name, "-", err.Error()))
		}
		bar.Add(1)
	}
	bar.Finish()

	reportFailed(failed, "variables")
	return nil
}

func (c VariableContext) listVariables(org, project string) ([]*model.Variable, error) {

	api := c.source
	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     org,
			"projectIdentifier": project,
			"size":              "1000",
		}).
		Get(BaseURL + "/ng/api/variables")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse(resp)
	}

	result := model.GetVariablesResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	variables := []*model.Variable{}
	for _, c := range result.Data.Content {
		variables = append(variables, c.Variable)
	}

	return variables, nil
}

func (c VariableContext) createVariable(variable *model.CreateVariableRequest) error {

	api := c.target
	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(variable).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
		}).
		Post(BaseURL + "/ng/api/variables")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}

	return nil
}
