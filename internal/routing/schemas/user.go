package schemas

import "github.com/getkin/kin-openapi/openapi3"

var UserSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewObjectSchema().Type,
		Properties: map[string]*openapi3.SchemaRef{
			"id": {
				Value: openapi3.NewUUIDSchema(),
			},
			"image": {
				Value: openapi3.NewStringSchema().
					WithFormat("uri"),
			},
			"name": {
				Value: openapi3.NewStringSchema().
					WithFormat("text").
					WithMin(3),
			},
			"username": {
				Value: openapi3.NewStringSchema().
					WithPattern(`(?:\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}\b)|(?:\b(?:\+?\d{1,3}[-.\s]?)?(?:\(?\d{2,4}\)?[-.\s]?)?\d{3,4}[-.\s]?\d{3,4}\b)`),
			},
			"bio": {
				Value: openapi3.NewStringSchema().
					WithFormat("text"),
			},
			"mfaEnabled": {
				Value: openapi3.NewBoolSchema(),
			},
			"mfaVerified": {
				Value: openapi3.NewBoolSchema(),
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
			"username",
			"mfaEnabled",
			"mfaVerified",
			"createdAt",
			"updatedAt",
		},
	},
}

var UsersSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewArraySchema().Type,
		Items: &openapi3.SchemaRef{
			Ref: "#/components/schemas/User",
		},
	},
}
