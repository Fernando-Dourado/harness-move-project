package services

// import (
// 	"encoding/json"
// 	"fmt"

// 	"github.com/jf781/harness-move-project/model"
// 	"github.com/schollz/progressbar/v3"
// )

// const ROLE = "authz/api/roles"

// type RoleContext struct {
// 	api           *ApiRequest
// 	sourceOrg     string
// 	sourceProject string
// 	targetOrg     string
// 	targetProject string
// }

// func NewRoleOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string) RoleContext {
// 	return RoleContext{
// 		api:           api,
// 		sourceOrg:     sourceOrg,
// 		sourceProject: sourceProject,
// 		targetOrg:     targetOrg,
// 		targetProject: targetProject,
// 	}
// }

// func (c RoleContext) Move() error {

// 	roles, err := c.api.listRoles(c.sourceOrg, c.sourceProject)
// 	if err != nil {
// 		return err
// 	}

// 	bar := progressbar.Default(int64(len(roles)), "Roles")
// 	var failed []string

// 	for _, r := range roles {

// 		rolePrincipal := model.NewRoleAssignmentPrincipal{
// 			Identifier: r.Principal.Identifier,
// 			Type:       r.Principal.Type,
// 		}

// 		role := &model.NewRoleAssignment{
// 			ResourceGroupIdentifier: r.ResourceGroupIdentifier,
// 			RoleIdentifier:          r.RoleIdentifier,
// 			Principal:               rolePrincipal,
// 			OrgIdentifier:           c.targetOrg,
// 			ProjectIdentifier:       c.targetProject,
// 		}

// 		err = c.api.createRole(role)

// 		if err != nil {
// 			failed = append(failed, fmt.Sprintln(r.Identifier, "-", err.Error()))
// 		}
// 		bar.Add(1)
// 	}
// 	bar.Finish()

// 	reportFailed(failed, "Roles:")
// 	return nil
// }

// func (api *ApiRequest) listRoles(org, project string) ([]*model.RoleAssignmentListContent, error) {

// 	resp, err := api.Client.R().
// 		SetHeader("x-api-key", api.Token).
// 		SetHeader("Content-Type", "application/json").
// 		SetQueryParams(map[string]string{
// 			"accountIdentifier": api.Account,
// 			"orgIdentifier":     org,
// 			"projectIdentifier": project,
// 			"size":              "1000",
// 		}).
// 		Get(BaseURL + ROLEASSIGNMENT)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if resp.IsError() {
// 		return nil, handleErrorResponse(resp)
// 	}

// 	result := model.GetRoleResponse{}
// 	err = json.Unmarshal(resp.Body(), &result)
// 	if err != nil {
// 		return nil, err
// 	}

// 	roles := []*model.type RoleAssignmentContent struct {
// 		{}
// 	for _, c := range result.Data.Content {
// 		roles = append(roles, &c.Roles)
// 	}

// 	return roles, nil
// }

// func (api *ApiRequest) createRole(role *model.Roles) error {

// 	resp, err := api.Client.R().
// 		SetHeader("x-api-key", api.Token).
// 		SetHeader("Content-Type", "application/json").
// 		SetBody(role).
// 		SetQueryParams(map[string]string{
// 			"accountIdentifier": api.Account,
// 			"orgIdentifier":     role.OrgIdentifier,
// 			"projectIdentifier": role.ProjectIdentifier,
// 		}).
// 		Post(BaseURL + ROLEASSIGNMENT)

// 	if err != nil {
// 		return err
// 	}
// 	if resp.IsError() {
// 		return handleErrorResponse(resp)
// 	}

// 	return nil
// }
