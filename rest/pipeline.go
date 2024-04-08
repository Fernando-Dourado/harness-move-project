package rest

import (
	"encoding/json"
	"fmt"

	"github.com/Fernando-Dourado/harness-move-project/model"
	"github.com/go-resty/resty/v2"
)

type PipelineRestContext struct {
	client  *resty.Client
	token   string
	account string
}

type PipelineRestClient interface {
	ListPipelines(org, project string) ([]*model.PipelineListContent, error)

	GetPipeline(org, project, pipeIdentifier string) (*model.PipelineGetData, error)

	CreatePipeline(org, project, yaml string) error
}

func NewPipelineClient(client *resty.Client, token, account string) *PipelineRestContext {
	return &PipelineRestContext{
		client:  client,
		token:   token,
		account: account,
	}
}

func (c PipelineRestContext) ListPipelines(org, project string) ([]*model.PipelineListContent, error) {

	resp, err := c.client.R().
		SetHeader("x-api-key", c.token).
		SetHeader("Content-Type", "application/json").
		SetBody(`{"filterType": "PipelineSetup"}`).
		SetQueryParams(map[string]string{
			"accountIdentifier": c.account,
			"orgIdentifier":     org,
			"projectIdentifier": project,
			"size":              "1000",
		}).
		Post(BaseURL_ + "/pipeline/api/pipelines/list")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse_(resp)
	}

	result := model.PipelineListResult{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	return result.Data.Content, nil
}

func (c PipelineRestContext) GetPipeline(org, project, pipeIdentifier string) (*model.PipelineGetData, error) {

	resp, err := c.client.R().
		SetHeader("x-api-key", c.token).
		SetHeader("Load-From-Cache", "false").
		SetQueryParams(map[string]string{
			"accountIdentifier": c.account,
			"orgIdentifier":     org,
			"projectIdentifier": project,
		}).
		Get(BaseURL_ + fmt.Sprintf("/pipeline/api/pipelines/%s", pipeIdentifier))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse_(resp)
	}

	result := model.PipelineGetResult{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (c PipelineRestContext) CreatePipeline(org, project, yaml string) error {
	return nil
}
