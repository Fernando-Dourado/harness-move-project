package operation

import (
	"github.com/Fernando-Dourado/harness-move-project/services"
	"github.com/go-resty/resty/v2"
)

type (
	Config struct {
		Token   string
		Account string
	}

	// NOT SURE OF WHICH NAME TO CHOSE TO THAT TYPE
	NoName struct {
		Org     string
		Project string
	}

	Move struct {
		Config Config
		Source NoName
		Target NoName
	}
)

func (o *Move) Exec() error {

	api := services.ApiRequest{
		Client:  resty.New(),
		Token:   o.Config.Token,
		Account: o.Config.Account,
	}

	// SOURCE AND TARGET MUST EXIST
	if err := api.ValidateProject(o.Source.Org, o.Source.Project); err != nil {
		return err
	}
	if err := api.ValidateProject(o.Target.Org, o.Target.Project); err != nil {
		return err
	}

	// list pipelines
	// extract pipeline yamls

	return nil
}
