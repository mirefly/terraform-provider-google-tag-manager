package gtm

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"google.golang.org/api/tagmanager/v2"
)

var parameterSchema = genParameterSchema()
var conditionSchema = schema.ListNestedAttribute{
	Optional: true,
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"type":      schema.StringAttribute{Required: true},
			"parameter": parameterSchema,
		},
	},
}

func wrapParameterSchema(list schema.ListNestedAttribute) schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"key":   schema.StringAttribute{Optional: true},
				"type":  schema.StringAttribute{Required: true},
				"value": schema.StringAttribute{Optional: true},
				"list":  list,
				"map":   list,
			},
		},
	}
}

func genParameterSchema() schema.ListNestedAttribute {
	var s = schema.ListNestedAttribute{Optional: true, NestedObject: schema.NestedAttributeObject{}}

	for i := 0; i < 5; i++ {
		s = wrapParameterSchema(s)
	}

	return s
}

type ResourceParameterModel struct {
	Key   types.String              `tfsdk:"key"`
	Type  types.String              `tfsdk:"type"`
	Value types.String              `tfsdk:"value"`
	List  []*ResourceParameterModel `tfsdk:"list"`
	Map   []*ResourceParameterModel `tfsdk:"map"`
}

func unwrapParameter(resourceParameter []*ResourceParameterModel) []*tagmanager.Parameter {
	var parameter []*tagmanager.Parameter

	for _, p := range resourceParameter {
		var list, mmap []*tagmanager.Parameter

		if p.List != nil {
			list = unwrapParameter(p.List)
		}

		if p.Map != nil {
			mmap = unwrapParameter(p.Map)
		}

		parameter = append(parameter, &tagmanager.Parameter{
			Key:   p.Key.ValueString(),
			Type:  p.Type.ValueString(),
			Value: p.Value.ValueString(),
			List:  list,
			Map:   mmap,
		})
	}

	return parameter
}

func wrapParameter(parameter []*tagmanager.Parameter) []*ResourceParameterModel {
	var resourceParameter []*ResourceParameterModel = make([]*ResourceParameterModel, len(parameter))

	for i, p := range parameter {
		var list, mmap []*ResourceParameterModel

		if p.List != nil {
			list = wrapParameter(p.List)
		}

		if p.Map != nil {
			mmap = wrapParameter(p.Map)
		}

		resourceParameter[i] = &ResourceParameterModel{
			Key:   nullableStringValue(p.Key),
			Type:  nullableStringValue(p.Type),
			Value: nullableStringValue(p.Value),
			List:  list,
			Map:   mmap,
		}
	}

	return resourceParameter
}

func nullableStringValue(s string) types.String {
	if s != "" {
		return types.StringValue(s)
	} else {
		return types.StringNull()
	}
}

type ResourceConditionModel struct {
	Type      types.String              `tfsdk:"type"`
	Parameter []*ResourceParameterModel `tfsdk:"parameter"`
}

func unwrapCondition(resourceCondition []*ResourceConditionModel) []*tagmanager.Condition {
	condition := make([]*tagmanager.Condition, len(resourceCondition))

	for i, rc := range resourceCondition {
		var parameter []*tagmanager.Parameter
		if rc.Parameter != nil {
			parameter = unwrapParameter(rc.Parameter)
		}

		condition[i] = &tagmanager.Condition{
			Type:      rc.Type.ValueString(),
			Parameter: parameter,
		}
	}
	return condition
}

func wrapCondition(condition []*tagmanager.Condition) []*ResourceConditionModel {
	resourceCondition := make([]*ResourceConditionModel, len(condition))

	for i, c := range condition {
		var resourceParameter []*ResourceParameterModel
		if c.Parameter != nil {
			resourceParameter = wrapParameter(c.Parameter)
		}

		resourceCondition[i] = &ResourceConditionModel{
			Type:      nullableStringValue(c.Type),
			Parameter: resourceParameter,
		}
	}

	return resourceCondition
}

func wrapStringArray(list []string) []types.String {
	var rv []types.String

	for _, v := range list {
		rv = append(rv, types.StringValue(v))
	}

	return rv
}

func unwrapStringArray(list []types.String) []string {
	var rv []string

	for _, v := range list {
		rv = append(rv, v.ValueString())
	}

	return rv
}
