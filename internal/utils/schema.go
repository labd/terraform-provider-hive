package utils

import (
	"bytes"

	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/formatter"
	"github.com/vektah/gqlparser/v2/parser"
)

// CleanSchema parses and re-prints a GraphQL schema to normalize formatting.
func CleanSchema(schema string) (string, error) {
	doc, err := parser.ParseSchema(&ast.Source{
		Input: schema,
	})
	if err != nil {
		return "", err
	}

	// exec format
	var buf bytes.Buffer
	formatter.NewFormatter(&buf, formatter.WithCompacted(), formatter.WithIndent("")).FormatSchemaDocument(doc)

	// validity check
	_, err = parser.ParseSchema(&ast.Source{
		Input: buf.String(),
	})
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
