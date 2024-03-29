package services

import (
	"strings"

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

type OperationContext struct {
}

func createYaml(yaml, sourceOrg, sourceProject, targetOrg, targetProject string) string {
	var out string
	out = strings.ReplaceAll(yaml, "orgIdentifier: "+sourceOrg, "orgIdentifier: "+targetOrg)
	out = strings.ReplaceAll(out, "projectIdentifier: "+sourceProject, "projectIdentifier: "+targetProject)
	return out
}
