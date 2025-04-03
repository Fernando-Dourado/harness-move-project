package services

import (
	"encoding/json"
	"fmt"

	"github.com/Fernando-Dourado/harness-move-project/model"
	"github.com/go-resty/resty/v2"
)

const GET_PROJECT = "/ng/api/projects/{identifier}"

func (s *SourceRequest) ValidateSource(org, project string) error {
	return validateOrgProject(s.Client, s.Token, s.Account, org, project)
}

func (t *TargetRequest) ValidateTarget(org, project string) error {
	return validateOrgProject(t.Client, t.Token, t.Account, org, project)
}

func validateOrgProject(c *resty.Client, token, account, org, project string) error {
	resp, err := c.R().
		SetHeader("x-api-key", token).
		SetHeader("Content-Type", "application/json").
		SetPathParam("identifier", project).
		SetQueryParams(map[string]string{
			"accountIdentifier": account,
			"orgIdentifier":     org,
		}).
		Get(BaseURL + GET_PROJECT)
	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}
	result := model.GetProjectResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return err
	}

	if result.Data == nil {
		return fmt.Errorf("org %s or project %s not exist", org, project)
	}
	return nil
}
