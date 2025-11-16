package parameters

import "github.com/getkin/kin-openapi/openapi3"

var DisablePaginationParameter = &openapi3.ParameterRef{
	Value: &openapi3.Parameter{
		In:              "query",
		Name:            "disablePagination",
		Description:     "Disable pagination.",
		AllowEmptyValue: false,
		Required:        false,
		Schema: &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type:    openapi3.NewBoolSchema().Type,
				Default: false,
			},
		},
	},
}
