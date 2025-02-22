package sdk

import (
	"context"
	"fmt"

	"github.com/labd/terraform-provider-hive/internal/client"
)

type PublishAppInput struct {
	Name      string
	Version   string
	Documents string
}

type PublishAppResult struct {
	Id         string
	AppName    string
	AppVersion string
	Status     string
}

/**
 * PublishApp() activates the app deployment of the given name and version.
 */
func (hc *HiveClient) PublishApp(ctx context.Context, input *PublishAppInput) (*PublishAppResult, error) {

	data, err := client.ActivateAppDeployment(ctx, *hc.client, client.ActivateAppDeploymentInput{
		AppName:    input.Name,
		AppVersion: input.Version,
	})

	if err != nil {
		return nil, err
	}

	if data.ActivateAppDeployment.GetError() != nil {
		return nil, fmt.Errorf("failed to create app: %s", data.ActivateAppDeployment.GetError().Message)
	}

	result := PublishAppResult{
		Id:         data.ActivateAppDeployment.GetOk().ActivatedAppDeployment.GetId(),
		AppName:    data.ActivateAppDeployment.GetOk().ActivatedAppDeployment.GetName(),
		AppVersion: data.ActivateAppDeployment.GetOk().ActivatedAppDeployment.GetVersion(),
		Status:     string(data.ActivateAppDeployment.GetOk().ActivatedAppDeployment.GetStatus()),
	}

	return &result, nil
}
