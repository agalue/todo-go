// Package docs Code generated by swaggo/swag. DO NOT EDIT
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
            "name": "Alejandro Galue"
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
        "/api/v1/todos": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Get all the TODOs",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Todo"
                            }
                        }
                    },
                    "500": {
                        "description": "Backend error"
                    }
                }
            },
            "post": {
                "produces": [
                    "application/json"
                ],
                "summary": "Add a new TODO",
                "parameters": [
                    {
                        "description": "New TODO",
                        "name": "todo",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Base"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.Todo"
                        }
                    },
                    "500": {
                        "description": "Backend error"
                    }
                }
            }
        },
        "/api/v1/todos/{id}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Get a TODO",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "TODO ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Todo"
                        }
                    },
                    "404": {
                        "description": "Not found"
                    },
                    "500": {
                        "description": "Backend error"
                    }
                }
            },
            "put": {
                "produces": [
                    "application/json"
                ],
                "summary": "Update a TODO",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "TODO ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "TODO Status",
                        "name": "status",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Status"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Todo"
                        }
                    },
                    "404": {
                        "description": "Not found"
                    },
                    "500": {
                        "description": "Backend error"
                    }
                }
            },
            "delete": {
                "produces": [
                    "application/json"
                ],
                "summary": "Delete a TODO",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "TODO ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Todo"
                        }
                    },
                    "404": {
                        "description": "Not found"
                    },
                    "500": {
                        "description": "Backend error"
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Base": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "priority": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "models.Status": {
            "type": "object",
            "properties": {
                "completed": {
                    "type": "boolean"
                }
            }
        },
        "models.Todo": {
            "type": "object",
            "properties": {
                "completed": {
                    "type": "boolean"
                },
                "created_at": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "priority": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.0.1",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "TODO API",
	Description:      "A Simple TODO API based on PostgreSQL",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}