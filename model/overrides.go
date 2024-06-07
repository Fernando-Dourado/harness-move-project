package model

type ListServiceOverridesRequest struct {
	Status        string                   `json:"status"`
	Data          ListServiceOverridesData `json:"data"`
	CorrelationID string                   `json:"correlationId"`
}

type ListServiceOverridesData struct {
	TotalPages    int64              `json:"totalPages"`
	TotalItems    int64              `json:"totalItems"`
	PageItemCount int64              `json:"pageItemCount"`
	PageSize      int64              `json:"pageSize"`
	Content       []*ServiceOverride `json:"content"`
	PageIndex     int64              `json:"pageIndex"`
	Empty         bool               `json:"empty"`
}

type ServiceOverride struct {
	AccountID         string `json:"accountId"`
	OrgIdentifier     string `json:"orgIdentifier"`
	ProjectIdentifier string `json:"projectIdentifier"`
	EnvironmentRef    string `json:"environmentRef"`
	ServiceRef        string `json:"serviceRef"`
	YAML              string `json:"yaml"`
}

type CreateServiceOverrideRequest struct {
	OrgIdentifier     string `json:"orgIdentifier"`
	ProjectIdentifier string `json:"projectIdentifier"`
	EnvironmentRef    string `json:"environmentIdentifier"`
	ServiceRef        string `json:"serviceIdentifier"`
	YAML              string `json:"yaml"`
}
