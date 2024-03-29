package model

type ListEnvironmentResponse struct {
	Status        string              `json:"status"`
	Data          ListEnvironmentData `json:"data"`
	CorrelationID string              `json:"correlationId"`
}

type ListEnvironmentData struct {
	TotalPages    int64                     `json:"totalPages"`
	TotalItems    int64                     `json:"totalItems"`
	PageItemCount int64                     `json:"pageItemCount"`
	PageSize      int64                     `json:"pageSize"`
	Content       []*ListEnvironmentContent `json:"content"`
	PageIndex     int64                     `json:"pageIndex"`
	Empty         bool                      `json:"empty"`
}

type ListEnvironmentContent struct {
	Environment ListEnvironment `json:"environment"`
}

type ListEnvironment struct {
	Account           string    `json:"accountId"`
	OrgIdentifier     string    `json:"orgIdentifier"`
	ProjectIdentifier string    `json:"projectIdentifier"`
	Identifier        string    `json:"identifier"`
	Name              string    `json:"name"`
	Description       *string   `json:"description,omitempty"`
	Type              string    `json:"type"`
	Deleted           bool      `json:"deleted"`
	Yaml              string    `json:"yaml"`
	Color             string    `json:"color"`
	StoreType         StoreType `json:"storeType"`
}

type CreateEnvironmentRequest struct {
	OrgIdentifier     string  `json:"orgIdentifier"`
	ProjectIdentifier string  `json:"projectIdentifier"`
	Identifier        string  `json:"identifier"`
	Name              string  `json:"name"`
	Description       *string `json:"description,omitempty"`
	Color             string  `json:"color"`
	Type              string  `json:"type"`
	Yaml              string  `json:"yaml"`
}
