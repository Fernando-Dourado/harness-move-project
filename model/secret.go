package model

import "github.com/harness/harness-go-sdk/harness/nextgen"

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
