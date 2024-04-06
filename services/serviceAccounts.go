package services

import (
	"encoding/json"
	"fmt"

	"github.com/jf781/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
)

const SERVICEACCOUNTS = "/ng/api/serviceaccount"

type ServiceAccountContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
}

func NewServiceAccountOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string) ServiceAccountContext {
	return ServiceAccountContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
	}
}

func (c ServiceAccountContext) Move() error {

	serviceAccounts, err := c.api.listServiceAccounts(c.sourceOrg, c.sourceProject)
	if err != nil {
		return err
	}

	bar := progressbar.Default(int64(len(serviceAccounts)), "Service Accounts    ")
	var failed []string

	for _, sa := range serviceAccounts {

		sa.OrgIdentifier = c.targetOrg
		sa.ProjectIdentifier = c.targetProject

		err = c.api.createServiceAccount(sa)

		if err != nil {
			failed = append(failed, fmt.Sprintln(sa.Name, "-", err.Error()))
		}
		bar.Add(1)
	}
	bar.Finish()

	reportFailed(failed, "Service Accounts:")
	return nil
}

func (api *ApiRequest) listServiceAccounts(org, project string) ([]*model.GetServiceAccountData, error) {

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     org,
			"projectIdentifier": project,
		}).
		Get(BaseURL + SERVICEACCOUNTS)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse(resp)
	}

	result := model.GetServiceACcountResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (api *ApiRequest) createServiceAccount(servcieAccount *model.GetServiceAccountData) error {

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(servcieAccount).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     servcieAccount.OrgIdentifier,
			"projectIdentifier": servcieAccount.ProjectIdentifier,
		}).
		Post(BaseURL + SERVICEACCOUNTS)

	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}

	return nil
}
