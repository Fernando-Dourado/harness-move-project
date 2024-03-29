package services

import (
	"encoding/json"
	"fmt"

	"github.com/Fernando-Dourado/harness-move-project/model"
)

const LIST_SERVICES = "/ng/api/servicesV2"
const CREATE_SERVICES = "/ng/api/servicesV2"

type ServiceContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
}

func NewServiceOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string) ServiceContext {
	return ServiceContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
	}
}

func (c ServiceContext) Move() error {

	services, err := c.api.listServices(c.sourceOrg, c.sourceProject)
	if err != nil {
		return err
	}

	for _, s := range services {
		newYaml := createYaml(s.Service.Yaml, c.sourceOrg, c.sourceProject, c.targetOrg, c.targetProject)
		service := &model.CreateServiceRequest{
			OrgIdentifier:     c.targetOrg,
			ProjectIdentifier: c.targetProject,
			Identifier:        s.Service.Identifier,
			Name:              s.Service.Name,
			Description:       s.Service.Description,
			Yaml:              newYaml,
		}
		if err := c.api.createService(service); err != nil {
			return err
		}
	}

	return nil
}

func (api *ApiRequest) listServices(org, project string) ([]*model.ServiceListContent, error) {

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     org,
			"projectIdentifier": project,
			"size":              "1000",
		}).
		Get(BaseURL + LIST_SERVICES)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf(resp.Status())
	}

	result := model.ServiceListResult{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	return result.Data.Content, nil
}

func (api *ApiRequest) createService(service *model.CreateServiceRequest) error {

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(service).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
		}).
		Post(BaseURL + CREATE_SERVICES)
	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleCreateErrorResponse(resp)
	}

	return nil
}
