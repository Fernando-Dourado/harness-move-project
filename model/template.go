package model

type TemplateListResult []TemplateListResultElement

type TemplateListResultElement struct {
	Account        string  `json:"account"`
	Org            string  `json:"org"`
	Project        string  `json:"project"`
	Identifier     string  `json:"identifier"`
	Name           string  `json:"name"`
	Description    *string `json:"description,omitempty"`
	VersionLabel   string  `json:"version_label"`
	EntityType     string  `json:"entity_type"`
	Scope          string  `json:"scope"`
	StoreType      string  `json:"store_type"`
	StableTemplate bool    `json:"stable_template"`
}

type TemplateGetResult struct {
	Status        string           `json:"status"`
	Data          *TemplateGetData `json:"data"`
	CorrelationID string           `json:"correlationId"`
}

type TemplateGetData struct {
	Account           string    `json:"accountId"`
	OrgIdentifier     string    `json:"orgIdentifier"`
	ProjectIdentifier string    `json:"projectIdentifier"`
	Identifier        string    `json:"identifier"`
	Yaml              string    `json:"yaml"`
	VersionLabel      string    `json:"versionLabel"`
	StoreType         StoreType `json:"storeType"`
}

type CreateTemplateRequest struct {
	Yaml        string  `json:"template_yaml"`
	Identifier  string  `json:"identifier"`
	Name        string  `json:"name"`
	Label       string  `json:"label"`
	Description *string `json:"description,omitempty"`
	Stable      bool    `json:"is_stable"`
}
