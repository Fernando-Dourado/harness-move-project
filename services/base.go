package services

import "github.com/go-resty/resty/v2"

const BaseURL = "https://app.harness.io"

type ApiRequest struct {
	Client  *resty.Client
	Token   string
	Account string
}
