package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/jstermask/dynatrace_client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &extensionResource{}
	_ resource.ResourceWithConfigure = &extensionResource{}
)

type extensionResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name		types.String `tfsdk:"name"`
	LastUpdated types.String `tfsdk:"last_updated"`
	Payload     types.String `tfsdk:"payload"`
}

// NewExtensionResource is a helper function to simplify the provider implementation.
func NewExtensionResource() resource.Resource {
	return &extensionResource{}
}

// extensionResource is the resource implementation.
type extensionResource struct {
	client *dynatrace_client.DynatraceClient
}

// Metadata returns the resource type name.
func (r *extensionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_extension"
}

// Schema defines the schema for the resource.
func (r *extensionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"payload": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *extensionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan extensionResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

	createRequest := dynatrace_client.DynatraceExtensionCreateRequest {
		Payload: plan.Payload.ValueString(),
	}


	createResponse, err := r.client.CreateExtension(&createRequest)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create extension", fmt.Sprintf("Extension creation failed : %v", err))
		return
	}
	
	plan.Id = types.StringValue(createResponse.Id)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	plan.Name = types.StringValue(createResponse.Name)
	
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if(resp.Diagnostics.HasError()) {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *extensionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state extensionResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

	clientRequest := dynatrace_client.DynatraceExtensionGetBinaryRequest{
		Id: state.Id.ValueString(),
	}

	extensionPayload, err := r.client.GetExtensionBinary(&clientRequest)
	if err != nil {
		resp.Diagnostics.AddError("Unable to get latest extension", fmt.Sprintf("Extension %s get binary failed : %v", state.Id, err))
		return
	}

	state.Payload = types.StringValue(extensionPayload.Payload)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *extensionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *extensionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *extensionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*dynatrace_client.DynatraceClient)

	if !ok {
		resp.Diagnostics.AddError("Unexpected resource configuration type", fmt.Sprintf("Expected *dynatrace.DynatraceClient but got : %T. Please report this issue to provider developers.", req.ProviderData))
		return
	}

	r.client = client
}
