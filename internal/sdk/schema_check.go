package sdk

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/labd/terraform-provider-hive/internal/client"
)

type SchemaCheckInput struct {
	Service   string
	Schema    string
	Author    string
	Commit    string
	ContextId string
	Target    string
	Project   string
}

type SchemaCheckResult struct {
	Id    string
	Valid bool
	URL   string
}

func (hc *HiveClient) SchemaCheck(ctx context.Context, input *SchemaCheckInput) (*SchemaCheckResult, error) {

	meta := &client.SchemaCheckMetaInput{
		Author: input.Author,
		Commit: input.Commit,
	}

	if meta.Author == "" || meta.Commit == "" {
		gitInfo, err := GetLatestCommitInfo()
		if err == nil {
			if meta.Author == "" {
				meta.Author = gitInfo.Author
			}

			if meta.Commit == "" {
				meta.Commit = gitInfo.Hash
			}
		}
	}

	vars := client.SchemaCheckInput{
		Service:   input.Service,
		Sdl:       minifySchema(input.Schema),
		Meta:      meta,
		ContextId: input.ContextId,
		Target:    getTarget(ctx, hc.Organization, input.Project, input.Target),
	}

	data, err := client.SchemaCheck(ctx, *hc.client, vars)

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
