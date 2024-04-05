package services

import (
	"encoding/json"
	"fmt"

	"github.com/jf781/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
)

const ROLE = "/authz/api/roles"

type RoleContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
}

func NewRoleOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string) RoleContext {
	return RoleContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
	}
}

func (c RoleContext) Move() error {

	roles, err := c.api.listRoles(c.sourceOrg, c.sourceProject)
	if err != nil {
		return err
	}

	bar := progressbar.Default(int64(len(roles)), "Roles")
	var failed []string

	for _, r := range roles {

		role := &model.NewRole{
			Identifier:         r.Identifier,
			Name:               r.Name,
			Description:        r.Description,
			Tags:               r.Tags,
			Permissions:        r.Permissions,
			AllowedScopeLevels: r.AllowedScopeLevels,
			OrgIdentifier:      c.targetOrg,
			ProjectIdentifier:  c.targetProject,
		}

		err = c.api.createRole(role)

		if err != nil {
			failed = append(failed, fmt.Sprintln(r.Identifier, "-", err.Error()))
		}
		bar.Add(1)
	}
	bar.Finish()

	reportFailed(failed, "Roles:")
	return nil
}

func (api *ApiRequest) listRoles(org, project string) ([]*model.ExistingRoles, error) {

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     org,
			"projectIdentifier": project,
			"size":              "1000",
		}).
		Get(BaseURL + ROLE)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse(resp)
	}

	result := model.GetRolesResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	roles := []*model.ExistingRoles{}
	for _, c := range result.Data.Content {
		if !c.HarnessManaged {
			// Only add non-Harness managed roles
			roles = append(roles, &c.Role)
		}
	}

	return roles, nil
}

func (api *ApiRequest) createRole(role *model.NewRole) error {

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(role).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     role.OrgIdentifier,
			"projectIdentifier": role.ProjectIdentifier,
		}).
		Post(BaseURL + ROLE)

	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}

	return nil
}
