package sdk

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/elliotchance/pie/v2"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/labd/terraform-provider-hive/internal/client"
)

type CreateAppInput struct {
	Name      string
	Version   string
	Documents string
}

type CreateAppResult struct {
	Id         string
	AppName    string
	AppVersion string
	Status     string
}

func parseDocuments(documents string) ([]client.DocumentInput, error) {
	operations := map[string]string{}
	err := json.Unmarshal([]byte(documents), &operations)
	if err != nil {
		return nil, err
	}

	result := make([]client.DocumentInput, 0, len(operations))
	for key, value := range operations {
		result = append(result, client.DocumentInput{
			Body: value,
			Hash: key,
		})
	}

	return result, nil
}

/**
 * CreateApp() first creates a new app version and then pushes all docucuments
 * (batched) to the new version.
 */
func (hc *HiveClient) CreateApp(ctx context.Context, input *CreateAppInput) (*CreateAppResult, error) {
	documents, err := parseDocuments(input.Documents)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal documents: %w", err)
	}

	if len(documents) == 0 {
		return nil, fmt.Errorf("no operations found in documents")
	}

	data, err := client.CreateAppDeployment(ctx, *hc.client, client.CreateAppDeploymentInput{
		AppName:    input.Name,
		AppVersion: input.Version,
	})

	if err != nil {
		return nil, err
	}

	if data.CreateAppDeployment.GetError() != nil {
		return nil, fmt.Errorf("failed to create app: %s", data.CreateAppDeployment.GetError().Message)
	}

	for _, batch := range pie.Chunk(documents, 100) {

		data, err := client.AddDocumentsToAppDeployment(ctx, *hc.client, client.AddDocumentsToAppDeploymentInput{
			AppName:    input.Name,
			AppVersion: input.Version,
			Documents:  batch,
		})

		if err != nil {
			return nil, err
		}

		result := data.GetAddDocumentsToAppDeployment()

		if result.Error != nil {

			// Skip this error for now. Need to investigate this further.
			if result.Error.Message != "App deployment has already been activated and is locked for modifications" {
				return nil, fmt.Errorf("failed to add documents: %s", result.Error.Message)
			} else {
				tflog.Debug(ctx, spew.Sdump(result))
			}
		}
	}

	result := CreateAppResult{
		Id:         data.CreateAppDeployment.GetOk().CreatedAppDeployment.GetId(),
		AppName:    data.CreateAppDeployment.GetOk().CreatedAppDeployment.GetName(),
		AppVersion: data.CreateAppDeployment.GetOk().CreatedAppDeployment.GetVersion(),
		Status:     string(data.CreateAppDeployment.GetOk().CreatedAppDeployment.GetStatus()),
	}

	return &result, nil
}
