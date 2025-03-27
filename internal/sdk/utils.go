package sdk

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/labd/terraform-provider-hive/internal/client"
)

// minifySchema removes extra whitespace from the schema string.
func minifySchema(schema string) string {
	re := regexp.MustCompile(`\s+`)
	return strings.TrimSpace(re.ReplaceAllString(schema, " "))
}

func getTarget(ctx context.Context, input string) (*client.TargetReferenceInput, error)  {
	if input == "" {
        return nil, nil
    }

	split := strings.Split(input, "/")

 	if (len(split) != 3) {
		return nil, fmt.Errorf("not a valid target: %v", input)

	}

	return &client.TargetReferenceInput{
		BySelector: client.TargetSelectorInput{
			OrganizationSlug: split[0],
			ProjectSlug:      split[1],
			TargetSlug:       split[2],
		},
	}, nil
}