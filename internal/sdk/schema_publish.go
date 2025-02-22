package sdk

import (
	"context"
	"fmt"
	"net/url"
	"path"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/labd/terraform-provider-hive/internal/client"
)

type SchemaPublishInput struct {
	Service string
	Schema  string
	URL     string
	Commit  string
}

type SchemaPublishResult struct {
	Valid bool
	URL   string
	Id    string
}

func (hc *HiveClient) SchemaPublish(ctx context.Context, input *SchemaPublishInput) (*SchemaPublishResult, error) {

	data, err := client.SchemaPublish(ctx, *hc.client, client.SchemaPublishInput{
		Service: input.Service,
		Commit:  input.Commit,
		Author:  "Michael van Tellingen <m.vantellingen@labdigital.nl>",
		Sdl:     minifySchema(input.Schema),
		Url:     input.URL,
	}, false)

	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("error: %v", err))
		return nil, err
	}

	switch v := data.SchemaPublish.(type) {

	case *client.SchemaPublishSchemaPublishSchemaPublishSuccess:
		result := SchemaPublishResult{
			Valid: v.Valid,
			URL:   v.GetLinkToWebsite(),
			Id:    extractIdFromURL(v.GetLinkToWebsite()),
		}
		return &result, nil

	case *client.SchemaPublishSchemaPublishSchemaPublishError:
		result := SchemaPublishResult{
			Valid: v.Valid,
			URL:   v.GetLinkToWebsite(),
			Id:    extractIdFromURL(v.GetLinkToWebsite()),
		}
		return &result, nil

	case *client.SchemaPublishSchemaPublishGitHubSchemaPublishSuccess:
		result := SchemaPublishResult{
			Valid: true,
		}
		return &result, nil

	case *client.SchemaPublishSchemaPublishSchemaPublishMissingServiceError:
		return nil, fmt.Errorf("hive error: %s", v.GetMessage())

	case *client.SchemaPublishSchemaPublishSchemaPublishMissingUrlError:
		return nil, fmt.Errorf("hive error: %s", v.GetMessage())

	case *client.SchemaPublishSchemaPublishGitHubSchemaPublishError:
		return nil, fmt.Errorf("hive error: %s", v.GetMessage())
	}

	return nil, fmt.Errorf("unexpected type %T", data.SchemaPublish)
}

// Hive doesn't return the check id for this mutation, so we just extract it
// from the URL (which it does return).
func extractIdFromURL(value string) string {
	u, err := url.Parse(value)
	if err != nil {
		return ""
	}

	// Use path.Base to extract the last segment of the URL path.
	return path.Base(u.Path)
}
