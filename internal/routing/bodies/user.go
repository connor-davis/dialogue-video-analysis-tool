package bodies

import "github.com/getkin/kin-openapi/openapi3"

var CreateUserSchema = &openapi3.RequestBodyRef{
	Value: &openapi3.RequestBody{
		Content: openapi3.Content{
			"application/json": openapi3.NewMediaType().
				WithSchema(&openapi3.Schema{
					Type: openapi3.NewObjectSchema().Type,
					Properties: map[string]*openapi3.SchemaRef{
						"image": {
							Value: openapi3.NewStringSchema().
								WithFormat("uri"),
						},
						"name": {
							Value: openapi3.NewStringSchema().WithFormat("text").WithMin(3),
						},
						"username": {
							Value: openapi3.NewStringSchema().
								WithPattern(`(?:\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}\b)|(?:\b(?:\+?\d{1,3}[-.\s]?)?(?:\(?\d{2,4}\)?[-.\s]?)?\d{3,4}[-.\s]?\d{3,4}\b)`),
						},
						"bio": {
							Value: openapi3.NewStringSchema().
								WithFormat("text").WithNullable(),
						},
					},
					Required: []string{
						"name",
						"username",
					},
				}),
		},
		Description: "The payload to create a new user.",
		Required:    true,
	},
}

var UpdateUserSchema = &openapi3.RequestBodyRef{
	Value: &openapi3.RequestBody{
		Content: openapi3.Content{
			"application/json": openapi3.NewMediaType().
				WithSchema(&openapi3.Schema{
					Type: openapi3.NewObjectSchema().Type,
					Properties: map[string]*openapi3.SchemaRef{
						"image": {
							Value: openapi3.NewStringSchema().
								WithFormat("uri"),
						},
						"name": {
							Value: openapi3.NewStringSchema().WithFormat("text").WithMin(3),
						},
						"username": {
							Value: openapi3.NewStringSchema().
								WithPattern(`(?:\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}\b)|(?:\b(?:\+?\d{1,3}[-.\s]?)?(?:\(?\d{2,4}\)?[-.\s]?)?\d{3,4}[-.\s]?\d{3,4}\b)`),
						},
						"mfaEnabled": {
							Value: openapi3.NewBoolSchema(),
						},
						"mfaVerified": {
							Value: openapi3.NewBoolSchema(),
						},
						"bio": {
							Value: openapi3.NewStringSchema().
								WithFormat("text").WithNullable(),
						},
					},
					Required: []string{},
				}),
		},
		Description: "The payload to update an existing user.",
		Required:    true,
	},
}
