package services

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/Fernando-Dourado/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
)

var bar *progressbar.ProgressBar

type FileStoreContext struct {
	source        *SourceRequest
	target        *TargetRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
}

func NewFileStoreOperation(sourceApi *SourceRequest, targetApi *TargetRequest, st *SourceTarget) FileStoreContext {
	return FileStoreContext{
		source:        sourceApi,
		target:        targetApi,
		sourceOrg:     st.SourceOrg,
		sourceProject: st.SourceProject,
		targetOrg:     st.TargetOrg,
		targetProject: st.TargetProject,
	}
}

func (c FileStoreContext) Move() error {

	nodes, err := c.listNodes("Root", "Root", nil)
	if err != nil {
		return err
	}

	bar = progressbar.Default(int64(len(nodes)), "File Store")
	var failures []string

	for _, n := range nodes {
		if err := c.handleNode(n, failures); err != nil {
			failures = handeNodeFailure(n, failures, err)
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

func (c FileStoreContext) handleNode(n *model.FileStoreNode, failures []string) error {

	// CREATE FOLDER OR FILE
	if err := c.createNode(n); err != nil {
		return err
	}

	// A FILE DON'T HAVE CHILD NODES
	if n.Type == model.File {
		return nil
	}

	// SEARCH FOR CHILD NODES
	nodes, err := c.listNodes(n.Identifier, n.Name, &n.ParentIdentifier)
	if err != nil {
		return err
	}

	bar.ChangeMax(bar.GetMax() + len(nodes))

	// FOR EACH NODE MAKE A RECURSIVE CALL
	for _, n := range nodes {
		if err := c.handleNode(n, failures); err != nil {
			failures = handeNodeFailure(n, failures, err)
		}
		bar.Add(1)
	}

	return nil
}

func (c FileStoreContext) createNode(n *model.FileStoreNode) error {

	switch n.Type {
	case model.Folder:
		return c.createFolder(n)
	case model.File:
		// DOWNLOAD FILE
		b, err := c.downloadFile(n)
		if err != nil {
			return err
		}

		// CREATE/UPLOAD FILE
		return c.createFile(n, b)

	default:
		return fmt.Errorf("unsupported file store node type %s", n.Type)
	}
}

func (c FileStoreContext) downloadFile(n *model.FileStoreNode) ([]byte, error) {

	resp, err := c.source.Client.R().
		SetHeader("x-api-key", c.source.Token).
		SetPathParam("identifier", n.Identifier).
		SetQueryParams(createQueryParams(c.source.Account, c.sourceOrg, c.sourceProject)).
		Get(c.source.Url + "/ng/api/file-store/files/{identifier}/download")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse(resp)
	}

	return resp.Body(), nil
}

func (c FileStoreContext) createFile(n *model.FileStoreNode, b []byte) error {

	// ENSURE ONLY FILES CAN BE UPLOADED
	if n.Type != model.File {
		return fmt.Errorf("node %s is not a folder", n.Name)
	}

	reader := bytes.NewReader(b)

	resp, err := c.target.Client.R().
		SetHeader("x-api-key", c.target.Token).
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
		SetQueryParams(createQueryParams(c.target.Account, c.targetOrg, c.targetProject)).
		Post(c.target.Url + "/ng/api/file-store")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}

	return nil
}

func (c FileStoreContext) createFolder(n *model.FileStoreNode) error {

	// ENSURE ONLY FOLDERS CAN BE CREATED
	if n.Type != model.Folder {
		return fmt.Errorf("node %s is not a folder", n.Name)
	}

	resp, err := c.target.Client.R().
		SetHeader("x-api-key", c.target.Token).
		SetHeader("Content-Type", "multipart/form-data").
		SetMultipartFormData(map[string]string{
			"identifier":       n.Identifier,
			"name":             n.Name,
			"type":             string(n.Type),
			"parentIdentifier": n.ParentIdentifier,
			"description":      n.Description,
			"path":             n.Path,
		}).
		SetQueryParams(createQueryParams(c.target.Account, c.targetOrg, c.targetProject)).
		Post(c.target.Url + "/ng/api/file-store")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}

	return nil
}

func (c FileStoreContext) listNodes(identifier, name string, parentIdentifier *string) ([]*model.FileStoreNode, error) {

	req := model.GetFolderNodesRequest{
		Identifier:       identifier,
		ParentIdentifier: parentIdentifier,
		Name:             name,
		Type:             "FOLDER",
	}

	resp, err := c.source.Client.R().
		SetHeader("x-api-key", c.source.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		SetQueryParams(createQueryParams(c.source.Account, c.sourceOrg, c.sourceProject)).
		Post(c.source.Url + "/ng/api/file-store/folder")
	if err != nil {
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
