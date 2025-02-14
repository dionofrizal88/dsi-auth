package admin

import (
	"context"
	"github.com/dionofrizal88/dsi/auth/domain/entity"
)

var (
	RecoveryRequestEmailTemplate = "Hello %s,\n\nWe received a request to reset your password for your Digital Sekuriti Indonesia account. If you made this request, please click the link below to reset your password:\n\n🔗 %s\n\nThis link is valid for %s. If you did not request a password reset, please ignore this email or contact our support team.\n\nFor security reasons, do not share this link with anyone.\n\nBest regards,\nDigital Sekuriti Indonesia\nadmin@digitalsekuriti.id | https://digitalsekuriti.id/"
)

// Request struct is used to get request value.
type Request struct {
	Name     string `form:"name" json:"name"`
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
}

// Response struct is used to get response value.
type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    interface{} `json:"meta,omitempty"`
}

// transformToResponse is a function to transform user into response value.
func (co *Controller) transformToResponse(ctx context.Context, message string, user *entity.User) Response {
	var response Response
	response.Message = message
	response.Data = user

	return response
}
