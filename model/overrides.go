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

type OverridesV2Type string

const (
	OV2_Global       OverridesV2Type = "ENV_GLOBAL_OVERRIDE"
	OV2_Service      OverridesV2Type = "ENV_SERVICE_OVERRIDE"
	OV2_Infra        OverridesV2Type = "INFRA_GLOBAL_OVERRIDE"
	OV2_ServiceInfra OverridesV2Type = "INFRA_SERVICE_OVERRIDE"
)

type ListOverridesV2Response struct {
	Status  string              `json:"status"`
	Code    string              `json:"code"`
	Message string              `json:"message"`
	Data    ListOverridesV2Data `json:"data"`
}

type ListOverridesV2Data struct {
	TotalPages    int64                    `json:"totalPages"`
	TotalItems    int64                    `json:"totalItems"`
	PageItemCount int64                    `json:"pageItemCount"`
	PageSize      int64                    `json:"pageSize"`
	Content       []ListOverridesV2Content `json:"content"`
	PageIndex     int64                    `json:"pageIndex"`
	Empty         bool                     `json:"empty"`
}

type ListOverridesV2Content struct {
	Identifier        string `json:"identifier"`
	AccountID         string `json:"accountId"`
	OrgIdentifier     string `json:"orgIdentifier"`
	ProjectIdentifier string `json:"projectIdentifier"`
}

type OverridesV2 struct {
	Identifier        string          `json:"identifier,omitempty"`
	AccountId         string          `json:"accountId,omitempty"`
	OrgIdentifier     string          `json:"orgIdentifier,omitempty"`
	ProjectIdentifier string          `json:"projectIdentifier,omitempty"`
	EnvironmentRef    string          `json:"environmentRef,omitempty"`
	ServiceRef        string          `json:"serviceRef,omitempty"`
	InfraIdentifier   string          `json:"infraIdentifier,omitempty"`
	Yaml              string          `json:"yaml,omitempty"`
	Type              OverridesV2Type `json:"type"`
}

type GetOverridesV2Response struct {
	Status  string      `json:"status"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    OverridesV2 `json:"data"`
}
