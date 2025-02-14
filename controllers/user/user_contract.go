package user

import (
	"context"
	"github.com/dionofrizal88/dsi/auth/domain/entity"
)

// Request struct is used to get request value.
type Request struct {
	Name     string `form:"name" json:"name"`
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
}

// RecoveryRequest struct is used to get recovery request value.
type RecoveryRequest struct {
	Token       string `form:"token" json:"token" param:"token"`
	Email       string `form:"email" json:"email"`
	NewPassword string `form:"new_password" json:"new_password"`
}

// Response struct is used to get response value.
type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    interface{} `json:"meta,omitempty"`
}

// ResponseAuth struct is used to get response value.
type ResponseAuth struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	AccessToken string `json:"access_token"`
}

// transformToResponse is a function to transform user into response value.
func (co *Controller) transformToResponse(ctx context.Context, message string, user *entity.User) Response {
	var response Response
	response.Message = message
	response.Data = user

	return response
}

// transformToAuthResponse is a function to transform user into response value.
func (co *Controller) transformToAuthResponse(ctx context.Context, message, token string, user *entity.User) Response {
	var response Response

	authResponse := ResponseAuth{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		AccessToken: token,
	}

	response.Message = message
	response.Data = authResponse

	return response
}
