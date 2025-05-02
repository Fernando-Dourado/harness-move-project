package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/Fernando-Dourado/harness-move-project/model"
	"github.com/harness/harness-go-sdk/harness/nextgen"
	"github.com/schollz/progressbar/v3"
)

type SecretContext struct {
	source        *SourceRequest
	target        *TargetRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
}

func NewSecretOperation(sourceApi *SourceRequest, targetApi *TargetRequest, st *SourceTarget) SecretContext {
	return SecretContext{
		source:        sourceApi,
		target:        targetApi,
		sourceOrg:     st.SourceOrg,
		sourceProject: st.SourceProject,
		targetOrg:     st.TargetOrg,
		targetProject: st.TargetProject,
	}
}

func (sc SecretContext) Move() error {

	secrets, err := sc.listSecrets(sc.sourceOrg, sc.sourceProject)
	if err != nil {
		return err
	}

	// SORT SECRETS USING CUSTOM ORDER
	order := map[nextgen.SecretType]int{
		nextgen.SecretTypes.SecretText:       1,
		nextgen.SecretTypes.SecretFile:       2,
		nextgen.SecretTypes.SSHKey:           3,
		nextgen.SecretTypes.WinRmCredentials: 4,
	}
	sort.Slice(secrets, func(i, j int) bool {
		return order[secrets[i].Type_] < order[secrets[j].Type_]
	})

	bar := progressbar.Default(int64(len(secrets)), "Secrets")
	var failed []string
	for _, secret := range secrets {
		secret.OrgIdentifier = sc.targetOrg
		secret.ProjectIdentifier = sc.targetProject

		err = sc.createSecret(secret)
		if err != nil {
			failed = append(failed, fmt.Sprintln(secret.Name, "-", err.Error()))
		}

		bar.Add(1)
	}
	bar.Finish()

	reportFailed(failed, "secrets")
	return nil
}

func (sc SecretContext) listSecrets(org string, project string) ([]*nextgen.Secret, error) {

	api := sc.source
	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     org,
			"projectIdentifier": project,
			"size":              "1000",
		}).
		SetBody(model.ListRequestBody{
			FilterType: "Secret",
		},
		).
		Post(BaseURL + "/ng/api/v2/secrets/list/secrets")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse(resp)
	}

	result := model.ListSecretResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	secrets := []*nextgen.Secret{}
	for _, e := range result.Data.Content {
		secrets = append(secrets, e.Secret)
	}

	return secrets, nil
}

func (sc SecretContext) createSecret(secret *nextgen.Secret) error {
	switch secret.Type_ {
	case nextgen.SecretTypes.SecretText:
		return sc.createSecretText(secret)

	case nextgen.SecretTypes.SecretFile:
		return sc.createSecretFile(secret)

	case nextgen.SecretTypes.SSHKey:
		return sc.createSecretSSHKey(secret)

	default:
		return fmt.Errorf("secret type %s not supported", secret.Type_)
	}
}

func (sc SecretContext) createSecretText(secret *nextgen.Secret) error {

	secret.Text.Value = "PLEASE_FIX_ME"
	body := &model.CreateSecretRequest{
		Secret: secret,
	}

	api := sc.target
	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     sc.targetOrg,
			"projectIdentifier": sc.targetProject,
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

func (sc SecretContext) createSecretFile(secret *nextgen.Secret) error {

	content := []byte("PLEASE_FIX_ME")
	reader := bytes.NewReader(content)
	body := &model.CreateSecretRequest{
		Secret: secret,
	}

	resp, err := sc.target.Client.R().
		SetHeader("x-api-key", sc.target.Token).
		SetHeader("Content-Type", "multipart/form-data").
		SetMultipartField("file", "blob", "plain/text", reader).
		SetMultipartFormData(map[string]string{
			"spec": func() string {
				jsonData, _ := json.Marshal(body)
				return string(jsonData)
			}(),
		}).
		SetQueryParams(map[string]string{
			"accountIdentifier": sc.target.Account,
			"orgIdentifier":     sc.targetOrg,
			"projectIdentifier": sc.targetProject,
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

func (sc SecretContext) createSecretSSHKey(secret *nextgen.Secret) error {
	body := &model.CreateSecretRequest{
		Secret: secret,
	}

	api := sc.target
	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     sc.targetOrg,
			"projectIdentifier": sc.targetProject,
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
