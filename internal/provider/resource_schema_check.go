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
var _ resource.Resource = &HiveSchemaCheckResource{}
var _ resource.ResourceWithImportState = &HiveSchemaCheckResource{}

// NewHiveSchemaCheckResource is a helper function to simplify the provider implementation.
func NewHiveSchemaCheckResource() resource.Resource {
	return &HiveSchemaCheckResource{}
}

// HiveSchemaCheckResource defines the resource implementation.
type HiveSchemaCheckResource struct {
	client *sdk.HiveClient
}

// HiveSchemaCheckResourceModel describes the resource data model.
type HiveSchemaCheckResourceModel struct {
	Service   types.String `tfsdk:"service"`
	Commit    types.String `tfsdk:"commit"`
	Author    types.String `tfsdk:"author"`
	Schema    types.String `tfsdk:"schema"`
	ContextId types.String `tfsdk:"context_id"`
	Id        types.String `tfsdk:"id"`
}

// Metadata returns the resource type name.
func (r *HiveSchemaCheckResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schema_check"
}

// Schema defines the schema for the hive_schema_check resource.
func (r *HiveSchemaCheckResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource to perform a schema check against a GraphQL schema",

		Attributes: map[string]schema.Attribute{
			"service": schema.StringAttribute{
				MarkdownDescription: "The service name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"commit": schema.StringAttribute{
				MarkdownDescription: "The commit or version identifier",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"schema": schema.StringAttribute{
				MarkdownDescription: "The GraphQL schema content",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"author": schema.StringAttribute{
				MarkdownDescription: "The author of the version",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"context_id": schema.StringAttribute{
				MarkdownDescription: "Context ID allows retaining approved breaking changes with the lifecycle",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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
func (r *HiveSchemaCheckResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *HiveSchemaCheckResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data HiveSchemaCheckResourceModel

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
func (r *HiveSchemaCheckResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data HiveSchemaCheckResourceModel

	// Retrieve state.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save any updates back to state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update handles updates to the resource.
func (r *HiveSchemaCheckResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data HiveSchemaCheckResourceModel

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
func (r *HiveSchemaCheckResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

// ImportState allows the resource to be imported into Terraform.
func (r *HiveSchemaCheckResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *HiveSchemaCheckResource) ExecuteRequest(ctx context.Context, data *HiveSchemaCheckResourceModel) *diag.ErrorDiagnostic {
	result, err := r.client.SchemaCheck(ctx, &sdk.SchemaCheckInput{
		Service:   data.Service.ValueString(),
		Schema:    data.Schema.ValueString(),
		Commit:    data.Commit.ValueString(),
		Author:    data.Author.ValueString(),
		ContextId: data.ContextId.ValueString(),
	})

	if err != nil {
		d := diag.NewErrorDiagnostic("Schema check failed", err.Error())
		return &d
	}

	if result == nil {
		d := diag.NewErrorDiagnostic("Schema check failed", "The schema is not valid")
		return &d
	}

	if !result.Valid {
		d := diag.NewErrorDiagnostic("Schema check failed", fmt.Sprintf("The schema is not valid, see %s for more details", result.URL))
		return &d
	}

	data.Id = types.StringValue(result.Id)

	return nil
}
