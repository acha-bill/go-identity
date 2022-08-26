// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/identity": {
            "post": {
                "description": "Signs a document bundle for the current user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "GetIdentity",
                "parameters": [
                    {
                        "description": "req",
                        "name": "doc",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.GetIdentityReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.GetIdentityRes"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "domain.Identity": {
            "type": "object",
            "properties": {
                "dob": {
                    "type": "string",
                    "example": ""
                },
                "firstName": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "lastName": {
                    "type": "string"
                }
            }
        },
        "handler.GetIdentityReq": {
            "type": "object",
            "properties": {
                "documentHash": {
                    "type": "string"
                }
            }
        },
        "handler.GetIdentityRes": {
            "type": "object",
            "properties": {
                "documentHash": {
                    "type": "string"
                },
                "identity": {
                    "$ref": "#/definitions/domain.Identity"
                },
                "signature": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "go-identity demo",
	Description:      "identity demo",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
