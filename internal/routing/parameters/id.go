package parameters

import "github.com/getkin/kin-openapi/openapi3"

var IdParameter = &openapi3.ParameterRef{
	Value: &openapi3.Parameter{
		In:              "path",
		Name:            "id",
		Description:     "The entity id.",
		AllowEmptyValue: false,
		Required:        true,
		Schema: &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type:   openapi3.NewStringSchema().Type,
				Format: "uuid",
			},
		},
	},
}
