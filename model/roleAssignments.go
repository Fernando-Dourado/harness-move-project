package model

// Generated by https://quicktype.io

type GetRoleAssignmentResponse struct {
	Status        string      `json:"status"`
	Data          GetRoleAssignmentData `json:"data"`
	CorrelationID string      `json:"correlationId"`
}

type GetRoleAssignmentData struct {
	TotalPages    int64              `json:"totalPages"`
	TotalItems    int64              `json:"totalItems"`
	PageItemCount int64              `json:"pageItemCount"`
	PageSize      int64              `json:"pageSize"`
	Content       []*RoleAssignmentListContent `json:"content"`
	PageIndex     int64              `json:"pageIndex"`
	Empty         bool               `json:"empty"`
	PageToken     interface{}        `json:"pageToken"`
}

type RoleAssignmentListContent struct {
	RoleAssignment ExistingRoleAssignment `json:"roleAssignment"`
	Scope          ExistingRoleAssignmentPrincipal   `json:"scope"`
	LastModifiedAt int64                 `json:"lastModifiedAt"`
	HarnessManaged bool                  `json:"harnessManaged"`
}

type ExistingRoleAssignment struct {
	Identifier              string    `json:"identifier"`
	ResourceGroupIdentifier string    `json:"resourceGroupIdentifier"`
	RoleIdentifier          string    `json:"roleIdentifier"`
	Principal               ExistingRoleAssignmentPrincipal `json:"principal"`
	Disabled                bool      `json:"disabled"`
	Managed                 bool      `json:"managed"`
	Internal                bool      `json:"internal"`
	OrgIdentifier           string    `json:"orgIdentifier"`
	ProjectIdentifier       string    `json:"projectIdentifier"`
}

type ExistingRoleAssignmentPrincipal struct {
	ScopeLevel *string `json:"scopeLevel"`
	Identifier string  `json:"identifier"`
	Type       string  `json:"type"`
}

type RoleAssignmentScope struct {
	AccountIdentifier string `json:"accountIdentifier"`
	OrgIdentifier     string `json:"orgIdentifier"`
	ProjectIdentifier string `json:"projectIdentifier"`
}

type NewRoleAssignment struct {
	ResourceGroupIdentifier string                  `json:"resourceGroupIdentifier"`
	RoleIdentifier          string                  `json:"roleIdentifier"`
	Principal               NewRoleAssignmentPrincipal `json:"principal"`
	OrgIdentifier           string                  `json:"orgIdentifier"`
	ProjectIdentifier       string                  `json:"projectIdentifier"`
}

type NewRoleAssignmentPrincipal struct {
	Identifier string `json:"identifier"`
	Type       string `json:"type"`
}
