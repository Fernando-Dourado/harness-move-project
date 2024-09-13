package services

import (
	"encoding/json"

	"github.com/jf781/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
	"go.uber.org/zap"
)

const SERVICEACCOUNTS = "/ng/api/serviceaccount"

type ServiceAccountContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
	logger        *zap.Logger
}

func NewServiceAccountOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string, logger *zap.Logger) ServiceAccountContext {
	return ServiceAccountContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
		logger:        logger,
	}
}

func (c ServiceAccountContext) Copy() error {

	c.logger.Info("Copying Service accounts",
		zap.String("project", c.sourceProject),
	)

	serviceAccounts, err := c.api.listServiceAccounts(c.sourceOrg, c.sourceProject, c.logger)
	if err != nil {
		c.logger.Error("Failed to retrive service accounts",
			zap.String("Project", c.sourceProject),
			zap.Error(err),
		)
		return err
	}

	bar := progressbar.Default(int64(len(serviceAccounts)), "Service Accounts    ")

	for _, sa := range serviceAccounts {

		c.logger.Info("Processing service account",
			zap.String("service account", sa.Email),
			zap.String("targetProject", c.targetProject),
		)

		sa.OrgIdentifier = c.targetOrg
		sa.ProjectIdentifier = c.targetProject

		err = c.api.createServiceAccount(sa, c.logger)

		if err != nil {
			c.logger.Error("Failed to create service account",
				zap.String("service account", sa.Email),
				zap.Error(err),
			)
		}
		bar.Add(1)
	}
	bar.Finish()

	return nil
}

func (api *ApiRequest) listServiceAccounts(org, project string, logger *zap.Logger) ([]*model.GetServiceAccountData, error) {

	logger.Info("Fetching service accounts",
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
		}).
		Get(api.BaseURL + SERVICEACCOUNTS)
	if err != nil {
		logger.Error("Failed to request to list of service accounts",
			zap.Error(err),
		)
		return nil, err
	}
	if resp.IsError() {
		logger.Error("Error response from API when listing service accounts",
			zap.String("response",
				resp.String(),
			),
		)
		return nil, handleErrorResponse(resp)
	}

	result := model.GetServiceACcountResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		logger.Error("Failed to parse response from API",
			zap.Error(err),
		)
		return nil, err
	}

	return result.Data, nil
}

func (api *ApiRequest) createServiceAccount(serviceAccount *model.GetServiceAccountData, logger *zap.Logger) error {

	logger.Info("Creating service account",
		zap.String("service account", serviceAccount.Email),
		zap.String("project", serviceAccount.ProjectIdentifier),
	)

	IncrementApiCalls()

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(serviceAccount).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     serviceAccount.OrgIdentifier,
			"projectIdentifier": serviceAccount.ProjectIdentifier,
		}).
		Post(api.BaseURL + SERVICEACCOUNTS)

	if err != nil {
		logger.Error("Failed to send request to create ",
			zap.String("Service account", serviceAccount.Email),
			zap.Error(err),
		)
		return err
	}
	if resp.IsError() {
		var errorResponse map[string]interface{}
		if err := json.Unmarshal(resp.Body(), &errorResponse); err == nil {
			if code, ok := errorResponse["code"].(string); ok && code == "DUPLICATE_FIELD" {
				// Log as a warning and skip the error
				logger.Info("Duplicate service account found, ignoring error",
					zap.String("service account", serviceAccount.Name),
				)
				return nil
			}
		} else {
			logger.Error(
				"Error response from API when creating ",
				zap.String("Service account", serviceAccount.Name),
				zap.String("response",
					resp.String(),
				),
			)
		}
		return handleErrorResponse(resp)
	}

	return nil
}
