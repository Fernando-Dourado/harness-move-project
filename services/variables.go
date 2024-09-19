package services

import (
	"encoding/json"

	"github.com/jf781/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
	"go.uber.org/zap"
)

type VariableContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
	logger        *zap.Logger
}

func NewVariableOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string, logger *zap.Logger) VariableContext {
	return VariableContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
		logger:        logger,
	}
}

func (c VariableContext) Copy() error {

	c.logger.Info("Copying variables",
		zap.String("project", c.sourceProject),
	)

	variables, err := c.api.listVariables(c.sourceOrg, c.sourceProject, c.logger)
	if err != nil {
		c.logger.Error("Failed to retrive variables",
			zap.String("Project", c.sourceProject),
			zap.Error(err),
		)
		return err
	}

	bar := progressbar.Default(int64(len(variables)), "Variables")

	for _, v := range variables {

		IncrementVariablesTotal()

		c.logger.Info("Processing variable",
			zap.String("variable", v.Name),
			zap.String("targetProject", c.targetProject),
		)

		v.OrgIdentifier = c.targetOrg
		v.ProjectIdentifier = c.targetProject

		err = c.api.createVariable(&model.CreateVariableRequest{
			Variable: v,
		}, c.logger)
		if err != nil {
			c.logger.Error("Failed to create variable",
				zap.String("variable", v.Name),
				zap.Error(err),
			)
		} else {
			IncrementVariablesMoved()
		}
		bar.Add(1)
	}
	bar.Finish()

	return nil
}

func (api *ApiRequest) listVariables(org, project string, logger *zap.Logger) ([]*model.Variable, error) {

	logger.Info("Fetching variables",
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
			"size":              "1000",
		}).
		Get(api.BaseURL + "/ng/api/variables")
	if err != nil {
		logger.Error("Failed to request to list of variables",
			zap.Error(err),
		)
		return nil, err
	}
	if resp.IsError() {
		logger.Error("Error response from API when listing variables",
			zap.String("response",
				resp.String(),
			),
		)
		return nil, handleErrorResponse(resp)
	}

	result := model.GetVariablesResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		logger.Error("Failed to parse response from API",
			zap.Error(err),
		)
		return nil, err
	}

	variables := []*model.Variable{}
	for _, c := range result.Data.Content {
		variables = append(variables, c.Variable)
	}

	return variables, nil
}

func (api *ApiRequest) createVariable(variable *model.CreateVariableRequest, logger *zap.Logger) error {

	logger.Info("Creating variable",
		zap.String("variable", variable.Variable.Name),
		zap.String("project", variable.Variable.ProjectIdentifier),
	)

	IncrementApiCalls()

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(variable).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
		}).
		Post(api.BaseURL + "/ng/api/variables")
	if err != nil {
		logger.Error("Failed to send request to create ",
			zap.String("variable", variable.Variable.Name),
			zap.Error(err),
		)
		return err
	}
	if resp.IsError() {
		var errorResponse map[string]interface{}
		if err := json.Unmarshal(resp.Body(), &errorResponse); err == nil {
			if code, ok := errorResponse["code"].(string); ok && code == "DUPLICATE_FIELD" {
				// Log as a warning and skip the error
				logger.Info("Duplicate variable found, ignoring error",
					zap.String("variable", variable.Variable.Name),
				)
				return nil
			}
		} else {
			logger.Error(
				"Error response from API when creating ",
				zap.String("variable", variable.Variable.Name),
				zap.String("response",
					resp.String(),
				),
			)
		}
		return handleErrorResponse(resp)
	}

	return nil
}
