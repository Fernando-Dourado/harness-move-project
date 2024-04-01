package model

type GetVariablesResponse struct {
	Status        string           `json:"status"`
	Data          GetVariablesData `json:"data"`
	CorrelationID string           `json:"correlationId"`
}

type GetVariablesData struct {
	TotalPages    int64                  `json:"totalPages"`
	TotalItems    int64                  `json:"totalItems"`
	PageItemCount int64                  `json:"pageItemCount"`
	PageSize      int64                  `json:"pageSize"`
	Content       []*GetVariablesContent `json:"content"`
	PageIndex     int64                  `json:"pageIndex"`
	Empty         bool                   `json:"empty"`
}

type GetVariablesContent struct {
	Variable       *Variable `json:"variable"`
	CreatedAt      int64     `json:"createdAt"`
	LastModifiedAt int64     `json:"lastModifiedAt"`
}

type Variable struct {
	Identifier        string       `json:"identifier"`
	Name              string       `json:"name"`
	Description       *string      `json:"description,omitempty"`
	OrgIdentifier     string       `json:"orgIdentifier"`
	ProjectIdentifier string       `json:"projectIdentifier"`
	Type              string       `json:"type"`
	Spec              SpecVariable `json:"spec"`
}

type SpecVariable struct {
	ValueType string  `json:"valueType"`
	Type      string  `json:"type"`
	Value     *string `json:"fixedValue,omitempty"`
}

type CreateVariableRequest struct {
	Variable *Variable `json:"variable"`
}
