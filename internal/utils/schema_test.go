package utils

import (
	"testing"
)

func TestCleanSchema(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  string
		expectErr bool
	}{
		{
			name:     "simple type",
			input:    "type Query { hello: String }",
			expected: "type Query {\nhello: String\n}\n",
		},
		{
			name:     "normalizes whitespace",
			input:    "type   Query  {   hello:   String   }",
			expected: "type Query {\nhello: String\n}\n",
		},
		{
			name:     "multiple fields",
			input:    "type Query { hello: String world: Int }",
			expected: "type Query {\nhello: String\nworld: Int\n}\n",
		},
		{
			name:     "schema with description",
			input:    `"A simple query" type Query { hello: String }`,
			expected: "\"\"\"\nA simple query\n\"\"\"\ntype Query {\nhello: String\n}\n",
		},
		{
			name: "schema with comment blocks",
			input: `
			###
			#Test
			### 
			type Query { hello: String }`,
			expected: "type Query {\nhello: String\n}\n",
		},
		{
			name:     "empty string returns newline",
			input:    "",
			expected: "",
		},
		{
			name:     "type with schema extension and directive",
			input:    "extend schema\n@link(url: \"https://specs.apollo.dev/federation/v2.0\", import: [\"@key\"])\ntype Query { hello: String }",
			expected: "extend schema @link(url: \"https://specs.apollo.dev/federation/v2.0\", import: [\"@key\"])\ntype Query {\nhello: String\n}\n",
		},

		{
			name:      "invalid schema",
			input:     "not a valid schema {{{",
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := CleanSchema(tc.input)

			if tc.expectErr {
				if err == nil {
					t.Fatalf("expected error, got nil (result: %q)", result)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result != tc.expected {
				t.Errorf("expected:\n%q\ngot:\n%q", tc.expected, result)
			}
		})
	}
}
