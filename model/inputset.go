package model

type ListInputsetResponse struct {
	Status        string           `json:"status"`
	Data          ListInputsetData `json:"data"`
	CorrelationID string           `json:"correlationId"`
}

type ListInputsetData struct {
	TotalPages int64                  `json:"totalPages"`
	TotalItems int64                  `json:"totalItems"`
	Content    []*ListInputsetContent `json:"content"`
	PageIndex  int64                  `json:"pageIndex"`
	Empty      bool                   `json:"empty"`
}

type ListInputsetContent struct {
	Identifier            string                `json:"identifier"`
	Name                  string                `json:"name"`
	PipelineIdentifier    string                `json:"pipelineIdentifier"`
	InputSetType          string                `json:"inputSetType"`
	EntityValidityDetails EntityValidityDetails `json:"entityValidityDetails"`
}

type GetInputsetResponse struct {
	Status        string           `json:"status"`
	Data          *GetInputsetData `json:"data"`
	CorrelationID string           `json:"correlationId"`
}

type GetInputsetData struct {
	Account            string `json:"accountId"`
	OrgIdentifier      string `json:"orgIdentifier"`
	ProjectIdentifier  string `json:"projectIdentifier"`
	PipelineIdentifier string `json:"pipelineIdentifier"`
	Identifier         string `json:"identifier"`
	Yaml               string `json:"inputSetYaml"`
	Name               string `json:"name"`
	Outdated           bool   `json:"outdated"`
}
