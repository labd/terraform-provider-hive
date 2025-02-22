package sdk

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type SchemaCheckInput struct {
	Service string
	Schema  string
}

// CountNodes represents the warnings/changes object.
type CountNodes struct {
	Nodes []interface{} `json:"nodes"`
	Total int           `json:"total"`
}

// SchemaCheckDetail represents the inner schemaCheck details.
type SchemaCheckDetail struct {
	ID     string `json:"id"`
	WebURL string `json:"webUrl"`
}

// SchemaCheckResponse represents the main schemaCheck response.
type SchemaCheckResponse struct {
	SchemaCheck SchemaCheckDetail `json:"schemaCheck"`
	Typename    string            `json:"__typename"`
	Valid       bool              `json:"valid"`
	Initial     bool              `json:"initial"`
	Warnings    CountNodes        `json:"warnings"`
	Changes     CountNodes        `json:"changes"`
}

// Response is the top-level structure.
type SchemaCheckGraphQLResponse struct {
	Data struct {
		SchemaCheck SchemaCheckResponse `json:"schemaCheck"`
	} `json:"data"`
}

func (hc *HiveClient) SchemaCheck(ctx context.Context, input *SchemaCheckInput) (*SchemaCheckResponse, error) {
	query := `
		mutation schemaCheck($input: SchemaCheckInput!) {
			schemaCheck(input: $input) {
				__typename
				... on SchemaCheckSuccess {
					valid
					initial
					warnings {
						nodes {
							message
							source
							line
							column
						}
						total
					}
					changes {
						nodes {
							message(withSafeBasedOnUsageNote: false)
							criticality
							isSafeBasedOnUsage
							approval {
								approvedBy {
									id
									displayName
								}
							}
						}
						total
						...RenderChanges_schemaChanges
					}
					schemaCheck {
						id
						webUrl
					}
				}
				... on SchemaCheckError {
					valid
					changes {
						nodes {
							message(withSafeBasedOnUsageNote: false)
							criticality
							isSafeBasedOnUsage
						}
						total
						...RenderChanges_schemaChanges
					}
					warnings {
						nodes {
							message
							source
							line
							column
						}
						total
					}
					errors {
						nodes {
							message
						}
						total
					}
					schemaCheck {
						webUrl
					}
				}
				... on GitHubSchemaCheckSuccess {
					message
				}
				... on GitHubSchemaCheckError {
					message
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

	result := SchemaCheckGraphQLResponse{}

	err := hc.Execute(ctx, query, map[string]any{
		"input": map[string]any{
			"service": input.Service,
			"sdl":     minifySchema(input.Schema),
		},
	},
		&result,
	)

	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("error: %v", err))
		return nil, err
	}

	return &result.Data.SchemaCheck, nil

}
