package model

type GetUserGroupsResponse struct {
	Status        string            `json:"status"`
	Data          GetUserGroupsData `json:"data"`
	CorrelationID string            `json:"correlationId"`
}

type GetUserGroupsData struct {
	TotalPages    int64        `json:"totalPages"`
	TotalItems    int64        `json:"totalItems"`
	PageItemCount int64        `json:"pageItemCount"`
	PageSize      int64        `json:"pageSize"`
	Content       []*UserGroup `json:"content"`
	PageIndex     int64        `json:"pageIndex"`
	Empty         bool         `json:"empty"`
	PageToken     interface{}  `json:"pageToken"`
}

type UserGroup struct {
	AccountIdentifier   string        `json:"accountIdentifier"`
	OrgIdentifier       string        `json:"orgIdentifier"`
	ProjectIdentifier   string        `json:"projectIdentifier"`
	Identifier          string        `json:"identifier"`
	Name                string        `json:"name"`
	Users               []string      `json:"users"`
	NotificationConfigs []interface{} `json:"notificationConfigs"`
	ExternallyManaged   bool          `json:"externallyManaged"`
	Description         string        `json:"description"`
	Tags                Tags          `json:"tags"`
	HarnessManaged      bool          `json:"harnessManaged"`
	SsoLinked           bool          `json:"ssoLinked"`
}

type UserGroupLookup struct {
	Identifier        string `json:"identifier"`
	OrgIdentifier     string `json:"orgIdentifier"`
	ProjectIdentifier string `json:"projectIdentifier"`
}

type UserGroupEmail struct {
	Name              string `json:"name"`
	EmailAddress      string `json:"emails"`
	OrgIdentifier     string `json:"orgIdentifier,omitempty"`
	ProjectIdentifier string `json:"projectIdentifier,omitempty"`
}
