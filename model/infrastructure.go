package model

type InfraDefListResponse struct {
	Status        string           `json:"status"`
	Data          InfraDefListData `json:"data"`
	CorrelationID string           `json:"correlationId"`
}

type InfraDefListData struct {
	TotalPages    int64                  `json:"totalPages"`
	TotalItems    int64                  `json:"totalItems"`
	PageItemCount int64                  `json:"pageItemCount"`
	PageSize      int64                  `json:"pageSize"`
	Content       []*InfraDefListContent `json:"content"`
	PageIndex     int64                  `json:"pageIndex"`
	Empty         bool                   `json:"empty"`
}

type InfraDefListContent struct {
	Infrastructure        Infrastructure        `json:"infrastructure"`
	EntityValidityDetails EntityValidityDetails `json:"entityValidityDetails"`
}

type Infrastructure struct {
	Account           string  `json:"accountId"`
	Identifier        string  `json:"identifier"`
	OrgIdentifier     string  `json:"orgIdentifier"`
	ProjectIdentifier string  `json:"projectIdentifier"`
	Name              string  `json:"name"`
	Description       *string `json:"description,omitempty"`
	Type              string  `json:"type"`
	DeploymentType    string  `json:"deploymentType"`
	Yaml              string  `json:"yaml"`
}

type CreateInfrastructureRequest struct {
	Name              string  `json:"name"`
	Identifier        string  `json:"identifier"`
	Description       *string `json:"description,omitempty"`
	OrgIdentifier     string  `json:"orgIdentifier"`
	ProjectIdentifier string  `json:"projectIdentifier"`
	EnvironmentRef    string  `json:"environmentRef"`
	DeploymentType    string  `json:"deploymentType"`
	Type              string  `json:"type"`
	Yaml              string  `json:"yaml"`
}
