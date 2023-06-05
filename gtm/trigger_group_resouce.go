package gtm

import (
	"context"
	"terraform-provider-google-tag-manager/gtm/api"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"google.golang.org/api/tagmanager/v2"
)

var (
	_ resource.ResourceWithConfigure = &triggerGroupResource{}
)

func NewTriggerGroupResource() resource.Resource {
	return &triggerGroupResource{}
}

type triggerGroupResource struct {
	client *api.ClientInWorkspace
}

// Configure adds the provider configured client to the resource.
func (r *triggerGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*api.ClientInWorkspace)
}

// Metadata returns the resource type name.
func (r *triggerGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_trigger_group"
}

// Schema defines the schema for the resource.
func (r *triggerGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"elements": schema.MapNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: triggerResourceSchemaAttributes,
				},
				Required: true,
			},
		},
	}
}

type resourceTriggerGroupModel struct {
	Elements map[string]resourceTriggerModel `tfsdk:"elements"`
}

func toResourceTriggerGroup(triggers []*tagmanager.Trigger) resourceTriggerGroupModel {
	var resourceTriggerGroup resourceTriggerGroupModel = resourceTriggerGroupModel{
		Elements: make(map[string]resourceTriggerModel, len(triggers)),
	}

	for _, trigger := range triggers {
		resourceTriggerGroup.Elements[trigger.Name] = toResourceTrigger(trigger)
	}

	return resourceTriggerGroup
}

// Create creates the resource and sets the initial Terraform state.
func (r *triggerGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resourceTriggerGroupModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createdElements := make(map[string]resourceTriggerModel)
	for _, element := range plan.Elements {
		trigger, err := r.client.CreateTrigger(toApiTrigger(element))
		if err != nil {
			resp.Diagnostics.AddError("Error Creating Trigger", err.Error())
			break
		}

		createdElements[trigger.Name] = toResourceTrigger(trigger)
	}

	plan.Elements = createdElements
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *triggerGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resourceTriggerGroupModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var nameSet = make(map[string]struct{})
	for _, element := range state.Elements {
		nameSet[element.Name.ValueString()] = struct{}{}
	}

	triggers, err := r.client.ListTriggers()
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Trigger Group", err.Error())
		return
	}

	triggersInCurrentGroup := make([]*tagmanager.Trigger, 0)
	for _, trigger := range triggers {
		if _, ok := nameSet[trigger.Name]; ok {
			triggersInCurrentGroup = append(triggersInCurrentGroup, trigger)
		}
	}

	diags = resp.State.Set(ctx, toResourceTriggerGroup(triggersInCurrentGroup))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *triggerGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state resourceTriggerGroupModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete all the triggers which doesn't exist in the new plan
	for _, element := range state.Elements {
		if _, ok := plan.Elements[element.Name.ValueString()]; !ok {
			tflog.Info(ctx, "Deleting Trigger: "+element.Name.ValueString())

			err := r.client.DeleteTrigger(element.Id.ValueString())
			if err != nil {
				resp.Diagnostics.AddError("Error Deleting Trigger", err.Error())
				break
			} else {
				delete(state.Elements, element.Name.ValueString())
			}
		}
	}

	// Create new triggers which doesn't exist in the state
	for _, element := range plan.Elements {
		if _, ok := state.Elements[element.Name.ValueString()]; !ok {
			tflog.Info(ctx, "Creating Trigger: "+element.Name.ValueString())

			trigger, err := r.client.CreateTrigger(toApiTrigger(element))
			if err != nil {
				resp.Diagnostics.AddError("Error Creating Trigger", err.Error())
				break
			} else {
				state.Elements[trigger.Name] = toResourceTrigger(trigger)
			}
		}
	}

	// Update trigger if not the same in plan and state
	for _, stateEl := range state.Elements {
		if planEl, ok := plan.Elements[stateEl.Name.ValueString()]; ok {
			tflog.Info(ctx, "Updating Trigger: "+stateEl.Name.ValueString())

			if !planEl.Equal(stateEl) {
				trigger, err := r.client.UpdateTrigger(stateEl.Id.ValueString(), toApiTrigger(planEl))
				if err != nil {
					resp.Diagnostics.AddError("Error Updating Trigger", err.Error())
					break
				} else {
					state.Elements[trigger.Name] = toResourceTrigger(trigger)
				}
			}
		}
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *triggerGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resourceTriggerGroupModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	elementsLeft := make(map[string]resourceTriggerModel)
	for name, element := range state.Elements {
		elementsLeft[name] = element
	}

	for name, element := range state.Elements {
		err := r.client.DeleteTrigger(element.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Error Deleting Trigger Group", err.Error())
			break
		}
		delete(elementsLeft, name)
	}

	state.Elements = elementsLeft

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
