package services

import (
	"encoding/json"
	"fmt"

	"github.com/Fernando-Dourado/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
)

const LIST_SERVICES = "/ng/api/servicesV2"
const CREATE_SERVICES = "/ng/api/servicesV2"

type ServiceContext struct {
	source        *SourceRequest
	target        *TargetRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
}

func NewServiceOperation(sourceApi *SourceRequest, targetApi *TargetRequest, sourceOrg, sourceProject, targetOrg, targetProject string) ServiceContext {
	return ServiceContext{
		source:        sourceApi,
		target:        targetApi,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
	}
}

func (c ServiceContext) Move() error {

	services, err := listServices(c.source, c.sourceOrg, c.sourceProject)
	if err != nil {
		return err
	}

	bar := progressbar.Default(int64(len(services)), "Services    ")
	var failed []string

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
		if err := createService(c.target, service); err != nil {
			failed = append(failed, fmt.Sprintln(s.Service.Name, "-", err.Error()))
		}
		bar.Add(1)
	}
	bar.Finish()

	reportFailed(failed, "services:")
	return nil
}

func listServices(s *SourceRequest, org, project string) ([]*model.ServiceListContent, error) {

	resp, err := s.Client.R().
		SetHeader("x-api-key", s.Token).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"accountIdentifier": s.Account,
			"orgIdentifier":     org,
			"projectIdentifier": project,
			"size":              "1000",
		}).
		Get(s.Url + LIST_SERVICES)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse(resp)
	}

	result := model.ServiceListResult{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	return result.Data.Content, nil
}

func createService(t *TargetRequest, service *model.CreateServiceRequest) error {

	resp, err := t.Client.R().
		SetHeader("x-api-key", t.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(service).
		SetQueryParams(map[string]string{
			"accountIdentifier": t.Account,
		}).
		Post(t.Url + CREATE_SERVICES)
	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}

	return nil
}
