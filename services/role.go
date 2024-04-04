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
		if r.Principal.ScopeLevel != nil {
			// Skip roles that are not project level
			continue
		}

		rolePrincipal := model.CreateRoleAssignmentPrincipal{
			Identifier: r.Principal.Identifier,
			Type:       r.Principal.Type,
		}

		role := &model.CreateRoleAssignment{
			ResourceGroupIdentifier: r.ResourceGroupIdentifier,
			RoleIdentifier:          r.RoleIdentifier,
			Principal:               rolePrincipal,
			OrgIdentifier:           c.targetOrg,
			ProjectIdentifier:       c.targetProject,
		}

		fmt.Printf("Role: %+v\n", role)
		
		err = c.api.CreateRoleAssignment(role)

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
	//fmt.Println("Role Body:", role)
	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}

	return nil
}
