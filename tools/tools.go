//go:build tools
// +build tools

package tools

import (
	_ "github.com/Khan/genqlient/generate"
	_ "github.com/hashicorp/copywrite"
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
)
