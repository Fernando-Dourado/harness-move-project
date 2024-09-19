package services

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"github.com/jf781/harness-move-project/model"
)

// const BaseURL = "https://app.harness.io"

type ApiRequest struct {
	Client  *resty.Client
	Token   string
	Account string
	BaseURL string
}

type Operation interface {
	Copy() error
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
		out = fmt.Sprintln(yaml, " projectIdentifier:", targetProject)
	}

	return out
}

func createYamlQuotes(yaml, sourceOrg, sourceProject, targetOrg, targetProject string) string {
	var out string

	if strings.Contains(yaml, "orgIdentifier: ") {
		out = strings.ReplaceAll(yaml, "orgIdentifier: \""+sourceOrg+"\"", "orgIdentifier: \""+targetOrg+"\"")
	} else {
		out = fmt.Sprintln(yaml, " orgIdentifier:", targetOrg)
	}

	if strings.Contains(yaml, "projectIdentifier: ") {
		out = strings.ReplaceAll(out, "projectIdentifier: \""+sourceProject, "projectIdentifier: \""+targetProject)
	} else {
		out = fmt.Sprintln(yaml, " projectIdentifier:", targetProject)
	}

	return out
}

func handleErrorResponse(resp *resty.Response) error {
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
