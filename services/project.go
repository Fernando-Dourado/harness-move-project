package services

import (
	"encoding/json"
	"fmt"

	"github.com/jf781/harness-move-project/model"
)

const GET_PROJECT = "/ng/api/projects/{identifier}"

func (api *ApiRequest) ValidateProject(org, project string) error {

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetPathParam("identifier", project).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
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
