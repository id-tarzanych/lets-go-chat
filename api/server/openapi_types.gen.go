// Package server provides primitives to interact with the openapi HTTP API.
//
// Code generated by unknown module path version unknown version DO NOT EDIT.
package server

// ActiveUsersResponse defines model for ActiveUsersResponse.
type ActiveUsersResponse struct {
	Count int `json:"count"`
}

// CreateUserRequest defines model for CreateUserRequest.
type CreateUserRequest struct {
	Password string `json:"password"`
	UserName string `json:"userName"`
}

// CreateUserResponse defines model for CreateUserResponse.
type CreateUserResponse struct {
	Id       *string `json:"id,omitempty"`
	UserName *string `json:"userName,omitempty"`
}

// LoginUserRequest defines model for LoginUserRequest.
type LoginUserRequest struct {
	// The password for login in clear text
	Password string `json:"password"`

	// The user name for login
	UserName string `json:"userName"`
}

// LoginUserResponse defines model for LoginUserResponse.
type LoginUserResponse struct {
	// A url for websoket API with a one-time token for starting chat
	Url string `json:"url"`
}

// WsRTMStartParams defines parameters for WsRTMStart.
type WsRTMStartParams struct {
	// One time token for a loged user
	Token string `json:"token"`
}

// CreateUserJSONBody defines parameters for CreateUser.
type CreateUserJSONBody CreateUserRequest

// LoginUserJSONBody defines parameters for LoginUser.
type LoginUserJSONBody LoginUserRequest

// CreateUserJSONRequestBody defines body for CreateUser for application/json ContentType.
type CreateUserJSONRequestBody CreateUserJSONBody

// LoginUserJSONRequestBody defines body for LoginUser for application/json ContentType.
type LoginUserJSONRequestBody LoginUserJSONBody