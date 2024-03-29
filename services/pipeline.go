package services

import (
	"encoding/json"

	"github.com/Fernando-Dourado/harness-move-project/model"
)

const LIST_PIPELINES = "/pipeline/api/pipelines/list"

func (api *ApiRequest) ListPipelines(org, project string) ([]*model.PipelineListContent, error) {

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(`{"filterType": "PipelineSetup"}`).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     org,
			"projectIdentifier": project,
			"size":              "1000",
		}).
		Post(BaseURL + LIST_PIPELINES)
	if err != nil {
		return nil, err
	}

	result := model.PipelineListResult{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	return result.Data.Content, nil
}
