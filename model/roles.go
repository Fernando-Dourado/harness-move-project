package model

type GetRolesResponse struct {
	Status        string       `json:"status"`
	Data          GetRolesData `json:"data"`
	CorrelationID string       `json:"correlationId"`
}

type GetRolesData struct {
	TotalPages    int64           `json:"totalPages"`
	TotalItems    int64           `json:"totalItems"`
	PageItemCount int64           `json:"pageItemCount"`
	PageSize      int64           `json:"pageSize"`
	Content       []*RolesContent `json:"content"`
	PageIndex     int64           `json:"pageIndex"`
	Empty         bool            `json:"empty"`
	PageToken     interface{}     `json:"pageToken"`
}

type RolesContent struct {
	Role                              ExistingRoles `json:"role"`
	Scope                             *RolesScope   `json:"scope"`
	HarnessManaged                    bool          `json:"harnessManaged"`
	CreatedAt                         int64         `json:"createdAt"`
	LastModifiedAt                    int64         `json:"lastModifiedAt"`
	RoleAssignedToUserCount           int64         `json:"roleAssignedToUserCount"`
	RoleAssignedToUserGroupCount      int64         `json:"roleAssignedToUserGroupCount"`
	RoleAssignedToServiceAccountCount int64         `json:"roleAssignedToServiceAccountCount"`
}

type ExistingRoles struct {
	Identifier         string              `json:"identifier"`
	Name               string              `json:"name"`
	Permissions        []string            `json:"permissions"`
	AllowedScopeLevels []AllowedScopeLevel `json:"allowedScopeLevels"`
	Description        string              `json:"description"`
	Tags               *Tags               `json:"tags"`
}

type Tags struct {
}

type RolesScope struct {
	AccountIdentifier string `json:"accountIdentifier"`
	OrgIdentifier     string `json:"orgIdentifier"`
	ProjectIdentifier string `json:"projectIdentifier"`
}

type AllowedScopeLevel string

type NewRole struct {
	Identifier         string              `json:"identifier"`
	Name               string              `json:"name"`
	AllowedScopeLevels []AllowedScopeLevel `json:"allowedScopeLevels"`
	Permissions        []string            `json:"permissions"`
	Description        string              `json:"description"`
	Tags               *Tags               `json:"tags"`
	OrgIdentifier      string              `json:"orgIdentifier"`
	ProjectIdentifier  string              `json:"projectIdentifier"`
}
