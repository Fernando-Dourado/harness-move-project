package services

import (
	"encoding/json"
	"fmt"

	"github.com/jf781/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
)

const ROLES = "/authz/api/roleassignments"

type RoleContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
}

func NewAccessControlOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string) RoleContext {
	return RoleContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
	}
}

func (c RoleContext) Move() error {

	roles, err := c.api.ListRoles(c.sourceOrg, c.sourceProject)
	if err != nil {
		return err
	}

	// report(roles)
	// return nil

	bar := progressbar.Default(int64(len(roles)), "Roles")
	var failed []string

	for _, r := range roles {
		r.OrgIdentifier = c.targetOrg
		r.ProjectIdentifier = c.targetProject

		fmt.Printf("Role Assignment: %+v \n", r.Principal)

		err = c.api.CreateRoleAssignment(&model.CreateRoleAssignment{
			role: CreateRoleAssignment{
				ResourceGroupIdentifier: r.ResourceGroupIdentifier,
				Principal:               r.Principal,
				Disabled:                r.Disabled,
				Managed:                 r.Managed,
				Internal:                r.Internal,
				OrgIdentifier:           r.OrgIdentifier,
				ProjectIdentifier:       r.ProjectIdentifier,
			},
		})
		if err != nil {
			failed = append(failed, fmt.Sprintln(r.Identifier, "-", err.Error()))
		}
		bar.Add(1)
	}
	bar.Finish()

	reportFailed(failed, "Roles:")
	return nil
}

func (api *ApiRequest) ListRoles(org, project string) ([]*model.RoleAssignmentContent, error) {

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     org,
			"projectIdentifier": project,
			"size":              "1000",
		}).
		Get(BaseURL + ROLES)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse(resp)
	}

	result := model.GetRoleResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	roles := []*model.RoleAssignmentContent{}
	for _, c := range result.Data.Content {
		roles = append(roles, &c.RoleAssignment)
	}

	return roles, nil
}

func (api *ApiRequest) CreateRoleAssignment(role *model.CreateRoleAssignment) error {

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(role).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
		}).
		Post(BaseURL + ROLES)
	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}

	return nil
}
