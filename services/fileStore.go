package services

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/jf781/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
	"go.uber.org/zap"
)

var bar *progressbar.ProgressBar

type FileStoreContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
	logger        *zap.Logger
}

func NewFileStoreOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string, logger *zap.Logger) FileStoreContext {
	return FileStoreContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
		logger:        logger,
	}
}

func (c FileStoreContext) Copy() error {

	c.logger.Info("Copying file stores",
		zap.String("project", c.sourceProject),
	)

	nodes, err := c.listNodes("Root", "Root", nil, c.logger)
	if err != nil {
		c.logger.Error("Failed to list file store nodes", zap.Error(err))
		return err
	}

	bar = progressbar.Default(int64(len(nodes)), "File Store")
	var failures []string

	for _, n := range nodes {

		IncrementFileStoresTotal()

		if err := c.handleNode(n, failures, c.logger); err != nil {
			c.logger.Error("Failed to handle file", zap.Error(err))
			failures = handeNodeFailure(n, failures, err)
		} else {
			IncrementFileStoresMoved()
		}
		bar.Add(1)
	}
	bar.Finish()

	reportFailed(failures, "file store nodes:")
	return nil
}

func handeNodeFailure(node *model.FileStoreNode, failures []string, err error) []string {
	return append(failures, fmt.Sprintf("%s (%s) - %s", node.Name, node.Path, err.Error()))
}

func (c FileStoreContext) handleNode(n *model.FileStoreNode, failures []string, logger *zap.Logger) error {

	// CREATE FOLDER OR FILE
	if err := c.createNode(n, c.logger); err != nil {
		logger.Error("Failed to create or folder", zap.Error(err))
		return err
	}

	// A FILE DON'T HAVE CHILD NODES
	if n.Type == model.File {
		logger.Info("File does not have any child nodes")
		return nil
	}

	// SEARCH FOR CHILD NODES
	nodes, err := c.listNodes(n.Identifier, n.Name, &n.ParentIdentifier, c.logger)
	if err != nil {
		logger.Error("Failed identify child nodes", zap.Error(err))
		return err
	}

	bar.ChangeMax(bar.GetMax() + len(nodes))

	// FOR EACH NODE MAKE A RECURSIVE CALL
	for _, n := range nodes {
		if err := c.handleNode(n, failures, c.logger); err != nil {
			logger.Error("Error creating file or directory", zap.Error(err))
			failures = handeNodeFailure(n, failures, err)
		}
		bar.Add(1)
	}

	return nil
}

func (c FileStoreContext) createNode(n *model.FileStoreNode, logger *zap.Logger) error {

	switch n.Type {
	case model.Folder:
		return c.createFolder(n, c.logger)
	case model.File:
		// DOWNLOAD FILE
		b, err := c.downloadFile(n, c.logger)
		if err != nil {
			logger.Error("Failed download file", zap.Error(err))
			return err
		}

		// CREATE/UPLOAD FILE
		return c.createFile(n, b, c.logger)

	default:
		return fmt.Errorf("unsupported file store node type %s", n.Type)
	}
}

func (c FileStoreContext) downloadFile(n *model.FileStoreNode, logger *zap.Logger) ([]byte, error) {

	IncrementApiCalls()

	resp, err := c.api.Client.R().
		SetHeader("x-api-key", c.api.Token).
		SetPathParam("identifier", n.Identifier).
		SetQueryParams(map[string]string{
			"accountIdentifier": c.api.Account,
			"orgIdentifier":     c.sourceOrg,
			"projectIdentifier": c.sourceProject,
		}).
		Get(c.api.BaseURL + "/ng/api/file-store/files/{identifier}/download")
	if err != nil {
		logger.Error("Failed to request download file",
			zap.Error(err),
		)
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse(resp)
	}

	return resp.Body(), nil
}

func (c FileStoreContext) createFile(n *model.FileStoreNode, b []byte, logger *zap.Logger) error {

	IncrementApiCalls()

	// ENSURE ONLY FILES CAN BE UPLOADED
	if n.Type != model.File {
		return fmt.Errorf("node %s is not a folder", n.Name)
	}

	reader := bytes.NewReader(b)

	resp, err := c.api.Client.R().
		SetHeader("x-api-key", c.api.Token).
		SetHeader("Content-Type", "multipart/form-data").
		SetMultipartField("content", "blob", "plain/text", reader).
		SetMultipartFormData(map[string]string{
			"identifier":       n.Identifier,
			"name":             n.Name,
			"type":             string(n.Type),
			"parentIdentifier": n.ParentIdentifier,
			"description":      n.Description,
			"path":             n.Path,
			"fileUsage":        n.FileUsage,
			"mimeType":         *n.MimeType,
		}).
		SetQueryParams(map[string]string{
			"accountIdentifier": c.api.Account,
			"orgIdentifier":     c.targetOrg,
			"projectIdentifier": c.targetProject,
		}).
		Post(c.api.BaseURL + "/ng/api/file-store")
	if err != nil {
		logger.Error("Failed to upload file",
			zap.Error(err),
		)
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}

	return nil
}

func (c FileStoreContext) createFolder(n *model.FileStoreNode, logger *zap.Logger) error {

	IncrementApiCalls()

	// ENSURE ONLY FOLDERS CAN BE CREATED
	if n.Type != model.Folder {
		return fmt.Errorf("node %s is not a folder", n.Name)
	}

	resp, err := c.api.Client.R().
		SetHeader("x-api-key", c.api.Token).
		SetHeader("Content-Type", "multipart/form-data").
		SetMultipartFormData(map[string]string{
			"identifier":       n.Identifier,
			"name":             n.Name,
			"type":             string(n.Type),
			"parentIdentifier": n.ParentIdentifier,
			"description":      n.Description,
			"path":             n.Path,
		}).
		SetQueryParams(map[string]string{
			"accountIdentifier": c.api.Account,
			"orgIdentifier":     c.targetOrg,
			"projectIdentifier": c.targetProject,
		}).
		Post(c.api.BaseURL + "/ng/api/file-store")
	if err != nil {
		logger.Error("Failed to create folder",
			zap.Error(err),
		)
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}

	return nil
}

func (c FileStoreContext) listNodes(identifier, name string, parentIdentifier *string, logger *zap.Logger) ([]*model.FileStoreNode, error) {

	IncrementApiCalls()

	req := model.GetFolderNodesRequest{
		Identifier:       identifier,
		ParentIdentifier: parentIdentifier,
		Name:             name,
		Type:             "FOLDER",
	}

	resp, err := c.api.Client.R().
		SetHeader("x-api-key", c.api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		SetQueryParams(map[string]string{
			"accountIdentifier": c.api.Account,
			"orgIdentifier":     c.sourceOrg,
			"projectIdentifier": c.sourceProject,
		}).
		Post(c.api.BaseURL + "/ng/api/file-store/folder")
	if err != nil {
		logger.Error("Failed to retrieve list of nodes",
			zap.Error(err),
		)
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse(resp)
	}

	result := model.GetFolderNodesResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	return result.Data.Children, nil
}
