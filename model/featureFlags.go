package model

type FeatureFlagListResult struct {
	ItemCount int64         `json:"itemCount"`
	PageCount int64         `json:"pageCount"`
	PageIndex int64         `json:"pageIndex"`
	PageSize  int64         `json:"pageSize"`
	Features  []FeatureFlag `json:"tags"`
}

type FeatureFlag struct {
	Archived             bool                     `json:"archived"`
	CreatedAt            int64                    `json:"createdAt"`
	DefaultOffVariation  string                   `json:"defaultOffVariation"`
	DefaultOnVariation   string                   `json:"defaultOnVariation"`
	Description          string                   `json:"description"`
	EnvProperties        FeatureFlagEnvProperties `json:"envProperties"`
	Evaluation           string                   `json:"evaluation"`
	EvaluationIdentifier string                   `json:"evaluationIdentifier"`
	Identifier           string                   `json:"identifier"`
	Kind                 string                   `json:"kind"`
	ModifiedAt           int64                    `json:"modifiedAt"`
	Name                 string                   `json:"name"`
	Owner                []string                 `json:"owner"`
	Permanent            bool                     `json:"permanent"`
	Prerequisites        []interface{}            `json:"prerequisites"`
	Project              string                   `json:"project"`
	Services             []interface{}            `json:"services"`
	Tags                 interface{}              `json:"tags"`
	Variations           []Variation              `json:"variations"`
	OrgIdentifier         string              `json:"orgIdentifier"`
	ProjectIdentifier     string              `json:"projectIdentifier"`
}

type FeatureFlagEnvProperties struct {
	DefaultServe       DefaultServe `json:"defaultServe"`
	Environment        string       `json:"environment"`
	JiraEnabled        bool         `json:"jiraEnabled"`
	ModifiedAt         int64        `json:"modifiedAt"`
	OffVariation       string       `json:"offVariation"`
	PipelineConfigured bool         `json:"pipelineConfigured"`
	Rules              interface{}  `json:"rules"`
	State              string       `json:"state"`
	VariationMap       interface{}  `json:"variationMap"`
	Version            int64        `json:"version"`
}

type DefaultServe struct {
}

type Variation struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
	Value      string `json:"value"`
}
