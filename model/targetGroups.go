package model

type TargetGroupListResponse struct {
	ItemCount    int64          `json:"itemCount"`
	PageCount    int64          `json:"pageCount"`
	PageIndex    int64          `json:"pageIndex"`
	PageSize     int64          `json:"pageSize"`
	Version      int64          `json:"version"`
	TargetGroups []TargetGroups `json:"segments"`
}

type TargetGroups struct {
	CreatedAt    int64         `json:"createdAt"`
	Environment  string        `json:"environment"`
	Excluded     []Cluded      `json:"excluded"`
	Identifier   string        `json:"identifier"`
	Included     []Cluded      `json:"included"`
	ModifiedAt   int64         `json:"modifiedAt"`
	Name         string        `json:"name"`
	Rules        []Rule        `json:"rules"`
	ServingRules []ServingRule `json:"servingRules"`
	Tags         []Tag         `json:"tags"`
	Version      int64         `json:"version"`
	Org          string        `json:"org"`
	Project      string        `json:"project"`
	Account      string        `json:"account"`
}

type NewTargetGroup struct {
	CreatedAt    int64         `json:"createdAt"`
	Environment  string        `json:"environment"`
	Excluded     []Cluded      `json:"excluded"`
	Identifier   string        `json:"identifier"`
	Included     []string      `json:"included"`
	ModifiedAt   int64         `json:"modifiedAt"`
	Name         string        `json:"name"`
	Rules        []Rule        `json:"rules"`
	ServingRules []ServingRule `json:"servingRules"`
	Version      int64         `json:"version"`
	Org          string        `json:"org"`
	Project      string        `json:"project"`
	Account      string        `json:"account"`
}

// Also Referenced by Targets
type Cluded struct {
	Account     string            `json:"account"`
	Anonymous   bool              `json:"anonymous"`
	Attributes  Attributes        `json:"attributes"`
	CreatedAt   int64             `json:"createdAt"`
	Environment string            `json:"environment"`
	Identifier  string            `json:"identifier"`
	Name        string            `json:"name"`
	Org         string            `json:"org"`
	Project     string            `json:"project"`
	Segments    []ExcludedSegment `json:"segments"`
}

type ExcludedSegment struct {
}
