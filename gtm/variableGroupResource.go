package gtm

import (
	"context"
	"terraform-provider-google-tag-manager/gtm/api"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
			},
		},
	}
}

type variableGroupResourceModel struct {
	Elements map[string]*variableResourceModel `tfsdk:"elements"`
}

// Create creates the resource and sets the initial Terraform state.
func (r *variableGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan variableGroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for i, element := range plan.Elements {
		variable, err := r.client.CreateVariable(unwrapVariable(*element))
		if err != nil {
			resp.Diagnostics.AddError("Error Creating Variable", err.Error())
			return
		}

		plan.Elements[i] = wrapVariable(variable)
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *variableGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state variableGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	variables, err := r.client.ListVariables()
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Variable", err.Error())
		return
	}

	var newState variableGroupResourceModel = variableGroupResourceModel{
		Elements: make(map[string]*variableResourceModel),
	}
	for _, variable := range variables {
		if _, ok := state.Elements[variable.Name]; !ok {
		}

		newState.Elements[variable.Name] = wrapVariable(variable)
	}
	diags = req.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *variableGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//	var plan, state variableGroupResourceModel
	//
	//	diags := req.Plan.Get(ctx, &plan)
	//	resp.Diagnostics.Append(diags...)
	//
	//	diags = req.State.Get(ctx, &state)
	//	resp.Diagnostics.Append(diags...)
	//
	//	if resp.Diagnostics.HasError() {
	//		return
	//	}
	//
	//	variable, err := r.client.UpdateVariable(state.Id.ValueString(), unwrapVariable(plan))
	//	if err != nil {
	//		resp.Diagnostics.AddError("Error Updating Variable", err.Error())
	//		return
	//	}
	//
	//	wrapVariable(variable, &plan)
	//	diags = resp.State.Set(ctx, plan)
	//	resp.Diagnostics.Append(diags...)
	//	if resp.Diagnostics.HasError() {
	//		return
	//	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *variableGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	//	var state variableGroupResourceModel
	//
	//	diags := req.State.Get(ctx, &state)
	//	resp.Diagnostics.Append(diags...)
	//	if resp.Diagnostics.HasError() {
	//		return
	//	}
	//
	//	err := r.client.DeleteVariable(state.Id.ValueString())
	//	if err != nil {
	//		resp.Diagnostics.AddError("Error Deleting Variable", err.Error())
	//		return
	//	}
}
