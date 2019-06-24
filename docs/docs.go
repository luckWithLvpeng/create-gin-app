// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2019-06-24 19:15:15.648883 +0800 CST m=+0.043052127

package docs

import (
	"bytes"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "电子物证",
        "contact": {},
        "license": {},
        "version": "1.0"
    },
    "host": "{{.Host}}",
    "basePath": "/v1",
    "paths": {
        "/role": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "tags": [
                    "role"
                ],
                "summary": "获取所有角色",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/controllers.Response"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "role"
                ],
                "summary": "添加角色",
                "parameters": [
                    {
                        "description": "role json",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.RoleForAdd"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/controllers.Response"
                        }
                    }
                }
            }
        },
        "/role/{id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "tags": [
                    "role"
                ],
                "summary": "获取单个角色",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "角色 id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/controllers.Response"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "role"
                ],
                "summary": "编辑角色",
                "parameters": [
                    {
                        "description": "role json",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.RoleForAdd"
                        }
                    },
                    {
                        "type": "integer",
                        "description": "role id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/controllers.Response"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "tags": [
                    "role"
                ],
                "summary": "删除角色",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "role id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/controllers.Response"
                        }
                    }
                }
            }
        },
        "/user": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "获取用户",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "now page",
                        "name": "nowPage",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "page size",
                        "name": "pageSize",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/controllers.Response"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "添加用户",
                "parameters": [
                    {
                        "description": "user for add",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.UserForAdd"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/controllers.Response"
                        }
                    }
                }
            }
        },
        "/user/login": {
            "post": {
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "tags": [
                    "user"
                ],
                "summary": "用户登录",
                "parameters": [
                    {
                        "type": "string",
                        "description": "user name",
                        "name": "Username",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "user password",
                        "name": "Password",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/controllers.Response"
                        }
                    }
                }
            }
        },
        "/user/logout": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "tags": [
                    "user"
                ],
                "summary": "用户退出, 把token加入黑名单",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/controllers.Response"
                        }
                    }
                }
            }
        },
        "/user/refreshToken": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "tags": [
                    "user"
                ],
                "summary": "用户获取新的token,新旧token会同时生效,旧的token 5分钟之后被销毁",
                "parameters": [
                    {
                        "type": "string",
                        "description": "refresh token",
                        "name": "RefreshToken",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/controllers.Response"
                        }
                    }
                }
            }
        },
        "/user/{id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "获取单个用户",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "user id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/controllers.Response"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "更新用户",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "user id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "user for update",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.UserForUpdate"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/controllers.Response"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "删除单个用户",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "user id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/controllers.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controllers.Response": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {
                    "type": "object"
                },
                "msg": {
                    "type": "string"
                }
            }
        },
        "models.RoleForAdd": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "roleName": {
                    "type": "string"
                }
            }
        },
        "models.UserForAdd": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "mobile": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "roleID": {
                    "type": "integer"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "models.UserForUpdate": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "mobile": {
                    "type": "string"
                },
                "roleID": {
                    "type": "integer"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo swaggerInfo

type s struct{}

func (s *s) ReadDoc() string {
	t, err := template.New("swagger_info").Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, SwaggerInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
