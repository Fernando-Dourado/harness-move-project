package model

import "github.com/harness/harness-go-sdk/harness/nextgen"

type ListConnectorResponse struct {
	Status        string            `json:"status"`
	Data          ListConnectorData `json:"data"`
	CorrelationID string            `json:"correlationId"`
}

type ListConnectorData struct {
	TotalPages    int64                        `json:"totalPages"`
	TotalItems    int64                        `json:"totalItems"`
	PageItemCount int64                        `json:"pageItemCount"`
	PageSize      int64                        `json:"pageSize"`
	Content       []*nextgen.ConnectorResponse `json:"content"`
	PageIndex     int64                        `json:"pageIndex"`
	Empty         bool                         `json:"empty"`
}

type CreateConnectorRequest struct {
	Connector *nextgen.ConnectorInfo `json:"connector"`
}
