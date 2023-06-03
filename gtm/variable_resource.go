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

var variableResourceSchemaAttributes = map[string]schema.Attribute{
	"name":      schema.StringAttribute{Required: true},
	"type":      schema.StringAttribute{Required: true},
	"id":        schema.StringAttribute{Computed: true},
	"notes":     schema.StringAttribute{Optional: true},
	"parameter": parameterSchema,
}

// Schema defines the schema for the resource.
func (r *variableResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: variableResourceSchemaAttributes,
	}
}

type resourceVariableModel struct {
	Name      types.String             `tfsdk:"name"`
	Type      types.String             `tfsdk:"type"`
	Id        types.String             `tfsdk:"id"`
	Notes     types.String             `tfsdk:"notes"`
	Parameter []ResourceParameterModel `tfsdk:"parameter"`
}

func (r resourceVariableModel) Equal(o resourceVariableModel) bool {
	if !r.Name.Equal(o.Name) ||
		!r.Type.Equal(o.Type) ||
		!r.Id.Equal(o.Id) ||
		!r.Notes.Equal(o.Notes) ||
		len(r.Parameter) != len(o.Parameter) {
		return false
	}

	for i := range r.Parameter {
		if !r.Parameter[i].Equal(o.Parameter[i]) {
			return false
		}
	}

	return true
}

func toResourceVariable(variable *tagmanager.Variable) resourceVariableModel {
	return resourceVariableModel{
		Name:      types.StringValue(variable.Name),
		Type:      types.StringValue(variable.Type),
		Id:        types.StringValue(variable.VariableId),
		Notes:     nullableStringValue(variable.Notes),
		Parameter: toResourceParameter(variable.Parameter),
	}
}
func toApiVariable(resource resourceVariableModel) *tagmanager.Variable {
	return &tagmanager.Variable{
		Name:       resource.Name.ValueString(),
		Type:       resource.Type.ValueString(),
		VariableId: resource.Id.String(),
		Notes:      resource.Notes.ValueString(),
		Parameter:  toApiParameter(resource.Parameter),
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *variableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resourceVariableModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	variable, err := r.client.CreateVariable(toApiVariable(plan))
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Variable", err.Error())
		return
	}

	diags = resp.State.Set(ctx, toResourceVariable(variable))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *variableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resourceVariableModel
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

	diags = resp.State.Set(ctx, toResourceVariable(variable))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *variableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state resourceVariableModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	variable, err := r.client.UpdateVariable(state.Id.ValueString(), toApiVariable(plan))
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Variable", err.Error())
		return
	}

	diags = resp.State.Set(ctx, toResourceVariable(variable))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *variableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resourceVariableModel

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
