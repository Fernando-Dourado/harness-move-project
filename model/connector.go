package model

type ConnectorListResult struct {
	Status        string           `json:"status"`
	Data          GetConnectorData `json:"data"`
	CorrelationID string           `json:"correlationId"`
}

type GetConnectorData struct {
	TotalPages    int64              `json:"totalPages"`
	TotalItems    int64              `json:"totalItems"`
	PageItemCount int64              `json:"pageItemCount"`
	PageSize      int64              `json:"pageSize"`
	Content       []ConnectorContent `json:"content"`
	PageIndex     int64              `json:"pageIndex"`
	Empty         bool               `json:"empty"`
	PageToken     interface{}        `json:"pageToken"`
}

type ConnectorContent struct {
	Connector             Connector                      `json:"connector"`
	CreatedAt             int64                          `json:"createdAt"`
	LastModifiedAt        int64                          `json:"lastModifiedAt"`
	Status                ConnectorStatus                `json:"status"`
	ActivityDetails       ConnectorActivityDetails       `json:"activityDetails"`
	HarnessManaged        bool                           `json:"harnessManaged"`
	GitDetails            ConnectorGitDetails            `json:"gitDetails"`
	EntityValidityDetails ConnectorEntityValidityDetails `json:"entityValidityDetails"`
	GovernanceMetadata    interface{}                    `json:"governanceMetadata"`
	IsFavorite            bool                           `json:"isFavorite"`
}

type ConnectorActivityDetails struct {
	LastActivityTime int64 `json:"lastActivityTime"`
}

type Connector struct {
	Name              string        `json:"name"`
	Identifier        string        `json:"identifier"`
	Description       *string       `json:"description"`
	AccountIdentifier string        `json:"accountIdentifier"`
	OrgIdentifier     string        `json:"orgIdentifier"`
	ProjectIdentifier string        `json:"projectIdentifier"`
	Tags              Tags          `json:"tags"`
	Type              string        `json:"type"`
	Spec              ConnectorSpec `json:"spec"`
}

type ConnectorSpec struct {
	AzureArtifactsURL    *string                  `json:"azureArtifactsUrl,omitempty"`
	Auth                 *ConnectorPurpleAuth     `json:"auth,omitempty"`
	DelegateSelectors    []string                 `json:"delegateSelectors"`
	ExecuteOnDelegate    *bool                    `json:"executeOnDelegate,omitempty"`
	Credential           *ConnectorCredential     `json:"credential,omitempty"`
	URL                  *string                  `json:"url,omitempty"`
	ValidationRepo       *string                  `json:"validationRepo"`
	Authentication       *ConnectorAuthentication `json:"authentication,omitempty"`
	APIAccess            *ConnectorAPIAccess      `json:"apiAccess"`
	Proxy                *bool                    `json:"proxy,omitempty"`
	ProxyURL             interface{}              `json:"proxyUrl"`
	Type                 *string                  `json:"type,omitempty"`
	DockerRegistryURL    *string                  `json:"dockerRegistryUrl,omitempty"`
	ProviderType         *string                  `json:"providerType,omitempty"`
	AzureEnvironmentType *string                  `json:"azureEnvironmentType,omitempty"`
	Credentials          interface{}              `json:"credentials"`
	Default              *bool                    `json:"default,omitempty"`
}

type ConnectorAPIAccess struct {
	Type string        `json:"type"`
	Spec APIAccessSpec `json:"spec"`
}

type APIAccessSpec struct {
	TokenRef string `json:"tokenRef"`
}

type ConnectorPurpleAuth struct {
	Spec *ConnectorAPIAccess `json:"spec,omitempty"`
	Type *string             `json:"type,omitempty"`
}

type ConnectorAuthentication struct {
	Type string                      `json:"type"`
	Spec ConnectorAuthenticationSpec `json:"spec"`
}

type ConnectorAuthenticationSpec struct {
	Type      *string              `json:"type,omitempty"`
	Spec      *ConnectorPurpleSpec `json:"spec,omitempty"`
	SSHKeyRef *string              `json:"sshKeyRef,omitempty"`
}

type ConnectorPurpleSpec struct {
	Username    string      `json:"username"`
	UsernameRef interface{} `json:"usernameRef"`
	TokenRef    string      `json:"tokenRef"`
}

type ConnectorCredential struct {
	Type string                   `json:"type"`
	Spec *ConnectorCredentialSpec `json:"spec"`
}

type ConnectorCredentialSpec struct {
	ApplicationID string              `json:"applicationId"`
	TenantID      string              `json:"tenantId"`
	Auth          ConnectorFluffyAuth `json:"auth"`
}

type ConnectorFluffyAuth struct {
	Type string            `json:"type"`
	Spec ConnectorAuthSpec `json:"spec"`
}

type ConnectorAuthSpec struct {
	SecretRef string `json:"secretRef"`
}

type ConnectorEntityValidityDetails struct {
	Valid       bool        `json:"valid"`
	InvalidYAML interface{} `json:"invalidYaml"`
}

type ConnectorGitDetails struct {
	ObjectID                 string `json:"objectId"`
	Branch                   string `json:"branch"`
	RepoIdentifier           string `json:"repoIdentifier"`
	RootFolder               string `json:"rootFolder"`
	FilePath                 string `json:"filePath"`
	RepoName                 string `json:"repoName"`
	CommitID                 string `json:"commitId"`
	FileURL                  string `json:"fileUrl"`
	RepoURL                  string `json:"repoUrl"`
	ParentEntityConnectorRef string `json:"parentEntityConnectorRef"`
	ParentEntityRepoName     string `json:"parentEntityRepoName"`
	IsHarnessCodeRepo        string `json:"isHarnessCodeRepo"`
}

type ConnectorStatus struct {
	Status          string                 `json:"status"`
	ErrorSummary    *string                `json:"errorSummary"`
	Errors          []ConnectorStatusError `json:"errors"`
	TestedAt        int64                  `json:"testedAt"`
	LastTestedAt    int64                  `json:"lastTestedAt"`
	LastConnectedAt int64                  `json:"lastConnectedAt"`
	LastAlertSent   interface{}            `json:"lastAlertSent"`
}

type ConnectorStatusError struct {
	Reason  string `json:"reason"`
	Message string `json:"message"`
	Code    int64  `json:"code"`
}
