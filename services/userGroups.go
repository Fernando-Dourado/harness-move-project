package services

import (
	"encoding/json"

	"github.com/jf781/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
	"go.uber.org/zap"
)

const USERGROUP = "/ng/api/user-groups"
const USERLOOKUP = "/ng/api/user/aggregate/"

type UserGroupContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
	logger        *zap.Logger
}

func NewUserGroupOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string, logger *zap.Logger) UserGroupContext {
	return UserGroupContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
		logger:        logger,
	}
}

func (c UserGroupContext) Copy() error {

	c.logger.Info("Copying user group",
		zap.String("project", c.sourceProject),
	)

	groups, err := c.api.listUserGroups(c.sourceOrg, c.sourceProject, c.logger)
	if err != nil {
		c.logger.Error("Failed to retrive user group",
			zap.String("Project", c.sourceProject),
			zap.Error(err),
		)
		return err
	}

	bar := progressbar.Default(int64(len(groups)), "User Groups    ")

	for _, g := range groups {

		c.logger.Info("Processing user group",
			zap.String("user group", g.Name),
			zap.String("sourceProject", c.sourceProject),
		)

		for i := range g.Users {
			user := &model.UserGroupLookup{
				Identifier:        g.Users[i],
				OrgIdentifier:     c.sourceOrg,
				ProjectIdentifier: c.sourceProject,
			}

			userEmail, err := c.api.getUsersEmail(user, c.logger)

			if err != nil {
				c.logger.Error("Failed to get user email",
					zap.String("user group", user.Identifier),
					zap.Error(err),
				)
			}

			g.Users = append(g.Users, userEmail.EmailAddress)
		}

		g.OrgIdentifier = c.targetOrg
		g.ProjectIdentifier = c.targetProject

		err = c.api.addUserGroup(g, c.logger)

		if err != nil {
			c.logger.Error("Failed to create group",
				zap.String("group", g.Name),
				zap.Error(err),
			)

		}
		bar.Add(1)
	}
	bar.Finish()

	return nil
}

func (api *ApiRequest) listUserGroups(org, project string, logger *zap.Logger) ([]*model.UserGroup, error) {

	logger.Info("Fetching user groups",
		zap.String("org", org),
		zap.String("project", project),
	)

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
		logger.Error("Failed to request to list of user groups",
			zap.Error(err),
		)
		return nil, err
	}
	if resp.IsError() {

		logger.Error("Error response from API when listing user groups",
			zap.String("response",
				resp.String(),
			),
		)
		return nil, handleErrorResponse(resp)
	}

	result := model.GetUserGroupsResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		logger.Error("Failed to parse response from API",
			zap.Error(err),
		)

		return nil, err
	}

	userGroups := []*model.UserGroup{}
	for _, c := range result.Data.Content {
		if !c.HarnessManaged {
			userGroups = append(userGroups, c)
		} else {
			logger.Warn("Skipping user group because it is Harness managed",
				zap.String("user group", c.Name),
				zap.Bool("harnessManaged", c.HarnessManaged),
			)
		}
	}

	return userGroups, nil
}

func (api *ApiRequest) getUsersEmail(user *model.UserGroupLookup, logger *zap.Logger) (*model.UserGroupEmail, error) {

	logger.Info("Fetching user details",
		zap.String("user", user.Identifier),
	)

	IncrementApiCalls()

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
		logger.Error("Failed to request to list of user groups",
			zap.Error(err),
		)
		return nil, err
	}
	if resp.IsError() {
		logger.Error("Error response from API when listing user groups",
			zap.String("response",
				resp.String(),
			),
		)
		return nil, handleErrorResponse(resp)
	}

	result := model.UserGroupEmail{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		logger.Error("Failed to parse response from API",
			zap.Error(err),
		)
		return nil, err
	}

	return &result, nil
}

func (api *ApiRequest) addUserGroup(userGroup *model.UserGroup, logger *zap.Logger) error {

	logger.Info("Creating user group",
		zap.String("user group", userGroup.Name),
		zap.String("project", userGroup.ProjectIdentifier),
	)

	IncrementApiCalls()

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
		logger.Error("Failed to send request to create ",
			zap.String("user group", userGroup.Name),
			zap.Error(err),
		)
		return err
	}
	if resp.IsError() {
		var errorResponse map[string]interface{}
		if err := json.Unmarshal(resp.Body(), &errorResponse); err == nil {
			if code, ok := errorResponse["code"].(string); ok && code == "DUPLICATE_FIELD" {
				// Log as a warning and skip the error
				logger.Info("Duplicate user group found, ignoring error",
					zap.String("user group", userGroup.Name),
				)
				return nil
			}
		} else {
			logger.Error(
				"Error response from API when creating ",
				zap.String("user group", userGroup.Name),
				zap.String("response",
					resp.String(),
				),
			)
		}
		return handleErrorResponse(resp)
	}

	return nil
}
