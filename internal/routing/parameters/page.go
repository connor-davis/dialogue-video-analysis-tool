package parameters

import "github.com/getkin/kin-openapi/openapi3"

var PageParameter = &openapi3.ParameterRef{
	Value: &openapi3.Parameter{
		In:              "query",
		Name:            "page",
		Description:     "The page number.",
		AllowEmptyValue: false,
		Required:        false,
		Schema: &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type:    openapi3.NewIntegerSchema().Type,
				Format:  "int64",
				Default: 1,
			},
		},
	},
}
