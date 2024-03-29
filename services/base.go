package services

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Fernando-Dourado/harness-move-project/model"
	"github.com/go-resty/resty/v2"
)

const BaseURL = "https://app.harness.io"

type ApiRequest struct {
	Client  *resty.Client
	Token   string
	Account string
}

type Operation interface {
	Move() error
}

func createYaml(yaml, sourceOrg, sourceProject, targetOrg, targetProject string) string {
	var out string
	out = strings.ReplaceAll(yaml, "orgIdentifier: "+sourceOrg, "orgIdentifier: "+targetOrg)
	out = strings.ReplaceAll(out, "projectIdentifier: "+sourceProject, "projectIdentifier: "+targetProject)
	return out
}

func handleCreateErrorResponse(resp *resty.Response) error {
	result := model.ErrorResponse{}
	err := json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return err
	}
	if result.Code == "DUPLICATE_FIELD" {
		return nil
	}
	if strings.Contains(result.Message, "already exists") {
		return nil
	}
	return fmt.Errorf("%s: %s", result.Code, result.Message)
}
