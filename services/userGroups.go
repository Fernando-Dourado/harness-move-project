package services

import (
	"encoding/json"
	"fmt"

	"github.com/jf781/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
)

const USERGROUP = "/ng/api/user-groups"
const USERLOOKUP = "/ng/api/user/aggregate/"

type UserGroupContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
}

func NewUserGroupOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string) UserGroupContext {
	return UserGroupContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
	}
}

func (c UserGroupContext) Move() error {

	groups, err := c.api.listUserGroups(c.sourceOrg, c.sourceProject)
	if err != nil {
		return err
	}

	bar := progressbar.Default(int64(len(groups)), "User Groups    ")
	var failed []string

	for _, g := range groups {

		for i := range g.Users {
			user := &model.UserGroupLookup{
				Identifier:        g.Users[i],
				OrgIdentifier:     c.sourceOrg,
				ProjectIdentifier: c.sourceProject,
			}

			userEmail, err := c.api.getUsersEmail(user)

			if err != nil {
				failed = append(failed, fmt.Sprintln(userEmail.Name, "-", err.Error()))
			}

			g.Users = append(g.Users, userEmail.EmailAddress)
		}

		g.OrgIdentifier = c.targetOrg
		g.ProjectIdentifier = c.targetProject

		err = c.api.addUserGroup(g)

		if err != nil {
			failed = append(failed, fmt.Sprintln(g.Name, "-", err.Error()))
		}
		bar.Add(1)
	}
	bar.Finish()

	reportFailed(failed, "User Groups:")
	return nil
}

func (api *ApiRequest) listUserGroups(org, project string) ([]*model.UserGroup, error) {

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     org,
			"projectIdentifier": project,
			"pageSize":          "100",
		}).
		Get(api.BaseURL + USERGROUP)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse(resp)
	}

	result := model.GetUserGroupsResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	userGroups := []*model.UserGroup{}
	for _, c := range result.Data.Content {
		if !c.HarnessManaged {
			userGroups = append(userGroups, c)
		}
	}

	return userGroups, nil
}

func (api *ApiRequest) getUsersEmail(user *model.UserGroupLookup) (*model.UserGroupEmail, error) {

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(user).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     user.OrgIdentifier,
			"projectIdentifier": user.ProjectIdentifier,
		}).
		Get(api.BaseURL + USERLOOKUP + "/" + user.Identifier)

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse(resp)
	}

	result := model.UserGroupEmail{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (api *ApiRequest) addUserGroup(userGroup *model.UserGroup) error {

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(userGroup).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     userGroup.OrgIdentifier,
			"projectIdentifier": userGroup.ProjectIdentifier,
		}).
		Post(api.BaseURL + USERGROUP)

	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}

	return nil
}
