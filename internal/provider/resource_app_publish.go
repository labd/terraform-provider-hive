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
var _ resource.Resource = &HiveAppPublishResource{}
var _ resource.ResourceWithImportState = &HiveAppPublishResource{}

// NewHiveAppPublishResource is a helper function to simplify the provider implementation.
func NewHiveAppPublishResource() resource.Resource {
	return &HiveAppPublishResource{}
}

// HiveAppPublishResource defines the resource implementation.
type HiveAppPublishResource struct {
	client *sdk.HiveClient
}

// HiveAppPublishResourceModel describes the resource data model.
type HiveAppPublishResourceModel struct {
	Id      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Version types.String `tfsdk:"version"`
}

// Metadata returns the resource type name.
func (r *HiveAppPublishResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_publish"
}

// Schema defines the schema for the hive_schema_check resource.
func (r *HiveAppPublishResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource to publush an app within Hive",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The service name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "The commit or version identifier",
				Required:            true,
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
func (r *HiveAppPublishResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *HiveAppPublishResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data HiveAppPublishResourceModel

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
func (r *HiveAppPublishResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data HiveAppPublishResourceModel

	// Retrieve state.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save any updates back to state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update handles updates to the resource.
func (r *HiveAppPublishResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data HiveAppPublishResourceModel

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
func (r *HiveAppPublishResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

// ImportState allows the resource to be imported into Terraform.
func (r *HiveAppPublishResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *HiveAppPublishResource) ExecuteRequest(ctx context.Context, data *HiveAppPublishResourceModel) *diag.ErrorDiagnostic {
	result, err := r.client.PublishApp(ctx, &sdk.PublishAppInput{
		Name:    data.Name.ValueString(),
		Version: data.Version.ValueString(),
	})

	if err != nil {
		d := diag.NewErrorDiagnostic("App creation failed", err.Error())
		return &d
	}

	data.Id = types.StringValue(result.Id)
	data.Name = types.StringValue(result.AppName)
	data.Version = types.StringValue(result.AppVersion)

	return nil
}
