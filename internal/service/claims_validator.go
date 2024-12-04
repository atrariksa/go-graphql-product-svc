package service

import (
	"github.com/golang-jwt/jwt/v4"
)

const AdminRole = "Admin"

type IClaimsValidator interface {
	IsAdmin(mapClaims jwt.MapClaims) bool
}

type ClaimsValidator struct {
}

func NewClaimsValidator() *ClaimsValidator {
	return &ClaimsValidator{}
}

func (cv *ClaimsValidator) IsAdmin(mapClaims jwt.MapClaims) bool {
	claims := mapClaims
	clientRole := claims["role"].(string)
	return clientRole == AdminRole
}
