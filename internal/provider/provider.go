package provider

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/labd/terraform-provider-hive/internal/sdk"
)

// Ensure HiveProvider satisfies various provider interfaces.
var _ provider.Provider = &HiveProvider{}
var _ provider.ProviderWithFunctions = &HiveProvider{}
var _ provider.ProviderWithEphemeralResources = &HiveProvider{}

// HiveProvider defines the provider implementation.
type HiveProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// HiveProviderModel describes the provider data model.
type HiveProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Token    types.String `tfsdk:"token"`
}

func (p *HiveProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "hive"
	resp.Version = p.version
}

func (p *HiveProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "The endpoint of the Hive API",
				Optional:            true,
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "The token to authenticate with the registry",
				Required:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *HiveProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data HiveProviderModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "https://app.graphql-hive.com/graphql"
	if !data.Endpoint.IsNull() {
		endpoint = data.Endpoint.ValueString()
	}

	token := os.Getenv("HIVE_TOKEN")
	if !data.Token.IsNull() {
		token = data.Token.ValueString()
	}

	tflog.Info(ctx, fmt.Sprintf("Configuring Hive provider with endpoint: %s", endpoint))

	httpClient := &http.Client{
		Transport: sdk.NewDebugTransport(http.DefaultTransport),
	}
	client := sdk.NewHiveClient(httpClient, endpoint, token)

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *HiveProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewHiveSchemaCheckResource,
		NewHiveSchemaPublishResource,
		NewHiveAppCreateResource,
		NewHiveAppPublishResource,
	}
}

func (p *HiveProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *HiveProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *HiveProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &HiveProvider{
			version: version,
		}
	}
}
