package utils

import (
	"fmt"

	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/printer"
)

// CleanSchema parses and re-prints a GraphQL schema to normalize formatting.
func CleanSchema(schema string) (string, error) {
	docs, err := parser.Parse(parser.ParseParams{
		Source: schema,
	})
	if err != nil {
		return "", err
	}

	cleanedSchema := printer.Print(docs)

	result, ok := cleanedSchema.(string)
	if !ok {
		return "", fmt.Errorf("expected string from printer, got %T", cleanedSchema)
	}

	return result, nil
}
