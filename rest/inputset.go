package rest

import (
	"github.com/Fernando-Dourado/harness-move-project/model"
	"github.com/go-resty/resty/v2"
)

type InputsetRestContext struct {
	api     *resty.Client
	token   string
	account string
}

type InputsetRestClient interface {
	ListInputsets(org, project, pipelineIdentifier string) ([]*model.ListInputsetContent, error)

	GetInputset(org, project, pipelineIdentifier, isIdentifier string) (*model.GetInputsetData, error)

	CreateInputset(org, project, pipelineIdentifier, yaml string) error
}

func NewInputsetClient(client *resty.Client, token, account string) *InputsetRestContext {
	return &InputsetRestContext{
		api:     client,
		token:   token,
		account: account,
	}
}

func (api *InputsetRestContext) ListInputsets(org, project, pipelineIdentifier string) ([]*model.ListInputsetContent, error) {
	return nil, nil
	// resp, err := api.Client.R().
	// 	SetHeader("x-api-key", api.Token).
	// 	SetHeader("Content-Type", "application/json").
	// 	SetQueryParams(map[string]string{
	// 		"accountIdentifier":  api.Account,
	// 		"orgIdentifier":      org,
	// 		"projectIdentifier":  project,
	// 		"pipelineIdentifier": pipelineIdentifier,
	// 		"inputSetType":       "ALL",
	// 		"size":               "1000",
	// 	}).
	// 	Get(BaseURL + "/pipeline/api/inputSets")
	// if err != nil {
	// 	return nil, err
	// }
	// if resp.IsError() {
	// 	return nil, handleErrorResponse(resp)
	// }

	// result := &model.ListInputsetResponse{}
	// err = json.Unmarshal(resp.Body(), &result)
	// if err != nil {
	// 	return nil, err
	// }

	// return result.Data.Content, nil
}

func (r *InputsetRestContext) GetInputset(org, project, pipelineIdentifier, isIdentifier string) (*model.GetInputsetData, error) {
	return nil, nil
	// resp, err := api.Client.R().
	// 	SetHeader("x-api-key", api.Token).
	// 	SetHeader("Content-Type", "application/json").
	// 	SetHeader("Load-From-Cache", "false").
	// 	SetPathParam("inputset", isIdentifier).
	// 	SetQueryParams(map[string]string{
	// 		"accountIdentifier":  api.Account,
	// 		"orgIdentifier":      org,
	// 		"projectIdentifier":  project,
	// 		"pipelineIdentifier": pipelineIdentifier,
	// 	}).
	// 	Get(BaseURL + "/pipeline/api/inputSets/{inputset}")
	// if err != nil {
	// 	return nil, err
	// }
	// if resp.IsError() {
	// 	return nil, handleErrorResponse(resp)
	// }

	// result := &model.GetInputsetResponse{}
	// if err = json.Unmarshal(resp.Body(), &result); err != nil {
	// 	return nil, err
	// }

	// return result.Data, nil
}

func (r InputsetRestContext) CreateInputset(org, project, pipelineIdentifier, yaml string) error {

	// resp, err := api.Client.R().
	// 	SetHeader("x-api-key", api.Token).
	// 	SetHeader("Content-Type", "application/yaml").
	// 	SetBody(yaml).
	// 	SetQueryParams(map[string]string{
	// 		"accountIdentifier":  api.Account,
	// 		"orgIdentifier":      org,
	// 		"projectIdentifier":  project,
	// 		"pipelineIdentifier": pipelineIdentifier,
	// 	}).
	// 	Post(BaseURL + "/pipeline/api/inputSets")
	// if err != nil {
	// 	return err
	// }
	// if resp.IsError() {
	// 	return handleErrorResponse(resp)
	// }

	return nil
}
