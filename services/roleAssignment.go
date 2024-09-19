package services

import (
	"encoding/json"

	"github.com/jf781/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
	"go.uber.org/zap"
)

const ROLEASSIGNMENT = "/authz/api/roleassignments"

type RoleAssignmentContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
	logger        *zap.Logger
}

func NewRoleAssignmentOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string, logger *zap.Logger) RoleAssignmentContext {
	return RoleAssignmentContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
		logger:        logger,
	}
}

func (c RoleAssignmentContext) Copy() error {

	c.logger.Info("Copying role assignments",
		zap.String("project", c.sourceProject),
	)

	roleAssignments, err := c.api.listRoleAssignments(c.sourceOrg, c.sourceProject, c.logger)
	if err != nil {
		c.logger.Error("Failed to retrive role assignments",
			zap.String("Project", c.sourceProject),
			zap.Error(err),
		)
		return err
	}

	bar := progressbar.Default(int64(len(roleAssignments)), "Roles Assignments    ")

	for _, r := range roleAssignments {

		IncrementRoleAssignmentsTotal()

		c.logger.Info("Processing role assignment",
			zap.String("role assignment", r.RoleIdentifier),
			zap.String("targetProject", c.targetProject),
		)

		role := &model.NewRoleAssignment{
			Identifier:              r.Identifier,
			ResourceGroupIdentifier: r.ResourceGroupIdentifier,
			RoleIdentifier:          r.RoleIdentifier,
			Principal:               r.Principal,
			OrgIdentifier:           c.targetOrg,
			ProjectIdentifier:       c.targetProject,
		}

		err = c.api.createRoleAssignment(role, c.logger)

		if err != nil {
			c.logger.Error("Failed to create role assignment",
				zap.String("role assignment", r.Identifier),
				zap.Error(err),
			)
		} else {
			IncrementRoleAssignmentsMoved()
		}
		bar.Add(1)
	}
	bar.Finish()

	return nil
}

func (api *ApiRequest) listRoleAssignments(org, project string, logger *zap.Logger) ([]*model.ExistingRoleAssignment, error) {

	logger.Info("Fetching role assignments",
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
		Get(api.BaseURL + ROLEASSIGNMENT)
	if err != nil {
		logger.Error("Failed to request to list of role assignments",
			zap.Error(err),
		)
		return nil, err
	}
	if resp.IsError() {
		logger.Error("Error response from API when listing role assignments",
			zap.String("response",
				resp.String(),
			),
		)
		return nil, handleErrorResponse(resp)
	}

	result := model.GetRoleAssignmentResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		logger.Error("Failed to parse response from API",
			zap.Error(err),
		)
		return nil, err
	}

	roleAssignments := []*model.ExistingRoleAssignment{}
	for _, c := range result.Data.Content {
		newRoleAssignment := c.RoleAssignment
		roleAssignments = append(roleAssignments, &newRoleAssignment)
	}

	return roleAssignments, nil
}

func (api *ApiRequest) createRoleAssignment(role *model.NewRoleAssignment, logger *zap.Logger) error {

	logger.Info("Creating role assignment",
		zap.String("role assignment", role.RoleIdentifier),
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
		Post(api.BaseURL + ROLEASSIGNMENT)

	if err != nil {
		logger.Error("Failed to send request to create ",
			zap.String("Role assignment", role.Identifier),
			zap.Error(err),
		)
		return err
	}
	if resp.IsError() {
		var errorResponse map[string]interface{}
		if err := json.Unmarshal(resp.Body(), &errorResponse); err == nil {
			if code, ok := errorResponse["code"].(string); ok && code == "DUPLICATE_FIELD" {
				// Log as a warning and skip the error
				logger.Info("Duplicate role assignment found, ignoring error",
					zap.String("role assignment", role.Identifier),
				)
				return nil
			}
		} else {
			logger.Error(
				"Error response from API when creating ",
				zap.String("Role assignment", role.Identifier),
				zap.String("response",
					resp.String(),
				),
			)
		}
		return handleErrorResponse(resp)
	}

	return nil
}
