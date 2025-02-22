package sdk

import (
	"context"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type CreateAppInput struct {
	Name    string
	Version string
}

type CreateAppResponse struct {
	Ok    *CreateAppResponseOk    `json:"ok"`
	Error *CreateAppResponseError `json:"error"`
}

type CreatedAppDeployment struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Version string `json:"version"`
	Status  string `json:"status"`
}

type CreateAppResponseOk struct {
	CreatedAppDeployment CreatedAppDeployment `json:"createdAppDeployment"`
}

type CreateAppResponseError struct {
	Message string `json:"message"`
}

type CreateAppGraphQLResponse struct {
	Data struct {
		CreateAppDeployment CreateAppResponse `json:"createAppDeployment"`
	} `json:"data"`
}

/**
 * CreateApp() first creates a new app version and then pushes all docucuments
 * (batched) to the new version.
 */
func (hc *HiveClient) CreateApp(ctx context.Context, input *CreateAppInput) (*CreateAppResponse, error) {
	query := `
		mutation CreateAppDeployment($input: CreateAppDeploymentInput!) {
			createAppDeployment(input: $input) {
				ok {
					createdAppDeployment {
						id
						name
						version
						status
					}
				}
				error {
					message
				}
			}
		}
	`

	result := CreateAppGraphQLResponse{}

	err := hc.Execute(ctx, query, map[string]any{
		"input": map[string]any{
			"appName":    input.Name,
			"appVersion": input.Version,
		},
	},
		&result,
	)

	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("error: %v", err))
		return nil, err
	}

	tflog.Debug(ctx, fmt.Sprintf("result: %s", spew.Sdump(result)))

	return &result.Data.CreateAppDeployment, nil

}
