package	services

import (
	"encoding/json"
	// "fmt"

	"github.com/jf781/harness-move-project/model"
	// "github.com/schollz/progressbar/v3"
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

	return roles
}

func (api *ApiRequest) ListRoles(org, project string) ([]*model.RoleListContent, error) {

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     org,
			"projectIdentifier": project,
			"size":              "1000",
		}).
		Get(BaseURL + LIST_SERVICES)
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

	return result.Data.Content, nil
}