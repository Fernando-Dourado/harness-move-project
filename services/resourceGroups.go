package services

import (
	"encoding/json"

	"github.com/jf781/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
	"go.uber.org/zap"
)

const RESOURCEGROUP = "/resourcegroup/api/v2/resourcegroup"

type ResourceGroupContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
	logger        *zap.Logger
}

func NewResourceGroupOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string, logger *zap.Logger) ResourceGroupContext {
	return ResourceGroupContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
		logger:        logger,
	}
}

func (c ResourceGroupContext) Copy() error {

	resourceGroups, err := c.api.listResourceGroups(c.sourceOrg, c.sourceProject, c.logger)
	if err != nil {
		c.logger.Error("Failed to retrive resource groups",
			zap.String("Project", c.sourceProject),
			zap.Error(err),
		)
		return err
	}

	bar := progressbar.Default(int64(len(resourceGroups)), "Resource Groups")

	for _, rg := range resourceGroups {

		c.logger.Info("Processing resource group",
			zap.String("resource group", rg.Name),
			zap.String("targetProject", c.targetProject),
		)

		rg.OrgIdentifier = c.targetOrg
		rg.ProjectIdentifier = c.targetProject

		for i := range rg.IncludedScopes {
			rg.IncludedScopes[i].OrgIdentifier = &c.targetOrg
			rg.IncludedScopes[i].ProjectIdentifier = &c.targetProject
		}

		newResourceGroup := &model.NewResourceGroupContent{
			ResourceGroup: rg,
		}

		err = c.api.createResourceGroup(newResourceGroup, c.logger)

		if err != nil {
			c.logger.Error("Failed to create resource group",
				zap.String("resource group", rg.Name),
				zap.Error(err),
			)
		}
		bar.Add(1)
	}
	bar.Finish()

	return nil
}

func (api *ApiRequest) listResourceGroups(org, project string, logger *zap.Logger) ([]*model.ResourceGroup, error) {

	logger.Info("Fetching resource groups",
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
		Get(api.BaseURL + RESOURCEGROUP)
	if err != nil {
		logger.Error("Failed to request to list of resource groups",
			zap.Error(err),
		)
		return nil, err
	}
	if resp.IsError() {
		logger.Error("Error response from API when listing resource groups",
			zap.String("response",
				resp.String(),
			),
		)
		return nil, handleErrorResponse(resp)
	}

	result := model.GetResourceGroupResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		logger.Error("Failed to parse response from API",
			zap.Error(err),
		)
		return nil, err
	}

	resourceGroups := []*model.ResourceGroup{}
	for _, c := range result.Data.Content {
		if !c.HarnessManaged {
			// Only add non-Harness managed Resoure Groups
			newResourceGroup := c.ResourceGroup
			resourceGroups = append(resourceGroups, &newResourceGroup)
		}
	}

	return resourceGroups, nil
}

func (api *ApiRequest) createResourceGroup(rg *model.NewResourceGroupContent, logger *zap.Logger) error {

	logger.Info("Creating resource group",
		zap.String("resource group", rg.ResourceGroup.Name),
		zap.String("project", rg.ResourceGroup.ProjectIdentifier),
	)

	IncrementApiCalls()

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(rg).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     rg.ResourceGroup.OrgIdentifier,
			"projectIdentifier": rg.ResourceGroup.ProjectIdentifier,
		}).
		Post(api.BaseURL + RESOURCEGROUP)

	if err != nil {
		logger.Error("Failed to send request to create ",
			zap.String("Resource group", rg.ResourceGroup.Name),
			zap.Error(err),
		)
		return err
	}
	if resp.IsError() {
		var errorResponse map[string]interface{}
		if err := json.Unmarshal(resp.Body(), &errorResponse); err == nil {
			if code, ok := errorResponse["code"].(string); ok && code == "DUPLICATE_FIELD" {
				// Log as a warning and skip the error
				logger.Info("Duplicate resource group found, ignoring error",
					zap.String("resource group", rg.ResourceGroup.Name),
				)
				return nil
			}
		} else {
			logger.Error(
				"Error response from API when creating ",
				zap.String("Resource group", rg.ResourceGroup.Name),
				zap.String("response",
					resp.String(),
				),
			)
		}
		return handleErrorResponse(resp)
	}

	return nil
}
