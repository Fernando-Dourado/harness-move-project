package model

type GetEnvGroupResponse struct {
	Status        string          `json:"status"`
	Data          GetEnvGroupData `json:"data"`
	MetaData      interface{}     `json:"metaData"`
	CorrelationID string          `json:"correlationId"`
}

type GetEnvGroupData struct {
	TotalPages    int64             `json:"totalPages"`
	TotalItems    int64             `json:"totalItems"`
	PageItemCount int64             `json:"pageItemCount"`
	PageSize      int64             `json:"pageSize"`
	Content       []EnvGroupContent `json:"content"`
	PageIndex     int64             `json:"pageIndex"`
	Empty         bool              `json:"empty"`
	PageToken     interface{}       `json:"pageToken"`
}

type EnvGroupContent struct {
	EnvGroup       EnvGroup `json:"envGroup"`
	CreatedAt      int64    `json:"createdAt"`
	LastModifiedAt int64    `json:"lastModifiedAt"`
}

type EnvGroupEnvironments struct {
	Environment           EnvGroup    `json:"environment"`
	CreatedAt             int64       `json:"createdAt"`
	LastModifiedAt        int64       `json:"lastModifiedAt"`
	EntityValidityDetails interface{} `json:"entityValidityDetails"`
}

type EnvGroup struct {
	AccountID         string                 `json:"accountId"`
	OrgIdentifier     string                 `json:"orgIdentifier"`
	ProjectIdentifier string                 `json:"projectIdentifier"`
	Identifier        string                 `json:"identifier"`
	Name              string                 `json:"name"`
	Description       string                 `json:"description"`
	Color             string                 `json:"color"`
	Deleted           bool                   `json:"deleted"`
	Tags              Tags                   `json:"tags"`
	EnvIdentifiers    []string               `json:"envIdentifiers,omitempty"`
	EnvResponse       []EnvGroupEnvironments `json:"envResponse,omitempty"`
	YAML              string                 `json:"yaml"`
	GitDetails        *EnvGroupGitDetails    `json:"gitDetails,omitempty"`
	Type              *string                `json:"type,omitempty"`
	StoreType         *string                `json:"storeType,omitempty"`
}

type EnvGroupGitDetails struct {
	ObjectID                 interface{} `json:"objectId"`
	Branch                   interface{} `json:"branch"`
	RepoIdentifier           interface{} `json:"repoIdentifier"`
	RootFolder               interface{} `json:"rootFolder"`
	FilePath                 interface{} `json:"filePath"`
	RepoName                 interface{} `json:"repoName"`
	CommitID                 interface{} `json:"commitId"`
	FileURL                  interface{} `json:"fileUrl"`
	RepoURL                  interface{} `json:"repoUrl"`
	ParentEntityConnectorRef interface{} `json:"parentEntityConnectorRef"`
	ParentEntityRepoName     interface{} `json:"parentEntityRepoName"`
	IsHarnessCodeRepo        interface{} `json:"isHarnessCodeRepo"`
}

type CreateEnvGroup struct {
	OrgIdentifier     string `json:"orgIdentifier"`
	ProjectIdentifier string `json:"projectIdentifier"`
	Identifier        string `json:"identifier"`
	Color             string `json:"color"`
	YAML              string `json:"yaml"`
}
