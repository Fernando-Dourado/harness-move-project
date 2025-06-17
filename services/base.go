package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/Fernando-Dourado/harness-move-project/model"
	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
)

// BaseURL is the default Harness API URL used when no specific URL is provided
const BaseURL = "https://app.harness.io"

var (
	ErrEntityNotFound = errors.New("entity not found")
)

type SourceRequest struct {
	Client  *resty.Client
	Token   string
	Account string
	Url     string
}

type TargetRequest struct {
	Client  *resty.Client
	Token   string
	Account string
	Url     string
}

type SourceTarget struct {
	SourceOrg     string
	SourceProject string
	TargetOrg     string
	TargetProject string
}

type Operation interface {
	Move() error
}

type OperationFactory interface {
	NewProjectOperation(sourceApi *SourceRequest, targetApi *TargetRequest, st *SourceTarget) ProjectContext

	NewVariableOperation(sourceApi *SourceRequest, targetApi *TargetRequest, st *SourceTarget) VariableContext
}

func createYaml(yaml, sourceOrg, sourceProject, targetOrg, targetProject string) string {
	var out string

	if strings.Contains(yaml, "orgIdentifier: ") {
		out = strings.ReplaceAll(yaml, "orgIdentifier: "+sourceOrg, "orgIdentifier: "+targetOrg)
	} else {
		out = fmt.Sprintln(yaml, " orgIdentifier:", targetOrg)
	}

	if strings.Contains(yaml, "projectIdentifier: ") {
		out = strings.ReplaceAll(out, "projectIdentifier: "+sourceProject, "projectIdentifier: "+targetProject)
	} else {
		if len(targetProject) > 0 {
			out = fmt.Sprintln(yaml, " projectIdentifier:", targetProject)
		}
	}

	return out
}

func handleErrorResponse(resp *resty.Response) error {
	result := model.ErrorResponse{}
	err := json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return err
	}
	if result.Code == "ENTITY_NOT_FOUND" {
		return ErrEntityNotFound
	}
	if result.Code == "DUPLICATE_FIELD" {
		return nil
	}
	if strings.Contains(result.Message, "already exists") {
		return nil
	}
	return fmt.Errorf("%s: %s", result.Code, removeNewLine(result.Message))
}

func removeNewLine(value string) string {
	return strings.ReplaceAll(value, "\n", "")
}

func reportFailed(failed []string, description string) {
	if len(failed) > 0 {
		fmt.Println(color.RedString(fmt.Sprintf("Failed %s %d", description, len(failed))))
		fmt.Println(color.RedString(strings.Join(failed, "\n")))
	}
}

func createQueryParams(account, org, project string) map[string]string {
	params := map[string]string{
		"accountIdentifier": account,
		"orgIdentifier":     org,
	}
	if len(project) > 0 {
		params["projectIdentifier"] = project
	}
	return params
}
