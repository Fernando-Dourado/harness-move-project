package model

type Pageable struct {
	Offset     int64 `json:"offset"`
	PageSize   int64 `json:"pageSize"`
	PageNumber int64 `json:"pageNumber"`
	Paged      bool  `json:"paged"`
	Unpaged    bool  `json:"unpaged"`
}
