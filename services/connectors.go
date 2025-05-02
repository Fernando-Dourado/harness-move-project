package services

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/Fernando-Dourado/harness-move-project/model"
	"github.com/harness/harness-go-sdk/harness/nextgen"
	"github.com/schollz/progressbar/v3"
)

type ConnectorContext struct {
	source        *SourceRequest
	target        *TargetRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
}

type _ListConnectorsBody struct {
	FilterType string `json:"filterType"`
}

func NewConnectorOperation(sourceApi *SourceRequest, targetApi *TargetRequest, st *SourceTarget) ConnectorContext {
	return ConnectorContext{
		source:        sourceApi,
		target:        targetApi,
		sourceOrg:     st.SourceOrg,
		sourceProject: st.SourceProject,
		targetOrg:     st.TargetOrg,
		targetProject: st.TargetProject,
	}
}

func (c ConnectorContext) Move() error {

	connectors, err := c.listConnectors(c.sourceOrg, c.sourceProject)
	if err != nil {
		return err
	}

	bar := progressbar.Default(int64(len(connectors)), "Connectors")
	var failed []string

	for _, conn := range connectors {
		conn.OrgIdentifier = c.targetOrg
		conn.ProjectIdentifier = c.targetProject

		err = c.createConnector(&model.CreateConnectorRequest{
			Connector: conn,
		})
		if err != nil {
			failed = append(failed, fmt.Sprintln(conn.Name, "-", err.Error()))
		}
		bar.Add(1)
	}
	bar.Finish()

	reportFailed(failed, "connectors")
	return nil
}

func (c ConnectorContext) listConnectors(org, project string) ([]*nextgen.ConnectorInfo, error) {

	body := _ListConnectorsBody{
		FilterType: "Connector",
	}

	api := c.source
	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"accountIdentifier":                    api.Account,
			"orgIdentifier":                        org,
			"projectIdentifier":                    project,
			"size":                                 "1000",
			"includeAllConnectorsAvailableAtScope": "false",
		}).
		SetBody(body).
		Post(BaseURL + "/ng/api/connectors/listV2")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse(resp)
	}

	result := model.ListConnectorResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	connectors := []*nextgen.ConnectorInfo{}
	for _, conn := range result.Data.Content {
		if !conn.HarnessManaged {
			connectors = append(connectors, conn.Connector)
		}
	}

	return connectors, nil
}

func (c ConnectorContext) createConnector(connector *model.CreateConnectorRequest) error {

	api := c.target
	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(connector).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
		}).
		Post(BaseURL + "/ng/api/connectors")
	if err != nil {
		return err
	}
	if resp.IsError() {
		err = handleConnectorErrorResponse(resp)

		if secretErr, ok := err.(*MissingSecretError); ok {
			return c.handleMissingSecret(secretErr.Id, connector)
		}
		return err
	}

	return nil
}

func (c ConnectorContext) handleMissingSecret(secretRef string, connector *model.CreateConnectorRequest) error {

	secret, err := c.getSecret(secretRef)
	if err != nil {
		return err
	}
	secret.OrgIdentifier = c.targetOrg
	secret.ProjectIdentifier = c.targetProject

	newSecret := &model.CreateSecretRequest{
		Secret: secret,
	}
	if secret.Type_ == nextgen.SecretTypes.SecretText {
		err = c.createSecretText(newSecret)
	} else if secret.Type_ == nextgen.SecretTypes.SSHKey {
		err = c.createSSHKey(newSecret)
	} else {
		return fmt.Errorf("secret type %s not supported", secret.Type_)
	}

	if err != nil {
		return err
	}

	// Retry creating the connector after creating the secret
	return c.createConnector(connector)
}

func (c ConnectorContext) getSecret(secretRef string) (*nextgen.Secret, error) {

	api := c.source
	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetPathParam("secretRef", secretRef).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     c.sourceOrg,
			"projectIdentifier": c.sourceProject,
		}).
		Get(BaseURL + "/ng/api/v2/secrets/{secretRef}")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse(resp)
	}

	result := model.GetSecretResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	return result.Data.Secret, nil
}

func (c ConnectorContext) createSecretText(secret *model.CreateSecretRequest) error {

	// SET A PLACEHOLDER FOR THE SECRET VALUE
	secret.Secret.Text.Value = "PLEASE_FIX_ME"

	api := c.target
	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(secret).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     c.targetOrg,
			"projectIdentifier": c.targetProject,
			"privateSecret":     "false",
		}).
		Post(BaseURL + "/ng/api/v2/secrets")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}

	return nil
}

func (c ConnectorContext) createSSHKey(secret *model.CreateSecretRequest) error {

	if secret.Secret.SSHKey.Auth.Type_ != nextgen.SSHAuthenticationTypes.SSH {
		return fmt.Errorf("ssh auth %s not supported", secret.Secret.SSHKey.Auth.Type_)
	}
	if secret.Secret.SSHKey.Auth.SSHConfig.KeyReferenceCredential == nil {
		return fmt.Errorf("ssh config type %s not supported", secret.Secret.SSHKey.Auth.SSHConfig.CredentialType)
	}

	// GET & CREATE SECRET

	sshId := secret.Secret.SSHKey.Auth.SSHConfig.KeyReferenceCredential.Key
	sshSecret, err := c.getSecret(sshId)
	if err != nil {
		return err
	}

	err = c.createSecretFile(&model.CreateSecretRequest{
		Secret: sshSecret,
	})
	if err != nil {
		return err
	}

	// --

	resp, err := c.target.Client.R().
		SetHeader("x-api-key", c.target.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(secret).
		SetQueryParams(map[string]string{
			"accountIdentifier": c.target.Account,
			"orgIdentifier":     c.targetOrg,
			"projectIdentifier": c.targetProject,
			"privateSecret":     "false",
		}).
		Post(BaseURL + "/ng/api/v2/secrets")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}

	return nil
}

func (c ConnectorContext) createSecretFile(secret *model.CreateSecretRequest) error {

	newSecret := &model.CreateSecretRequest{
		Secret: &nextgen.Secret{
			Type_:             nextgen.SecretTypes.SecretFile,
			OrgIdentifier:     c.targetOrg,
			ProjectIdentifier: c.targetProject,
			Identifier:        secret.Secret.Identifier,
			Name:              secret.Secret.Name,
			Description:       secret.Secret.Description,
			Tags:              secret.Secret.Tags,
			File: &nextgen.SecretFileSpe{
				SecretManagerIdentifier: secret.Secret.File.SecretManagerIdentifier,
			},
		},
	}

	content := []byte("PLEASE_FIX_ME")
	reader := bytes.NewReader(content)

	resp, err := c.target.Client.R().
		SetHeader("x-api-key", c.target.Token).
		SetHeader("Content-Type", "multipart/form-data").
		SetMultipartField("file", "blob", "plain/text", reader).
		SetMultipartFormData(map[string]string{
			"spec": func() string {
				jsonData, _ := json.Marshal(newSecret)
				return string(jsonData)
			}(),
		}).
		SetQueryParams(map[string]string{
			"accountIdentifier": c.target.Account,
			"orgIdentifier":     c.targetOrg,
			"projectIdentifier": c.targetProject,
			"privateSecret":     "false",
		}).
		Post(BaseURL + "/ng/api/v2/secrets/files")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}

	return nil
}
