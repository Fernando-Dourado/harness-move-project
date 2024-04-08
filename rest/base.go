package rest

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Fernando-Dourado/harness-move-project/model"
	"github.com/go-resty/resty/v2"
)

const BaseURL_ = "https://app.harness.io"

type ApiContext struct {
	Client  *resty.Client
	Token   string
	Account string
}

func handleErrorResponse_(resp *resty.Response) error {
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
	return fmt.Errorf("%s: %s", result.Code, strings.ReplaceAll(result.Message, "\n", ""))
}
