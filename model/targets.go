package model

type TargetListResponse struct {
	ItemCount int64    `json:"itemCount"`
	PageCount int64    `json:"pageCount"`
	PageIndex int64    `json:"pageIndex"`
	PageSize  int64    `json:"pageSize"`
	Version   int64    `json:"version"`
	Targets   []Target `json:"targets"`
}

type Target struct {
	Account     string          `json:"account"`
	Anonymous   bool            `json:"anonymous"`
	Attributes  Attributes      `json:"attributes"`
	CreatedAt   int64           `json:"createdAt"`
	Environment string          `json:"environment"`
	Identifier  string          `json:"identifier"`
	Name        string          `json:"name"`
	Org         string          `json:"org"`
	Project     string          `json:"project"`
	Segments    []TargetSegment `json:"segments"`
}

// Also Referenced by Target Groups
type Attributes struct {
	Age      int64  `json:"age"`
	Location string `json:"location"`
}

type TargetSegment struct {
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
}

// Also Referenced by Target Groups
type Rule struct {
	Attribute string        `json:"attribute"`
	ID        string         `json:"id"`
	Negate    bool          `json:"negate"`
	Op        string        `json:"op"`
	Values    []interface{} `json:"values"`
}

// Also Referenced by Target Groups
type ServingRule struct {
	Clauses  []interface{} `json:"clauses"`
	Priority int64         `json:"priority"`
	RuleID   string        `json:"ruleId"`
}
