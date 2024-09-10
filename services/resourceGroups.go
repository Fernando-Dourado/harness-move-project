package services

import (
	"encoding/json"
	"fmt"

	"github.com/jf781/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
)

const RESOURCEGROUP = "/resourcegroup/api/v2/resourcegroup"

type ResourceGroupContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
}

func NewResourceGroupOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string) ResourceGroupContext {
	return ResourceGroupContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
	}
}

func (c ResourceGroupContext) Move() error {

	resourceGroups, err := c.api.listResourceGroups(c.sourceOrg, c.sourceProject)
	if err != nil {
		return err
	}

	bar := progressbar.Default(int64(len(resourceGroups)), "Resource Groups")
	var failed []string

	for _, rg := range resourceGroups {

		rg.OrgIdentifier = c.targetOrg
		rg.ProjectIdentifier = c.targetProject

		for i := range rg.IncludedScopes {
			rg.IncludedScopes[i].OrgIdentifier = &c.targetOrg
			rg.IncludedScopes[i].ProjectIdentifier = &c.targetProject
		}

		newResourceGroup := &model.NewResourceGroupContent{
			ResourceGroup: rg,
		}

		err = c.api.createResourceGroup(newResourceGroup)

		if err != nil {
			failed = append(failed, fmt.Sprintln(rg.Identifier, "-", err.Error()))
		}
		bar.Add(1)
	}
	bar.Finish()

	reportFailed(failed, "Resource Groups:")
	return nil
}

func (api *ApiRequest) listResourceGroups(org, project string) ([]*model.ResourceGroup, error) {

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
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse(resp)
	}

	result := model.GetResourceGroupResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
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

func (api *ApiRequest) createResourceGroup(rg *model.NewResourceGroupContent) error {

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
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}

	return nil
}
