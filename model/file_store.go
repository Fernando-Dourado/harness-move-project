package model

type GetFolderNodesRequest struct {
	Identifier       string  `json:"identifier"`
	ParentIdentifier *string `json:"parentIdentifier,omitempty"`
	Name             string  `json:"name"`
	Type             string  `json:"type"`
}

type GetFolderNodesResponse struct {
	Status        string        `json:"status"`
	Data          FileStoreNode `json:"data"`
	CorrelationID string        `json:"correlationId"`
}

type FileStoreNode struct {
	Identifier       string           `json:"identifier"`
	ParentIdentifier string           `json:"parentIdentifier"`
	Name             string           `json:"name"`
	Type             FileStoreType    `json:"type"`
	Path             string           `json:"path"`
	LastModifiedAt   int64            `json:"lastModifiedAt"`
	LastModifiedBy   LastModifiedBy   `json:"lastModifiedBy"`
	FileUsage        string           `json:"fileUsage"`
	Description      string           `json:"description"`
	MimeType         *string          `json:"mimeType,omitempty"`
	Children         []*FileStoreNode `json:"children,omitempty"`
}

type LastModifiedBy struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type FileStoreType string

const (
	File   FileStoreType = "FILE"
	Folder FileStoreType = "FOLDER"
)
