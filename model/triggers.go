package model

type GetTriggerResposne struct {
	Status        string         `json:"status"`
	Data          GetTriggerData `json:"data"`
	MetaData      interface{}    `json:"metaData"`
	CorrelationID string         `json:"correlationId"`
}

type GetTriggerData struct {
	TotalPages    int64            `json:"totalPages"`
	TotalItems    int64            `json:"totalItems"`
	PageItemCount int64            `json:"pageItemCount"`
	PageSize      int64            `json:"pageSize"`
	Content       []TriggerContent `json:"content"`
	PageIndex     int64            `json:"pageIndex"`
	Empty         bool             `json:"empty"`
	PageToken     interface{}      `json:"pageToken"`
}

type TriggerContent struct {
	Name                  string              `json:"name"`
	Identifier            string              `json:"identifier"`
	Type                  string              `json:"type"`
	TriggerStatus         TriggerStatus       `json:"triggerStatus"`
	BuildDetails          TriggerBuildDetails `json:"buildDetails"`
	Tags                  Tags                `json:"tags"`
	Executions            []int64             `json:"executions"`
	YAML                  string              `json:"yaml"`
	WebhookURL            string              `json:"webhookUrl"`
	WebhookCurlCommand    string              `json:"webhookCurlCommand"`
	Enabled               bool                `json:"enabled"`
	YAMLVersion           string              `json:"yamlVersion"`
	PipelineInputOutdated bool                `json:"pipelineInputOutdated"`
	OrgIdentifier         string              `json:"orgIdentifier"`
	ProjectIdentifier     string              `json:"projectIdentifier"`
}

type TriggerBuildDetails struct {
	BuildType string `json:"buildType"`
}

type TriggerStatus struct {
	PollingSubscriptionStatus     TriggerPollingSubscriptionStatus `json:"pollingSubscriptionStatus"`
	ValidationStatus              TriggerValidationStatus          `json:"validationStatus"`
	WebhookAutoRegistrationStatus interface{}                      `json:"webhookAutoRegistrationStatus"`
	WebhookInfo                   interface{}                      `json:"webhookInfo"`
	Status                        string                           `json:"status"`
	DetailMessages                []string                         `json:"detailMessages"`
	LastPollingUpdate             int64                            `json:"lastPollingUpdate"`
	LastPolled                    []interface{}                    `json:"lastPolled"`
}

type TriggerPollingSubscriptionStatus struct {
	StatusResult          string        `json:"statusResult"`
	DetailedMessage       string        `json:"detailedMessage"`
	LastPolled            []interface{} `json:"lastPolled"`
	LastPollingUpdate     int64         `json:"lastPollingUpdate"`
	ErrorStatusValidUntil int64         `json:"errorStatusValidUntil"`
}

type TriggerValidationStatus struct {
	StatusResult    string      `json:"statusResult"`
	DetailedMessage interface{} `json:"detailedMessage"`
}
