// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "AAStar Support",
            "url": "https://aastar.xyz"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/auth": {
            "post": {
                "description": "Get AccessToken By ApiKey",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "parameters": [
                    {
                        "description": "AccessToken Model",
                        "name": "credential",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.ClientCredential"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/api/healthz": {
            "get": {
                "description": "Get Healthz",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Healthz"
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/api/v1/get-support-entrypoint": {
            "get": {
                "security": [
                    {
                        "JWT": []
                    }
                ],
                "description": "get the support entrypoint",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Sponsor"
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/api/v1/get-support-strategy": {
            "get": {
                "security": [
                    {
                        "JWT": []
                    }
                ],
                "description": "get the support strategy",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Sponsor"
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/api/v1/try-pay-user-operation": {
            "post": {
                "security": [
                    {
                        "JWT": []
                    }
                ],
                "description": "sponsor the userOp",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Sponsor"
                ],
                "parameters": [
                    {
                        "description": "UserOp Request",
                        "name": "tryPay",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.TryPayUserOpRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        }
    },
    "definitions": {
        "model.ClientCredential": {
            "type": "object",
            "properties": {
                "apiKey": {
                    "type": "string"
                }
            }
        },
        "model.TryPayUserOpRequest": {
            "type": "object",
            "properties": {
                "extra": {},
                "force_entrypoint_address": {
                    "type": "string"
                },
                "force_network": {
                    "$ref": "#/definitions/types.NetWork"
                },
                "force_strategy_id": {
                    "type": "string"
                },
                "force_token": {
                    "type": "string"
                },
                "user_operation": {
                    "$ref": "#/definitions/model.UserOperationItem"
                }
            }
        },
        "model.UserOperationItem": {
            "type": "object",
            "required": [
                "call_data",
                "call_gas_limit",
                "max_fee_per_gas",
                "max_priority_fee_per_gas",
                "nonce",
                "per_verification_gas",
                "sender",
                "verification_gas_list"
            ],
            "properties": {
                "call_data": {
                    "type": "string"
                },
                "call_gas_limit": {
                    "type": "string"
                },
                "init_code": {
                    "type": "string"
                },
                "max_fee_per_gas": {
                    "type": "string"
                },
                "max_priority_fee_per_gas": {
                    "type": "string"
                },
                "nonce": {
                    "type": "string"
                },
                "per_verification_gas": {
                    "type": "string"
                },
                "sender": {
                    "type": "string"
                },
                "signature": {
                    "type": "string"
                },
                "verification_gas_list": {
                    "type": "string"
                }
            }
        },
        "types.NetWork": {
            "type": "string",
            "enum": [
                "ethereum",
                "sepolia",
                "arbitrum"
            ],
            "x-enum-varnames": [
                "Ethereum",
                "Sepolia",
                "Arbitrum"
            ]
        }
    },
    "securityDefinitions": {
        "JWT": {
            "description": "Type 'Bearer \\\u003cTOKEN\\\u003e' to correctly set the AccessToken",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
