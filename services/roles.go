package services

import (
	"encoding/json"

	"github.com/jf781/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
	"go.uber.org/zap"
)

const ROLE = "/authz/api/roles"

type RoleContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
	logger        *zap.Logger
}

func NewRoleOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string, logger *zap.Logger) RoleContext {
	return RoleContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
		logger:        logger,
	}
}

func (c RoleContext) Copy() error {

	c.logger.Info("Copying Roles",
		zap.String("project", c.sourceProject),
	)

	roles, err := c.api.listRoles(c.sourceOrg, c.sourceProject, c.logger)
	if err != nil {
		c.logger.Error("Failed to retrive roles",
			zap.String("Project", c.sourceProject),
			zap.Error(err),
		)
		return err
	}

	bar := progressbar.Default(int64(len(roles)), "Roles    ")

	for _, r := range roles {

		IncrementRolesTotal()

		c.logger.Info("Processing role",
			zap.String("role", r.Name),
			zap.String("targetProject", c.targetProject),
		)

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

		err = c.api.createRole(role, c.logger)

		if err != nil {
			c.logger.Error("Failed to create role",
				zap.String("role", role.Name),
				zap.Error(err),
			)
		} else {
			IncrementRolesMoved()
		}
		bar.Add(1)
	}
	bar.Finish()

	return nil
}

func (api *ApiRequest) listRoles(org, project string, logger *zap.Logger) ([]*model.ExistingRoles, error) {

	logger.Info("Fetching roles",
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
			"limit":             "100",
		}).
		Get(api.BaseURL + ROLE)
	if err != nil {
		logger.Error("Failed to request to list of roles",
			zap.Error(err),
		)
		return nil, err
	}
	if resp.IsError() {
		logger.Error("Error response from API when listing roles",
			zap.String("response",
				resp.String(),
			),
		)
		return nil, handleErrorResponse(resp)
	}

	result := model.GetRolesResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		logger.Error("Failed to parse response from API",
			zap.Error(err),
		)
		return nil, err
	}

	roles := []*model.ExistingRoles{}
	for _, c := range result.Data.Content {
		if !c.HarnessManaged {
			// Only add non-Harness managed roles
			newRole := c.Role
			roles = append(roles, &newRole)
		} else {
			logger.Warn("Skipping role because it is harness managed",
				zap.String("role", c.Role.Name),
				zap.Bool("harnessManaged", c.HarnessManaged),
			)
		}
	}

	return roles, nil
}

func (api *ApiRequest) createRole(role *model.NewRole, logger *zap.Logger) error {

	logger.Info("Creating role",
		zap.String("role", role.Name),
		zap.String("project", role.ProjectIdentifier),
	)

	IncrementApiCalls()

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(role).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     role.OrgIdentifier,
			"projectIdentifier": role.ProjectIdentifier,
		}).
		Post(api.BaseURL + ROLE)

	if err != nil {
		logger.Error("Failed to send request to create ",
			zap.String("Role", role.Name),
			zap.Error(err),
		)
		return err
	}
	if resp.IsError() {
		var errorResponse map[string]interface{}
		if err := json.Unmarshal(resp.Body(), &errorResponse); err == nil {
			if code, ok := errorResponse["code"].(string); ok && code == "DUPLICATE_FIELD" {
				// Log as a warning and skip the error
				logger.Info("Duplicate role found, ignoring error",
					zap.String("role", role.Name),
				)
				return nil
			}
		} else {
			logger.Error(
				"Error response from API when creating ",
				zap.String("Role", role.Name),
				zap.String("response",
					resp.String(),
				),
			)
		}
		return handleErrorResponse(resp)
	}

	return nil
}
