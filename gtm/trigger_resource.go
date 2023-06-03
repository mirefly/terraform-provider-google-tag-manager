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
	_ resource.ResourceWithConfigure = &triggerResource{}
)

func NewTriggerResource() resource.Resource {
	return &triggerResource{}
}

type triggerResource struct {
	client *api.ClientInWorkspace
}

// Configure adds the provider configured client to the resource.
func (r *triggerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*api.ClientInWorkspace)
}

// Metadata returns the resource type name.
func (r *triggerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_trigger"
}

// Schema defines the schema for the resource.
func (r *triggerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name":                schema.StringAttribute{Required: true},
			"type":                schema.StringAttribute{Required: true},
			"id":                  schema.StringAttribute{Computed: true},
			"notes":               schema.StringAttribute{Optional: true},
			"custom_event_filter": conditionSchema,
		},
	}
}

type resourceTriggerModel struct {
	Name              types.String             `tfsdk:"name"`
	Type              types.String             `tfsdk:"type"`
	Id                types.String             `tfsdk:"id"`
	Notes             types.String             `tfsdk:"notes"`
	CustomEventFilter []resourceConditionModel `tfsdk:"custom_event_filter"`
}

func toResourceTrigger(trigger *tagmanager.Trigger) *resourceTriggerModel {
	return &resourceTriggerModel{
		Name:              types.StringValue(trigger.Name),
		Type:              types.StringValue(trigger.Type),
		Id:                types.StringValue(trigger.TriggerId),
		Notes:             nullableStringValue(trigger.Notes),
		CustomEventFilter: toResourceCondition(trigger.CustomEventFilter),
	}
}

func toApiTrigger(resource resourceTriggerModel) *tagmanager.Trigger {
	return &tagmanager.Trigger{
		Name:              resource.Name.ValueString(),
		Type:              resource.Type.ValueString(),
		TriggerId:         resource.Id.ValueString(),
		Notes:             resource.Notes.ValueString(),
		CustomEventFilter: toApiCondition(resource.CustomEventFilter),
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *triggerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resourceTriggerModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	trigger, err := r.client.CreateTrigger(toApiTrigger(plan))
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Trigger", err.Error())
		return
	}

	diags = resp.State.Set(ctx, toResourceTrigger(trigger))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *triggerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resourceTriggerModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	trigger, err := r.client.Trigger(state.Id.ValueString())
	if err == api.ErrNotExist {
		resp.State.RemoveResource(ctx)
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error Reading Trigger", err.Error())
		return
	}

	diags = resp.State.Set(ctx, toResourceTrigger(trigger))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *triggerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state resourceTriggerModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	trigger, err := r.client.UpdateTrigger(state.Id.ValueString(), toApiTrigger(plan))
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Trigger", err.Error())
		return
	}

	diags = resp.State.Set(ctx, toResourceTrigger(trigger))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *triggerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resourceTriggerModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteTrigger(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Deleting Trigger", err.Error())
		return
	}
}
