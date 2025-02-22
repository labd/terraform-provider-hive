package sdk

import (
	"context"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type SchemaPublishInput struct {
	Service string
	Schema  string
	URL     string
	Commit  string
}

// SchemaPublishResponse represents the main schemaCheck response.
type SchemaPublishResponse struct {
	Typename      string     `json:"__typename"`
	Id            string     `json:"id"`
	Valid         bool       `json:"valid"`
	Initial       bool       `json:"initial"`
	LinkToWebsite string     `json:"linkToWebsite"`
	Changes       CountNodes `json:"changes"`
}

// Response is the top-level structure.
type SchemaPublishGraphQLResponse struct {
	Data struct {
		SchemaPublish SchemaPublishResponse `json:"schemaPublish"`
	} `json:"data"`
}

func (hc *HiveClient) SchemaPublish(ctx context.Context, input *SchemaPublishInput) (*SchemaPublishResponse, error) {
	query := `
		mutation schemaPublish($input: SchemaPublishInput!, $usesGitHubApp: Boolean!) {
			schemaPublish(input: $input) {
				__typename
				... on SchemaPublishSuccess @skip(if: $usesGitHubApp) {
					initial
					valid
					successMessage: message
					linkToWebsite
					changes {
						nodes {
							message(withSafeBasedOnUsageNote: false)
							criticality
							isSafeBasedOnUsage
						}
						total
						...RenderChanges_schemaChanges
					}
				}
				... on SchemaPublishError @skip(if: $usesGitHubApp) {
					valid
					linkToWebsite
					changes {
						nodes {
							message(withSafeBasedOnUsageNote: false)
							criticality
							isSafeBasedOnUsage
						}
						total
						...RenderChanges_schemaChanges
					}
					errors {
						nodes {
							message
						}
						total
					}
				}
				... on SchemaPublishMissingServiceError @skip(if: $usesGitHubApp) {
					missingServiceError: message
				}
				... on SchemaPublishMissingUrlError @skip(if: $usesGitHubApp) {
					missingUrlError: message
				}
				... on GitHubSchemaPublishSuccess @include(if: $usesGitHubApp) {
					message
				}
				... on GitHubSchemaPublishError @include(if: $usesGitHubApp) {
					message
				}
				... on SchemaPublishRetry {
					reason
				}
			}
		}

		fragment RenderChanges_schemaChanges on SchemaChangeConnection {
			total
			nodes {
				criticality
				isSafeBasedOnUsage
				message(withSafeBasedOnUsageNote: false)
				approval {
					approvedBy {
						displayName
					}
				}
			}
		}
	`

	result := SchemaPublishGraphQLResponse{}

	err := hc.Execute(ctx, query, map[string]any{
		"input": map[string]any{
			"service": input.Service,
			"commit":  input.Commit,
			"author":  "Michael van Tellingen <m.vantellingen@labdigital.nl>",
			"sdl":     minifySchema(input.Schema),
			"url":     input.URL,
		},
		"usesGitHubApp": false,
	},
		&result,
	)

	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("error: %v", err))
		return nil, err
	}

	spew.Dump(result.Data.SchemaPublish)
	rv := &result.Data.SchemaPublish
	rv.Id = rv.LinkToWebsite
	return rv, nil

}
