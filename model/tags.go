package model

// type CreateTagRequest struct {
// 	Environment string `json:"environment"`
// 	OrgIdentifier string `json:"orgIdentifier"`
// 	ProjectIdentifier string `json:"projectIdentifier"`
// }

type TagListResult struct {
	ItemCount int64 `json:"itemCount"`
	PageCount int64 `json:"pageCount"`
	PageIndex int64 `json:"pageIndex"`
	PageSize  int64 `json:"pageSize"`
	Version   int64 `json:"version"`
	Tags      []Tag `json:"tags"`
}

type Tag struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

type CreateTagRequest struct {
	Identifier        string `json:"identifier"`
	Name              string `json:"name"`
	OrgIdentifier     string `json:"orgIdentifier"`
	ProjectIdentifier string `json:"projectIdentifier"`
}
