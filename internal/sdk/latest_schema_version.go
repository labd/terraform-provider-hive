package sdk

import (
	"context"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (hc *HiveClient) GetLatestSchemaVersion(ctx context.Context) {
	query := `
		query LatestSchemaVersion(
			$includeSDL: Boolean!
			$includeSupergraph: Boolean!
			$includeSubgraphs: Boolean!
			$target: TargetReferenceInput
		) {
			latestValidVersion(target: $target) {
				id
				valid
				sdl @include(if: $includeSDL)
				supergraph @include(if: $includeSupergraph)
				schemas @include(if: $includeSubgraphs) {
					nodes {
						__typename
						... on SingleSchema {
							id
							date
						}
						... on CompositeSchema {
							id
							date
							url
							service
						}
					}
					total
				}
			}
		}
	`

	result := map[string]any{}
	err := hc.Execute(ctx, query, map[string]interface{}{
		"includeSDL":        true,
		"includeSupergraph": true,
		"includeSubgraphs":  false,
	}, &result)
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("error: %v", err))
		return
	}

	spew.Dump(result)
	tflog.Debug(ctx, fmt.Sprintf("result: %v", result))

}
