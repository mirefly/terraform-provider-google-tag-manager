package gtm

import (
	"context"
	"terraform-provider-google-tag-manager/gtm/api"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"google.golang.org/api/tagmanager/v2"
)

var (
	_ resource.ResourceWithConfigure = &tagResource{}
)

func NewTagResource() resource.Resource {
	return &tagResource{}
}

type tagResource struct {
	client *api.ClientInWorkspace
}

// Configure adds the provider configured client to the resource.
func (r *tagResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*api.ClientInWorkspace)
}

// Metadata returns the resource type name.
func (r *tagResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tag"
}

// Schema defines the schema for the resource.
func (r *tagResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name":      schema.StringAttribute{Required: true},
			"type":      schema.StringAttribute{Required: true},
			"id":        schema.StringAttribute{Computed: true},
			"notes":     schema.StringAttribute{Optional: true},
			"parameter": parameterSchema,
			"firing_trigger_id": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

type resourceTagModel struct {
	Name            types.String             `tfsdk:"name"`
	Type            types.String             `tfsdk:"type"`
	Id              types.String             `tfsdk:"id"`
	Notes           types.String             `tfsdk:"notes"`
	Parameter       []ResourceParameterModel `tfsdk:"parameter"`
	FiringTriggerId []types.String           `tfsdk:"firing_trigger_id"`
}

func toResourceTag(tag *tagmanager.Tag) *resourceTagModel {
	return &resourceTagModel{
		Name:            types.StringValue(tag.Name),
		Type:            types.StringValue(tag.Type),
		Id:              types.StringValue(tag.TagId),
		Notes:           nullableStringValue(tag.Notes),
		Parameter:       toResourceParameter(tag.Parameter),
		FiringTriggerId: toResourceStringArray(tag.FiringTriggerId),
	}

}

func toApiTag(resource resourceTagModel) *tagmanager.Tag {
	return &tagmanager.Tag{
		Name:            resource.Name.ValueString(),
		Type:            resource.Type.ValueString(),
		TagId:           resource.Id.String(),
		Notes:           resource.Notes.ValueString(),
		Parameter:       toApiParameter(resource.Parameter),
		FiringTriggerId: unwrapStringArray(resource.FiringTriggerId),
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *tagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resourceTagModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tag, err := r.client.CreateTag(toApiTag(plan))
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Tag", err.Error())
		return
	}

	diags = resp.State.Set(ctx, toResourceTag(tag))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *tagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resourceTagModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tag, err := r.client.Tag(state.Id.ValueString())
	if err == api.ErrNotExist {
		resp.State.RemoveResource(ctx)
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error Reading Tag", err.Error())
		return
	}

	diags = resp.State.Set(ctx, toResourceTag(tag))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *tagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state resourceTagModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	tag, err := r.client.UpdateTag(state.Id.ValueString(), toApiTag(plan))
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Tag", err.Error())
		return
	}

	diags = resp.State.Set(ctx, toResourceTag(tag))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *tagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resourceTagModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteTag(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Deleting Tag", err.Error())
		return
	}
}
