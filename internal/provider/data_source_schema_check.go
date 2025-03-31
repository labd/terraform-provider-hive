package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/labd/terraform-provider-hive/internal/sdk"
)

var _ datasource.DataSource = &HiveSchemaCheckDataSource{}

func NewHiveSchemaCheckDataSource() datasource.DataSource {
	return &HiveSchemaCheckDataSource{}
}

type HiveSchemaCheckDataSource struct {
	client *sdk.HiveClient
}

func (r *HiveSchemaCheckDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schema_check"
}

type HiveSchemaCheckDataSourceModel struct {
	Service types.String `tfsdk:"service"`
	Commit  types.String `tfsdk:"commit"`
	Author  types.String `tfsdk:"author"`
	Schema  types.String `tfsdk:"schema"`
	Id      types.String `tfsdk:"id"`
	Target  types.String `tfsdk:"target"`
	Project types.String `tfsdk:"project"`
}

func (d *HiveSchemaCheckDataSource) Schema(ctx context.Context, _req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source to perform a schema check against a GraphQL schema",

		Attributes: map[string]schema.Attribute{
			"service": schema.StringAttribute{
				MarkdownDescription: "The service name",
				Required:            true,
			},
			"commit": schema.StringAttribute{
				MarkdownDescription: "The commit or version identifier",
				Optional:            true,
			},
			"schema": schema.StringAttribute{
				MarkdownDescription: "The GraphQL schema content",
				Required:            true,
			},
			"author": schema.StringAttribute{
				MarkdownDescription: "The author of the version",
				Optional:            true,
			},
			"project": schema.StringAttribute{
				MarkdownDescription: "The project name",
				Optional:            true,
			},
			"target": schema.StringAttribute{
				MarkdownDescription: "The target name",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The resource ID",
			},
		},
	}
}

func (r *HiveSchemaCheckDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sdk.HiveClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *sdk.HiveClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *HiveSchemaCheckDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data HiveSchemaCheckDataSourceModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.SchemaCheck(ctx, &sdk.SchemaCheckInput{
		Service: data.Service.ValueString(),
		Schema:  data.Schema.ValueString(),
		Commit:  data.Commit.ValueString(),
		Author:  data.Author.ValueString(),
		Target:  data.Target.ValueString(),
		Project: data.Project.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Schema check failed",
			fmt.Sprintf("Unable validate schema, got error: %s", err.Error()),
		)
		return
	}

	if result == nil {
		resp.Diagnostics.AddError("Schema check failed", "The schema is not valid")
		return
	}

	if !result.Valid {
		resp.Diagnostics.AddError("Schema check failed", fmt.Sprintf("The schema is not valid, see %s for more details", result.URL))
		return
	}

	data.Id = types.StringValue(result.Id)

	// Save any updates back to state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
