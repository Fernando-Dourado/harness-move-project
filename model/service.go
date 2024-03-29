package model

type ServiceListResult struct {
	Status        string          `json:"status"`
	Data          ServiceListData `json:"data"`
	CorrelationID string          `json:"correlationId"`
}

type ServiceListData struct {
	TotalPages    int64                 `json:"totalPages"`
	TotalItems    int64                 `json:"totalItems"`
	PageItemCount int64                 `json:"pageItemCount"`
	PageSize      int64                 `json:"pageSize"`
	Content       []*ServiceListContent `json:"content"`
	PageIndex     int64                 `json:"pageIndex"`
	Empty         bool                  `json:"empty"`
}

type ServiceListContent struct {
	Service               Service               `json:"service"`
	CreatedAt             int64                 `json:"createdAt"`
	LastModifiedAt        int64                 `json:"lastModifiedAt"`
	EntityValidityDetails EntityValidityDetails `json:"entityValidityDetails"`
}

type Service struct {
	AccountID         string    `json:"accountId"`
	Identifier        string    `json:"identifier"`
	OrgIdentifier     string    `json:"orgIdentifier"`
	ProjectIdentifier string    `json:"projectIdentifier"`
	Name              string    `json:"name"`
	Description       *string   `json:"description,omitempty"`
	Deleted           bool      `json:"deleted"`
	Yaml              string    `json:"yaml"`
	StoreType         StoreType `json:"storeType"`
}

// CREATE SERVICE

type CreateServiceRequest struct {
	Identifier        string  `json:"identifier"`
	OrgIdentifier     string  `json:"orgIdentifier"`
	ProjectIdentifier string  `json:"projectIdentifier"`
	Name              string  `json:"name"`
	Description       *string `json:"description,omitempty"`
	Yaml              string  `json:"yaml"`
}
