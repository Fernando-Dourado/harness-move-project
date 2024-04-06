package model

type GetServiceACcountResponse struct {
	Status        string                   `json:"status"`
	Data          []*GetServiceAccountData `json:"data"`
	CorrelationID string                   `json:"correlationId"`
}

type GetServiceAccountData struct {
	Identifier         string      `json:"identifier"`
	Name               string      `json:"name"`
	Email              string      `json:"email"`
	Description        string      `json:"description"`
	Tags               Tags        `json:"tags"`
	AccountIdentifier  string      `json:"accountIdentifier"`
	OrgIdentifier      string      `json:"orgIdentifier"`
	ProjectIdentifier  string      `json:"projectIdentifier"`
	GovernanceMetadata interface{} `json:"governanceMetadata"`
}
