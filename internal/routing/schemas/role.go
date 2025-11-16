package schemas

import "github.com/getkin/kin-openapi/openapi3"

var RoleSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewObjectSchema().Type,
		Properties: map[string]*openapi3.SchemaRef{
			"id": {
				Value: openapi3.NewUUIDSchema(),
			},
			"name": {
				Value: openapi3.NewStringSchema().
					WithFormat("text").
					WithMin(3),
			},
			"description": {
				Value: openapi3.NewStringSchema().
					WithFormat("text").
					WithMin(3),
			},
			"permissions": {
				Value: &openapi3.Schema{
					Type: openapi3.NewArraySchema().Type,
					Items: &openapi3.SchemaRef{
						Value: openapi3.NewStringSchema().WithFormat("text"),
					},
				},
			},
			"createdAt": {
				Value: openapi3.NewDateTimeSchema(),
			},
			"updatedAt": {
				Value: openapi3.NewDateTimeSchema(),
			},
		},
		Required: []string{
			"id",
			"name",
			"description",
			"permissions",
			"createdAt",
			"updatedAt",
		},
	},
}

var RolesSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewArraySchema().Type,
		Items: &openapi3.SchemaRef{
			Ref: "#/components/schemas/Role",
		},
	},
}
