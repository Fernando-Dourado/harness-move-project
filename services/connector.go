package services

import (
	"encoding/json"
	"fmt"

	"github.com/jf781/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
)

const CONNECTORLOOKUP = "/ng/api/connectors/listV2"
const CONNECTORCREATE = "/ng/api/connectors"

type ConnectorContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
}

func NewConnectorOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string) ConnectorContext {
	return ConnectorContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
	}
}

func (c ConnectorContext) Move() error {

	connectors, err := c.api.listConnectors(c.sourceOrg, c.sourceProject)
	if err != nil {
		return err
	}

	bar := progressbar.Default(int64(len(connectors)), "Connectors    ")
	var failed []string

	for _, cn := range connectors {

		cn.Connector.OrgIdentifier = c.targetOrg
		cn.Connector.ProjectIdentifier = c.targetProject

		err = c.api.addConnector(cn)

		if err != nil {
			failed = append(failed, fmt.Sprintln(cn.Connector.Name, "-", err.Error()))
		}
		bar.Add(1)
	}
	bar.Finish()

	reportFailed(failed, "Connectors:")
	return nil
}

func (api *ApiRequest) listConnectors(org, project string) ([]*model.ConnectorContent, error) {

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     org,
			"projectIdentifier": project,
			"pageSize":          "100",
		}).
		Post(api.BaseURL + CONNECTORLOOKUP)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse(resp)
	}

	result := model.ConnectorListResult{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	connectors := []*model.ConnectorContent{}
	for _, c := range result.Data.Content {
		if !c.HarnessManaged && c.Status.Status == "SUCCESS" {
			newConnectors := c
			connectors = append(connectors, &newConnectors)
		} else {
			fmt.Println("Skipping connector - Harness managed or Status is inactive: ", c.Connector.Name)
		}
	}

	return connectors, nil
}

func (api *ApiRequest) addConnector(connector *model.ConnectorContent) error {

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(connector).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
		}).
		Post(api.BaseURL + CONNECTORCREATE)

	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}

	return nil
}
