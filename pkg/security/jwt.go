package security

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/dionofrizal88/dsi/auth/config"
	"time"
)

type JWT struct {
	Config config.Configuration

	JWTClaims JWTClaims
}

type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

// NewJWT is to define jwt object
func NewJWT(cfg config.Configuration) *JWT {
	return &JWT{
		Config: cfg,
	}
}

// GetJWTClaims is func for get JWT claims.
func (j *JWT) GetJWTClaims() JWTClaims {
	return j.JWTClaims
}

// GenerateJWT is func for generate JWT.
func (j *JWT) GenerateJWT(id, email string) (string, error) {
	claims := JWTClaims{
		UserID: id,
		Email:  email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.Config.AppSecret))
}

// ValidateJWT is func for validate JWT.
func (j *JWT) ValidateJWT(tokenString string) bool {
	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.Config.AppSecret), nil
	})

	j.JWTClaims = *claims

	return err == nil && token.Valid
}

// DecodeJWT is func for decode JWT.
func (j *JWT) DecodeJWT(tokenString string) (map[string]interface{}, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	data := make(map[string]interface{})
	for key, value := range claims {
		data[key] = value
	}

	return data, nil
}
