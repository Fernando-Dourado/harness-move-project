package services

import (
	"encoding/json"

	"github.com/jf781/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
	"go.uber.org/zap"
)

const LISTUSER = "/ng/api/user/aggregate"
const ADDUSER = "/ng/api/user/users"

type UserScopeContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
	logger        *zap.Logger
}

func NewUserScopeOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string, logger *zap.Logger) UserScopeContext {
	return UserScopeContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
		logger:        logger,
	}
}

func (c UserScopeContext) Copy() error {

	c.logger.Info("Copying users",
		zap.String("project", c.sourceProject),
	)

	users, err := c.api.listUsers(c.sourceOrg, c.sourceProject, c.logger)
	if err != nil {
		c.logger.Error("Failed to retrive users",
			zap.String("Project", c.sourceProject),
			zap.Error(err),
		)
		return err
	}

	bar := progressbar.Default(int64(len(users)), "Users    ")

	for _, u := range users {

		IncrementUsersTotal()

		c.logger.Info("Processing user",
			zap.String("user", u.Name),
			zap.String("sourceProject", c.sourceProject),
		)

		userToAdd := &model.UserEmail{
			EmailAddress:      []string{u.Email},
			OrgIdentifier:     c.targetOrg,
			ProjectIdentifier: c.targetProject,
		}

		err = c.api.addUserToScope(userToAdd, c.logger)

		if err != nil {
			c.logger.Error("Failed to create user",
				zap.String("user", u.Name),
				zap.Error(err),
			)
		} else {
			IncrementUsersMoved()
		}
		bar.Add(1)
	}
	bar.Finish()

	return nil
}

func (api *ApiRequest) listUsers(org, project string, logger *zap.Logger) ([]*model.User, error) {

	logger.Info("Fetching users",
		zap.String("org", org),
		zap.String("project", project),
	)

	IncrementApiCalls()

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     org,
			"projectIdentifier": project,
			"pageSize":          "100",
		}).
		Post(api.BaseURL + LISTUSER)
	if err != nil {
		logger.Error("Failed to request to list of users",
			zap.Error(err),
		)
		return nil, err
	}
	if resp.IsError() {
		logger.Error("Error response from API when listing users",
			zap.String("response",
				resp.String(),
			),
		)
		return nil, handleErrorResponse(resp)
	}

	result := model.GetUserResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		logger.Error("Failed to parse response from API",
			zap.Error(err),
		)
		return nil, err
	}

	users := []*model.User{}
	for _, c := range result.Data.Content {
		newUser := c.User
		users = append(users, &newUser)
	}

	return users, nil
}

func (api *ApiRequest) addUserToScope(user *model.UserEmail, logger *zap.Logger) error {

	logger.Info("Creating user",
		zap.String("user", user.EmailAddress[0]),
		zap.String("project", user.ProjectIdentifier),
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
		Post(api.BaseURL + ADDUSER)

	if err != nil {
		logger.Error("Failed to send request to create ",
			zap.String("user", user.EmailAddress[0]),
			zap.Error(err),
		)
		return err
	}
	if resp.IsError() {
		var errorResponse map[string]interface{}
		if err := json.Unmarshal(resp.Body(), &errorResponse); err == nil {
			if code, ok := errorResponse["code"].(string); ok && code == "DUPLICATE_FIELD" {
				// Log as a warning and skip the error
				logger.Info("Duplicate user found, ignoring error",
					zap.String("user", user.EmailAddress[0]),
				)
				return nil
			}
		} else {
			logger.Error(
				"Error response from API when creating ",
				zap.String("user", user.EmailAddress[0]),
				zap.String("response",
					resp.String(),
				),
			)
		}
		return handleErrorResponse(resp)
	}

	return nil
}
