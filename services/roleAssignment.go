package services

import (
	"encoding/json"
	"fmt"

	"github.com/jf781/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
)

const ROLEASSIGNMENT = "/authz/api/roleassignments"

type RoleAssignmentContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
}

func NewRoleAssignmentOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string) RoleAssignmentContext {
	return RoleAssignmentContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
	}
}

func (c RoleAssignmentContext) Move() error {

	roleAssignments, err := c.api.listRoleAssignments(c.sourceOrg, c.sourceProject)
	if err != nil {
		return err
	}

	bar := progressbar.Default(int64(len(roleAssignments)), "Roles Assignments    ")
	var failed []string

	for _, r := range roleAssignments {

		role := &model.NewRoleAssignment{
			Identifier:              r.Identifier,
			ResourceGroupIdentifier: r.ResourceGroupIdentifier,
			RoleIdentifier:          r.RoleIdentifier,
			Principal:               r.Principal,
			OrgIdentifier:           c.targetOrg,
			ProjectIdentifier:       c.targetProject,
		}

		err = c.api.createRoleAssignment(role)

		if err != nil {
			failed = append(failed, fmt.Sprintln(r.Identifier, "-", err.Error()))
		}
		bar.Add(1)
	}
	bar.Finish()

	reportFailed(failed, "Role Assingments:")
	return nil
}

func (api *ApiRequest) listRoleAssignments(org, project string) ([]*model.ExistingRoleAssignment, error) {

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     org,
			"projectIdentifier": project,
			"pageSize":          "100",
		}).
		Get(api.BaseURL + ROLEASSIGNMENT)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse(resp)
	}

	result := model.GetRoleAssignmentResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	roleAssignments := []*model.ExistingRoleAssignment{}
	for _, c := range result.Data.Content {
		newRoleAssignment := c.RoleAssignment
		roleAssignments = append(roleAssignments, &newRoleAssignment)
	}

	return roleAssignments, nil
}

func (api *ApiRequest) createRoleAssignment(role *model.NewRoleAssignment) error {

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(role).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     role.OrgIdentifier,
			"projectIdentifier": role.ProjectIdentifier,
		}).
		Post(api.BaseURL + ROLEASSIGNMENT)

	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}

	return nil
}
