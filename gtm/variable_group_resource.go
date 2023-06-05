package gtm

import (
	"context"
	"terraform-provider-google-tag-manager/gtm/api"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"google.golang.org/api/tagmanager/v2"
)

var (
	_ resource.ResourceWithConfigure = &variableGroupResource{}
)

func NewVariableGroupResource() resource.Resource {
	return &variableGroupResource{}
}

type variableGroupResource struct {
	client *api.ClientInWorkspace
}

// Configure adds the provider configured client to the resource.
func (r *variableGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*api.ClientInWorkspace)
}

// Metadata returns the resource type name.
func (r *variableGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_variable_group"
}

// Schema defines the schema for the resource.
func (r *variableGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"elements": schema.MapNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: variableResourceSchemaAttributes,
				},
				Required: true,
			},
		},
	}
}

type resourceVariableGroupModel struct {
	Elements map[string]resourceVariableModel `tfsdk:"elements"`
}

func toResourceVariableGroup(variables []*tagmanager.Variable) resourceVariableGroupModel {
	var resourceVariableGroup resourceVariableGroupModel = resourceVariableGroupModel{
		Elements: make(map[string]resourceVariableModel, len(variables)),
	}

	for _, variable := range variables {
		resourceVariableGroup.Elements[variable.Name] = toResourceVariable(variable)
	}

	return resourceVariableGroup
}

// Create creates the resource and sets the initial Terraform state.
func (r *variableGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resourceVariableGroupModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createdElements := make(map[string]resourceVariableModel)
	for _, element := range plan.Elements {
		variable, err := r.client.CreateVariable(toApiVariable(element))
		if err != nil {
			resp.Diagnostics.AddError("Error Creating Variable", err.Error())
			break
		}

		createdElements[variable.Name] = toResourceVariable(variable)
	}

	plan.Elements = createdElements
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *variableGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resourceVariableGroupModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var nameSet = make(map[string]struct{})
	for _, element := range state.Elements {
		nameSet[element.Name.ValueString()] = struct{}{}
	}

	variables, err := r.client.ListVariables()
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Variable Group", err.Error())
		return
	}

	variablesInCurrentGroup := make([]*tagmanager.Variable, 0)
	for _, variable := range variables {
		if _, ok := nameSet[variable.Name]; ok {
			variablesInCurrentGroup = append(variablesInCurrentGroup, variable)
		}
	}

	diags = resp.State.Set(ctx, toResourceVariableGroup(variablesInCurrentGroup))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *variableGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state resourceVariableGroupModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete all the variables which doesn't exist in the new plan
	for _, element := range state.Elements {
		if _, ok := plan.Elements[element.Name.ValueString()]; !ok {
			err := r.client.DeleteVariable(element.Id.ValueString())
			if err != nil {
				resp.Diagnostics.AddError("Error Deleting Variable", err.Error())
				break
			} else {
				delete(state.Elements, element.Name.ValueString())
			}
		}
	}

	// Create new variables which doesn't exist in the state
	for _, element := range plan.Elements {
		if _, ok := state.Elements[element.Name.ValueString()]; !ok {
			variable, err := r.client.CreateVariable(toApiVariable(element))
			if err != nil {
				resp.Diagnostics.AddError("Error Creating Variable", err.Error())
				break
			} else {
				state.Elements[variable.Name] = toResourceVariable(variable)
			}
		}
	}

	// Update variable if not the same in plan and state
	for _, stateEl := range state.Elements {
		if planEl, ok := plan.Elements[stateEl.Name.ValueString()]; ok {
			if !planEl.Equal(stateEl) {
				variable, err := r.client.UpdateVariable(stateEl.Id.ValueString(), toApiVariable(planEl))
				if err != nil {
					resp.Diagnostics.AddError("Error Updating Variable", err.Error())
					break
				} else {
					state.Elements[variable.Name] = toResourceVariable(variable)
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
func (r *variableGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resourceVariableGroupModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	elementsLeft := make(map[string]resourceVariableModel)
	for name, element := range state.Elements {
		elementsLeft[name] = element
	}

	for name, element := range state.Elements {
		err := r.client.DeleteVariable(element.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Error Deleting Variable Group", err.Error())
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
