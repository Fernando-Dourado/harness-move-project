package model

import "github.com/harness/harness-go-sdk/harness/nextgen"

type ListSecretResponse struct {
	Status        string         `json:"status"`
	Data          ListSecretData `json:"data"`
	CorrelationID string         `json:"correlationId"`
}

type ListSecretData struct {
	TotalPages    int64                     `json:"totalPages"`
	TotalItems    int64                     `json:"totalItems"`
	PageItemCount int64                     `json:"pageItemCount"`
	PageSize      int64                     `json:"pageSize"`
	Content       []*nextgen.SecretResponse `json:"content"`
	PageIndex     int64                     `json:"pageIndex"`
	Empty         bool                      `json:"empty"`
}

type GetSecretResponse struct {
	Status        string        `json:"status"`
	Data          GetSecretData `json:"data"`
	CorrelationID string        `json:"correlationId"`
}

type GetSecretData struct {
	Secret *nextgen.Secret `json:"secret"`
}

type CreateSecretRequest struct {
	Secret *nextgen.Secret `json:"secret"`
}
