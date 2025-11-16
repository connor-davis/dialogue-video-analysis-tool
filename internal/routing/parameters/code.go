package parameters

import (
	"github.com/getkin/kin-openapi/openapi3"
)

var CodeParameter = &openapi3.ParameterRef{
	Value: &openapi3.Parameter{
		In:              "query",
		Name:            "code",
		Description:     "The authorization code received from the OAuth provider.",
		AllowEmptyValue: false,
		Required:        false,
		Schema: &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type: openapi3.NewStringSchema().Type,
			},
		},
	},
}
