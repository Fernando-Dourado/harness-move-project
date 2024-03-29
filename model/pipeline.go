package model

type PipelineListResult struct {
	Status        string           `json:"status"`
	Data          PipelineListData `json:"data"`
	CorrelationID string           `json:"correlationId"`
}

type PipelineListData struct {
	Content          []*PipelineListContent `json:"content"`
	Pageable         Pageable               `json:"pageable"`
	TotalElements    int64                  `json:"totalElements"`
	Last             bool                   `json:"last"`
	TotalPages       int64                  `json:"totalPages"`
	Size             int64                  `json:"size"`
	Number           int64                  `json:"number"`
	First            bool                   `json:"first"`
	NumberOfElements int64                  `json:"numberOfElements"`
	Empty            bool                   `json:"empty"`
}

type PipelineListContent struct {
	Name          string      `json:"name"`
	Identifier    string      `json:"identifier"`
	Version       int64       `json:"version"`
	NumOfStages   int64       `json:"numOfStages"`
	CreatedAt     int64       `json:"createdAt"`
	LastUpdatedAt int64       `json:"lastUpdatedAt"`
	GitDetails    *GitDetails `json:"gitDetails,omitempty"`
	StoreType     StoreType   `json:"storeType"`
	ConnectorRef  *string     `json:"connectorRef,omitempty"`
	Description   *string     `json:"description,omitempty"`
}

type GitDetails struct {
	FilePath string `json:"filePath"`
	RepoName string `json:"repoName"`
	RepoURL  string `json:"repoUrl"`
}

type Status string

type StageType string

type StoreType string

const (
	Inline StoreType = "INLINE"
	Remote StoreType = "REMOTE"
)
