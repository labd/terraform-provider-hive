package sdk

import (
	"regexp"
	"strings"
)

// minifySchema removes extra whitespace from the schema string.
func minifySchema(schema string) string {
	re := regexp.MustCompile(`\s+`)
	return strings.TrimSpace(re.ReplaceAllString(schema, " "))
}
