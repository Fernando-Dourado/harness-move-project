package model

type GetUserResponse struct {
	Status        string       `json:"status"`
	Data          GetUserData  `json:"data"`
	CorrelationID string       `json:"correlationId"`
}

type GetUserData struct {
	TotalPages    int64         `json:"totalPages"`
	TotalItems    int64         `json:"totalItems"`
	PageItemCount int64         `json:"pageItemCount"`
	PageSize      int64         `json:"pageSize"`
	Content       []UserContent `json:"content"`
	PageIndex     int64         `json:"pageIndex"`
	Empty         bool          `json:"empty"`
	PageToken     interface{}   `json:"pageToken"`
}

type UserContent struct {
	User                   User                          `json:"user"`
	RoleAssignmentMetadata []UserRoleAssignmentMetadatum `json:"roleAssignmentMetadata"`
}

type UserRoleAssignmentMetadatum struct {
	Identifier              string `json:"identifier"`
	RoleIdentifier          string `json:"roleIdentifier"`
	RoleName                string `json:"roleName"`
	ResourceGroupIdentifier string `json:"resourceGroupIdentifier"`
	ResourceGroupName       string `json:"resourceGroupName"`
	ManagedRole             bool   `json:"managedRole"`
	ManagedRoleAssignment   bool   `json:"managedRoleAssignment"`
}

type User struct {
	Name                           string `json:"name"`
	Email                          string `json:"email"`
	UUID                           string `json:"uuid"`
	Locked                         bool   `json:"locked"`
	Disabled                       bool   `json:"disabled"`
	ExternallyManaged              bool   `json:"externallyManaged"`
	TwoFactorAuthenticationEnabled bool   `json:"twoFactorAuthenticationEnabled"`
}

type UserEmail struct {
	EmailAddress      []string `json:"emails"`
	OrgIdentifier     string   `json:"orgIdentifier,omitempty"`
	ProjectIdentifier string   `json:"projectIdentifier,omitempty"`
}
