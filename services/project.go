package services

import (
	"encoding/json"
	"fmt"

	"github.com/Fernando-Dourado/harness-move-project/model"
)

const LIST_PROJECTS = "/ng/api/projects/list"

func (api *ApiRequest) ValidateProject(org, project string) error {

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"accountIdentifier":  api.Account,
			"orgIdentifiers":     org,
			"projectIdentifiers": project,
			"pageIndex":          "0",
			"pageSize":           "1",
		}).
		Get(BaseURL + LIST_PROJECTS)
	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}
	result := model.ProjectExistResult{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return err
	}

	if result.Data.PageItemCount != 1 || result.Data.Empty {
		return fmt.Errorf("org %s or project %s not exist", org, project)
	}
	return nil
}
