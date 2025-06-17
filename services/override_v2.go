package services

import (
	"encoding/json"
	"fmt"

	"github.com/Fernando-Dourado/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
)

type OverrideV2Context struct {
	source        *SourceRequest
	target        *TargetRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
}

func NewOverrideV2Operation(sourceApi *SourceRequest, targetApi *TargetRequest, st *SourceTarget) OverrideV2Context {
	return OverrideV2Context{
		source:        sourceApi,
		target:        targetApi,
		sourceOrg:     st.SourceOrg,
		sourceProject: st.SourceProject,
		targetOrg:     st.TargetOrg,
		targetProject: st.TargetProject,
	}
}

func (c OverrideV2Context) Move() error {

	// FETCH AND CREATE BY TYPE
	overrideTypes := []model.OverridesV2Type{
		model.OV2_Global,
		model.OV2_Service,
		model.OV2_Infra,
		model.OV2_ServiceInfra,
	}

	bar := progressbar.Default(int64(len(overrideTypes)), "Overrides V2")
	var failed = []string{}

	for _, overrideType := range overrideTypes {
		overrideIds, err := c.listOverrides(c.sourceOrg, c.sourceProject, overrideType)
		if err != nil {
			return err
		}
		bar.ChangeMax(bar.GetMax() + len(overrideIds))

		for _, id := range overrideIds {
			override, err := c.getOverride(id)
			if err != nil {
				failed = append(failed, fmt.Sprintln("Unable to get override type", overrideType, "identifier", id, "-", err.Error()))
			} else {
				override.OrgIdentifier = c.targetOrg
				override.ProjectIdentifier = c.targetProject

				if err = c.createOverride(override); err != nil {
					failed = append(failed, fmt.Sprintln("Unable to create override type", override.Type, "identifier", override.Identifier, "-", err.Error()))
				}
			}
			bar.Add(1)
		}
		bar.Add(1)
	}
	bar.Finish()

	reportFailed(failed, "overrides v2:")
	return nil
}

func (c OverrideV2Context) listOverrides(org, project string, overrideType model.OverridesV2Type) ([]string, error) {

	params := createQueryParams(c.source.Account, org, project)
	params["size"] = "1000"
	params["page"] = "0"
	params["type"] = string(overrideType)

	api := c.source
	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(params).
		Post(api.Url + "/ng/api/serviceOverrides/v2/list")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse(resp)
	}

	result := model.ListOverridesV2Response{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	overrides := []string{}
	for _, c := range result.Data.Content {
		overrides = append(overrides, c.Identifier)
	}

	return overrides, nil
}

func (c OverrideV2Context) getOverride(identifier string) (*model.OverridesV2, error) {

	params := createQueryParams(c.source.Account, c.sourceOrg, c.sourceProject)

	api := c.source
	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(params).
		SetPathParams(map[string]string{
			"identifier": identifier,
		}).
		Get(api.Url + "/ng/api/serviceOverrides/{identifier}")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse(resp)
	}

	result := model.GetOverridesV2Response{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

func (c OverrideV2Context) createOverride(override *model.OverridesV2) error {

	api := c.target
	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(override).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
		}).
		Post(api.Url + "/ng/api/serviceOverrides")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}

	return nil
}
