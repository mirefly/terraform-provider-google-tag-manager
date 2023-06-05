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
	_ resource.ResourceWithConfigure = &tagGroupResource{}
)

func NewTagGroupResource() resource.Resource {
	return &tagGroupResource{}
}

type tagGroupResource struct {
	client *api.ClientInWorkspace
}

// Configure adds the provider configured client to the resource.
func (r *tagGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*api.ClientInWorkspace)
}

// Metadata returns the resource type name.
func (r *tagGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tag_group"
}

// Schema defines the schema for the resource.
func (r *tagGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"elements": schema.MapNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: tagResourceSchemaAttributes,
				},
				Required: true,
			},
		},
	}
}

type resourceTagGroupModel struct {
	Elements map[string]resourceTagModel `tfsdk:"elements"`
}

func toResourceTagGroup(tags []*tagmanager.Tag) resourceTagGroupModel {
	var resourceTagGroup resourceTagGroupModel = resourceTagGroupModel{
		Elements: make(map[string]resourceTagModel, len(tags)),
	}

	for _, tag := range tags {
		resourceTagGroup.Elements[tag.Name] = toResourceTag(tag)
	}

	return resourceTagGroup
}

// Create creates the resource and sets the initial Terraform state.
func (r *tagGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resourceTagGroupModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createdElements := make(map[string]resourceTagModel)
	for _, element := range plan.Elements {
		tag, err := r.client.CreateTag(toApiTag(element))
		if err != nil {
			resp.Diagnostics.AddError("Error Creating Tag", err.Error())
			break
		}

		createdElements[tag.Name] = toResourceTag(tag)
	}

	plan.Elements = createdElements
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *tagGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resourceTagGroupModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var nameSet = make(map[string]struct{})
	for _, element := range state.Elements {
		nameSet[element.Name.ValueString()] = struct{}{}
	}

	tags, err := r.client.ListTags()
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Tag Group", err.Error())
		return
	}

	tagsInCurrentGroup := make([]*tagmanager.Tag, 0)
	for _, tag := range tags {
		if _, ok := nameSet[tag.Name]; ok {
			tagsInCurrentGroup = append(tagsInCurrentGroup, tag)
		}
	}

	diags = resp.State.Set(ctx, toResourceTagGroup(tagsInCurrentGroup))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *tagGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state resourceTagGroupModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete all the tags which doesn't exist in the new plan
	for _, element := range state.Elements {
		if _, ok := plan.Elements[element.Name.ValueString()]; !ok {
			tflog.Info(ctx, "Deleting Tag: "+element.Name.ValueString())

			err := r.client.DeleteTag(element.Id.ValueString())
			if err != nil {
				resp.Diagnostics.AddError("Error Deleting Tag", err.Error())
				break
			} else {
				delete(state.Elements, element.Name.ValueString())
			}
		}
	}

	// Create new tags which doesn't exist in the state
	for _, element := range plan.Elements {
		if _, ok := state.Elements[element.Name.ValueString()]; !ok {
			tflog.Info(ctx, "Creating Tag: "+element.Name.ValueString())

			tag, err := r.client.CreateTag(toApiTag(element))
			if err != nil {
				resp.Diagnostics.AddError("Error Creating Tag", err.Error())
				break
			} else {
				state.Elements[tag.Name] = toResourceTag(tag)
			}
		}
	}

	// Update tag if not the same in plan and state
	for _, stateEl := range state.Elements {
		if planEl, ok := plan.Elements[stateEl.Name.ValueString()]; ok {
			if !planEl.Equal(stateEl) {
				tflog.Info(ctx, "Updating Tag: "+stateEl.Name.ValueString())

				tag, err := r.client.UpdateTag(stateEl.Id.ValueString(), toApiTag(planEl))
				if err != nil {
					resp.Diagnostics.AddError("Error Updating Tag", err.Error())
					break
				} else {
					state.Elements[tag.Name] = toResourceTag(tag)
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
func (r *tagGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resourceTagGroupModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	elementsLeft := make(map[string]resourceTagModel)
	for name, element := range state.Elements {
		elementsLeft[name] = element
	}

	for name, element := range state.Elements {
		err := r.client.DeleteTag(element.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Error Deleting Tag Group", err.Error())
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
