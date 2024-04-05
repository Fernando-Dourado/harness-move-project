package model

type GetResourceGroupResponse struct {
	Status        string                `json:"status"`
	Data          *GetResourceGroupData `json:"data"`
	CorrelationID string                `json:"correlationId"`
}

type GetResourceGroupData struct {
	TotalPages    int64                  `json:"totalPages"`
	TotalItems    int64                  `json:"totalItems"`
	PageItemCount int64                  `json:"pageItemCount"`
	PageSize      int64                  `json:"pageSize"`
	Content       []ResourceGroupContent `json:"content"`
	PageIndex     int64                  `json:"pageIndex"`
	Empty         bool                   `json:"empty"`
	PageToken     interface{}            `json:"pageToken"`
}

type ResourceGroupContent struct {
	ResourceGroup  ResourceGroup `json:"resourceGroup"`
	CreatedAt      int64         `json:"createdAt"`
	LastModifiedAt int64         `json:"lastModifiedAt"`
	HarnessManaged bool          `json:"harnessManaged"`
}

type ResourceGroup struct {
	Identifier         string                       `json:"identifier"`
	Name               string                       `json:"name"`
	Color              string                       `json:"color"`
	Tags               Tags                         `json:"tags"`
	Description        string                       `json:"description"`
	AllowedScopeLevels []string                     `json:"allowedScopeLevels"`
	IncludedScopes     []ResourceGroupIncludedScope `json:"includedScopes"`
	ResourceFilter     ResourceGroupResourceFilter  `json:"resourceFilter"`
	AccountIdentifier  *string                      `json:"accountIdentifier,omitempty"`
	OrgIdentifier      string                       `json:"orgIdentifier,omitempty"`
	ProjectIdentifier  string                       `json:"projectIdentifier,omitempty"`
}

type ResourceGroupIncludedScope struct {
	Filter            string  `json:"filter"`
	AccountIdentifier *string `json:"accountIdentifier,omitempty"`
	OrgIdentifier     *string `json:"orgIdentifier,omitempty"`
	ProjectIdentifier *string `json:"projectIdentifier,omitempty"`
}

type ResourceGroupResourceFilter struct {
	Resources           []ResourceGroupResource `json:"resources"`
	IncludeAllResources bool                    `json:"includeAllResources"`
}

type ResourceGroupResource struct {
	ResourceType string   `json:"resourceType"`
	Identifiers  []string `json:"identifiers,omitempty"`
}

type NewResourceGroupContent struct {
	ResourceGroup *ResourceGroup `json:"resourceGroup"`
}
