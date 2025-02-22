package sdk

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/labd/terraform-provider-hive/internal/client"
)

type SchemaCheckInput struct {
	Service string
	Schema  string
}

type SchemaCheckResult struct {
	Id    string
	Valid bool
	URL   string
}

func (hc *HiveClient) SchemaCheck(ctx context.Context, input *SchemaCheckInput) (*SchemaCheckResult, error) {

	data, err := client.SchemaCheck(ctx, *hc.client, client.SchemaCheckInput{
		Service: input.Service,
		Sdl:     minifySchema(input.Schema),
		// Meta: client.SchemaCheckMetaInput{
		// 	Author: "Michael",
		// 	Commit: "123456",
		// },
		// Target: client.TargetReferenceInput{
		// 	BySelector: client.TargetSelectorInput{
		// 		OrganizationSlug: "..",
		// 		ProjectSlug:      "michael-sandbox",
		// 		TargetSlug:       "development",
		// 	},
		// },
	})

	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("error: %v", err))
		return nil, err
	}

	switch v := data.SchemaCheck.(type) {
	case *client.SchemaCheckSchemaCheckSchemaCheckSuccess:
		result := SchemaCheckResult{
			Valid: v.Valid,
			Id:    v.SchemaCheck.GetId(),
			URL:   v.SchemaCheck.GetWebUrl(),
		}
		return &result, nil

	case *client.SchemaCheckSchemaCheckSchemaCheckError:
		result := SchemaCheckResult{
			Valid: v.Valid,
			Id:    v.SchemaCheck.GetId(),
			URL:   v.SchemaCheck.GetWebUrl(),
		}
		return &result, nil

	case *client.SchemaCheckSchemaCheckGitHubSchemaCheckSuccess:
		result := SchemaCheckResult{Valid: true}
		return &result, nil

	case *client.SchemaCheckSchemaCheckGitHubSchemaCheckError:
		result := SchemaCheckResult{Valid: false}
		return &result, nil

	}

	return nil, fmt.Errorf("unexpected type %T", data.SchemaCheck)
}
