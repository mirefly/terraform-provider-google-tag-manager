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
	_ resource.ResourceWithConfigure = &variableResource{}
)

func NewVariableResource() resource.Resource {
	return &variableResource{}
}

type variableResource struct {
	client *api.ClientInWorkspace
}

// Configure adds the provider configured client to the resource.
func (r *variableResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*api.ClientInWorkspace)
}

// Metadata returns the resource type name.
func (r *variableResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_variable"
}

// Schema defines the schema for the resource.
func (r *variableResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name":      schema.StringAttribute{Required: true},
			"type":      schema.StringAttribute{Required: true},
			"id":        schema.StringAttribute{Computed: true},
			"notes":     schema.StringAttribute{Optional: true},
			"parameter": parameterSchema,
		},
	}
}

type variableResourceModel struct {
	Name      types.String              `tfsdk:"name"`
	Type      types.String              `tfsdk:"type"`
	Id        types.String              `tfsdk:"id"`
	Notes     types.String              `tfsdk:"notes"`
	Parameter []*ResourceParameterModel `tfsdk:"parameter"`
}

func overwriteVariableResource(variable *tagmanager.Variable, resource *variableResourceModel) {
	resource.Name = types.StringValue(variable.Name)
	resource.Type = types.StringValue(variable.Type)
	resource.Id = types.StringValue(variable.VariableId)

	resource.Parameter = wrapParameter(variable.Parameter)
}

func extractVariableParameter(resource variableResourceModel) []*tagmanager.Parameter {
	var parameter []*tagmanager.Parameter
	for _, p := range resource.Parameter {
		parameter = append(parameter, &tagmanager.Parameter{
			Key:   p.Key.ValueString(),
			Type:  p.Type.ValueString(),
			Value: p.Value.ValueString(),
		})
	}
	return parameter
}

func extractVariable(resource variableResourceModel) *tagmanager.Variable {
	parameter := extractVariableParameter(resource)
	return &tagmanager.Variable{
		Name:       resource.Name.ValueString(),
		Type:       resource.Type.ValueString(),
		VariableId: resource.Id.String(),
		Notes:      resource.Notes.ValueString(),
		Parameter:  parameter,
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *variableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan variableResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	variable, err := r.client.CreateVariable(extractVariable(plan))

	if err != nil {
		resp.Diagnostics.AddError("Error Creating Variable", err.Error())
		return
	}

	overwriteVariableResource(variable, &plan)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *variableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state variableResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	variable, err := r.client.Variable(state.Id.ValueString())
	if err == api.ErrNotExist {
		resp.State.RemoveResource(ctx)
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error Reading Variable", err.Error())
		return
	}

	overwriteVariableResource(variable, &state)
	diags = req.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *variableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state variableResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	variable, err := r.client.UpdateVariable(state.Id.ValueString(), extractVariable(plan))
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Variable", err.Error())
		return
	}

	overwriteVariableResource(variable, &plan)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *variableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state variableResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteVariable(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Deleting Variable", err.Error())
		return
	}
}
