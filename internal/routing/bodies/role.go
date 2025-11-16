package bodies

import "github.com/getkin/kin-openapi/openapi3"

var CreateRoleSchema = &openapi3.RequestBodyRef{
	Value: &openapi3.RequestBody{
		Content: openapi3.Content{
			"application/json": openapi3.NewMediaType().
				WithSchema(&openapi3.Schema{
					Type: openapi3.NewObjectSchema().Type,
					Properties: map[string]*openapi3.SchemaRef{
						"name": {
							Value: openapi3.NewStringSchema().WithFormat("text").WithMin(3),
						},
						"description": {
							Value: openapi3.NewStringSchema().WithFormat("text").WithMin(3),
						},
					},
					Required: []string{
						"name",
						"description",
					},
				}),
		},
		Description: "The payload to create a new role.",
		Required:    true,
	},
}

var UpdateRoleSchema = &openapi3.RequestBodyRef{
	Value: &openapi3.RequestBody{
		Content: openapi3.Content{
			"application/json": openapi3.NewMediaType().
				WithSchema(&openapi3.Schema{
					Type: openapi3.NewObjectSchema().Type,
					Properties: map[string]*openapi3.SchemaRef{
						"name": {
							Value: openapi3.NewStringSchema().WithFormat("text").WithMin(3),
						},
						"description": {
							Value: openapi3.NewStringSchema().WithFormat("text").WithMin(3),
						},
						"permissions": {
							Value: &openapi3.Schema{
								Type: openapi3.NewArraySchema().Type,
								Items: &openapi3.SchemaRef{
									Value: openapi3.NewStringSchema().WithFormat("text"),
								},
							},
						},
					},
				}),
		},
		Description: "The payload to update an existing role.",
		Required:    true,
	},
}
