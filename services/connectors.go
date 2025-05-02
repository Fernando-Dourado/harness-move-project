package services

import (
	"encoding/json"
	"fmt"

	"github.com/Fernando-Dourado/harness-move-project/model"
	"github.com/harness/harness-go-sdk/harness/nextgen"
	"github.com/schollz/progressbar/v3"
)

type ConnectorContext struct {
	source        *SourceRequest
	target        *TargetRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
}

func NewConnectorOperation(sourceApi *SourceRequest, targetApi *TargetRequest, st *SourceTarget) ConnectorContext {
	return ConnectorContext{
		source:        sourceApi,
		target:        targetApi,
		sourceOrg:     st.SourceOrg,
		sourceProject: st.SourceProject,
		targetOrg:     st.TargetOrg,
		targetProject: st.TargetProject,
	}
}

func (c ConnectorContext) Move() error {

	connectors, err := c.listConnectors(c.sourceOrg, c.sourceProject)
	if err != nil {
		return err
	}

	bar := progressbar.Default(int64(len(connectors)), "Connectors")
	var failed []string

	for _, conn := range connectors {
		conn.OrgIdentifier = c.targetOrg
		conn.ProjectIdentifier = c.targetProject

		err = c.createConnector(&model.CreateConnectorRequest{
			Connector: conn,
		})
		if err != nil {
			failed = append(failed, fmt.Sprintln(conn.Name, "-", err.Error()))
		}
		bar.Add(1)
	}
	bar.Finish()

	reportFailed(failed, "connectors")
	return nil
}

func (c ConnectorContext) listConnectors(org, project string) ([]*nextgen.ConnectorInfo, error) {

	api := c.source
	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"accountIdentifier":                    api.Account,
			"orgIdentifier":                        org,
			"projectIdentifier":                    project,
			"size":                                 "1000",
			"includeAllConnectorsAvailableAtScope": "false",
		}).
		SetBody(model.ListRequestBody{
			FilterType: "Connector",
		}).
		Post(BaseURL + "/ng/api/connectors/listV2")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse(resp)
	}

	result := model.ListConnectorResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	connectors := []*nextgen.ConnectorInfo{}
	for _, conn := range result.Data.Content {
		if !conn.HarnessManaged {
			connectors = append(connectors, conn.Connector)
		}
	}

	return connectors, nil
}

func (c ConnectorContext) createConnector(connector *model.CreateConnectorRequest) error {

	api := c.target
	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(connector).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
		}).
		Post(BaseURL + "/ng/api/connectors")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}

	return nil
}
