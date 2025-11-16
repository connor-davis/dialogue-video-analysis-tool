package parameters

import (
	"github.com/connor-davis/dialogue-video-analysis-tool/common"
	"github.com/getkin/kin-openapi/openapi3"
)

var ToParameter = &openapi3.ParameterRef{
	Value: &openapi3.Parameter{
		In:              "query",
		Name:            "to",
		Description:     "The url to redirect to after authentication.",
		AllowEmptyValue: false,
		Required:        false,
		Schema: &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type:    openapi3.NewStringSchema().Type,
				Default: common.EnvString("APP_BASE_URL", "http://localhost:3000"),
			},
		},
	},
}
