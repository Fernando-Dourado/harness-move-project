package model

type ProjectExistResult struct {
	Status        string           `json:"status"`
	Data          ProjectExistData `json:"data"`
	CorrelationID string           `json:"correlationId"`
}

type ProjectExistData struct {
	TotalPages    int64 `json:"totalPages"`
	TotalItems    int64 `json:"totalItems"`
	PageItemCount int64 `json:"pageItemCount"`
	PageSize      int64 `json:"pageSize"`
	PageIndex     int64 `json:"pageIndex"`
	Empty         bool  `json:"empty"`
}
