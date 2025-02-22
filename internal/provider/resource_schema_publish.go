package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/labd/terraform-provider-hive/internal/sdk"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &HiveSchemaPublishResource{}
var _ resource.ResourceWithImportState = &HiveSchemaPublishResource{}

// NewHiveSchemaPublishResource is a helper function to simplify the provider implementation.
func NewHiveSchemaPublishResource() resource.Resource {
	return &HiveSchemaPublishResource{}
}

// HiveSchemaPublishResource defines the resource implementation.
type HiveSchemaPublishResource struct {
	client *sdk.HiveClient
}

// HiveSchemaPublishResourceModel describes the resource data model.
type HiveSchemaPublishResourceModel struct {
	Service types.String `tfsdk:"service"`
	Commit  types.String `tfsdk:"commit"`
	Schema  types.String `tfsdk:"schema"`
	URL     types.String `tfsdk:"url"`
	Id      types.String `tfsdk:"id"`
}

// Metadata returns the resource type name.
func (r *HiveSchemaPublishResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schema_publish"
}

// Schema defines the schema for the hive_schema_publish resource.
func (r *HiveSchemaPublishResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource to perform a schema publish for a GraphQL schema",

		Attributes: map[string]schema.Attribute{
			"service": schema.StringAttribute{
				MarkdownDescription: "The service name",
				Required:            true,
			},
			"commit": schema.StringAttribute{
				MarkdownDescription: "The commit or version identifier",
				Required:            true,
			},
			"schema": schema.StringAttribute{
				MarkdownDescription: "The GraphQL schema content",
				Required:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "The GraphQL schema content",
				Required:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The resource ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Configure saves the provider configured HTTP client on the resource.
func (r *HiveSchemaPublishResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sdk.HiveClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *sdk.HiveClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Create handles the creation of the resource.
func (r *HiveSchemaPublishResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data HiveSchemaPublishResourceModel

	// Read Terraform plan data into the model.
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	diag := r.ExecuteRequest(ctx, &data)
	if diag != nil {
		resp.Diagnostics.Append(*diag)
		return
	}

	// Save the data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data. This is a no-op for this resource.
func (r *HiveSchemaPublishResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data HiveSchemaPublishResourceModel

	// Retrieve state.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save any updates back to state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update handles updates to the resource.
func (r *HiveSchemaPublishResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data HiveSchemaPublishResourceModel

	// Read Terraform plan data into the model.
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	diag := r.ExecuteRequest(ctx, &data)
	if diag != nil {
		resp.Diagnostics.Append(*diag)
		return
	}

	// Save updated data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete handles resource deletion. This is a no-op for this resource.
func (r *HiveSchemaPublishResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

// ImportState allows the resource to be imported into Terraform.
func (r *HiveSchemaPublishResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *HiveSchemaPublishResource) ExecuteRequest(ctx context.Context, data *HiveSchemaPublishResourceModel) *diag.ErrorDiagnostic {
	result, err := r.client.SchemaPublish(ctx, &sdk.SchemaPublishInput{
		Service: data.Service.ValueString(),
		Schema:  data.Schema.ValueString(),
		Commit:  data.Commit.ValueString(),
		URL:     data.URL.ValueString(),
	})

	if err != nil {
		d := diag.NewErrorDiagnostic("Schema publish failed", err.Error())
		return &d
	}

	if !result.Valid {
		d := diag.NewErrorDiagnostic("Schema publish failed", fmt.Sprintf("The schema is not valid, see %s for more details", result.URL))
		return &d
	}

	data.Id = types.StringValue(result.Id)

	return nil
}
