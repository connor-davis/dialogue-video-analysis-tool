package parameters

import "github.com/getkin/kin-openapi/openapi3"

var PreloadParameter = &openapi3.ParameterRef{
	Value: &openapi3.Parameter{
		In:              "query",
		Name:            "preload",
		Description:     "The related entities to preload.",
		AllowEmptyValue: true,
		Required:        false,
		Schema: &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type: openapi3.NewArraySchema().Type,
				Items: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: openapi3.NewStringSchema().Type,
					},
				},
			},
		},
	},
}
