package sdk

import (
	"context"
	"regexp"
	"strings"

	"github.com/labd/terraform-provider-hive/internal/client"
)

// minifySchema removes extra whitespace from the schema string.
func minifySchema(schema string) string {
	re := regexp.MustCompile(`\s+`)
	return strings.TrimSpace(re.ReplaceAllString(schema, " "))
}

func getTarget(ctx context.Context, organization string, project string, target string) *client.TargetReferenceInput {
	if organization == "" || project == "" || target == "" {
		return nil
	}

	return &client.TargetReferenceInput{
		BySelector: client.TargetSelectorInput{
			OrganizationSlug: organization,
			ProjectSlug:      project,
			TargetSlug:       target,
		},
	}
}
