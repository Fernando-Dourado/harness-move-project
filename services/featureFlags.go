package services

import (
	"encoding/json"
	"fmt"

	"github.com/jf781/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
)

const FEATFLAGS = "/cf/admin/features"

type FeatureContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
}

func NewFeatureFlagOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string) FeatureContext {
	return FeatureContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
	}
}

func (c FeatureContext) Move() error {

	// featureFlags := []*model.FeatureFlag{}

	featureFlags, err := c.api.listFeatureFlags(c.sourceOrg, c.sourceProject)
	if err != nil {
		return err
	}

	bar := progressbar.Default(int64(len(featureFlags)), "Feature Flags    ")
	var failed []string

	for _, f := range featureFlags {
		if f.Tags == nil {
			f.Tags = []string{}
		}

		err = c.api.createFeatureFlags(&model.CreateFeatureFlag{
			OrgIdentifier:       c.targetOrg,
			ProjectIdentifier:   c.targetProject,
			Archived:            f.Archived,
			CreatedAt:           f.CreatedAt,
			DefaultOffVariation: f.DefaultOffVariation,
			DefaultOnVariation:  f.DefaultOnVariation,
			Description:         f.Description,
			EnvProperties:       f.EnvProperties,
			Evaluation:          f.Evaluation,
			Identifier:          f.Identifier,
			Kind:                f.Kind,
			ModifiedAt:          f.ModifiedAt,
			Name:                f.Name,
			Owner:               fmt.Sprint(f.Owner),
			Permanent:           f.Permanent,
			Prerequisites:       f.Prerequisites,
			Project:             c.targetProject,
			Services:            f.Services,
			Tags:                f.Tags,
			Variations:          f.Variations,
		})

		if err != nil {
			failed = append(failed, fmt.Sprintln(f.Name, "-", err.Error()))
		}
		bar.Add(1)
	}
	bar.Finish()

	reportFailed(failed, "Feature Flags:")
	return nil
}

func (api *ApiRequest) listFeatureFlags(org, project string) ([]*model.FeatureFlag, error) {

	// api.Client.SetDebug(true)

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     org,
			"projectIdentifier": project,
		}).
		Get(BaseURL + FEATFLAGS)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse(resp)
	}

	result := model.FeatureFlagListResult{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	featureFlags := []*model.FeatureFlag{}
	for _, c := range result.Features {
		feature := c
		featureFlags = append(featureFlags, &feature)
	}

	return featureFlags, nil
}

func (api *ApiRequest) createFeatureFlags(featureFlag *model.CreateFeatureFlag) error {

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(featureFlag).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     featureFlag.OrgIdentifier,
			"projectIdentifier": featureFlag.ProjectIdentifier,
		}).
		Post(BaseURL + FEATFLAGS)

	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}

	return nil
}
