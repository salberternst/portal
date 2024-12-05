package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthenticatedUser struct {
	Id         string   `validate:"required"`
	FullName   string   `validate:"required"`
	Email      string   `validate:"required,email"`
	TenantId   string   `validate:"required"`
	CustomerId string   `validate:"omitempty"`
	Roles      []string `validate:"required"`
}

func NewAuthenticatedUser(ctx *gin.Context) AuthenticatedUser {
	return AuthenticatedUser{
		Id:         ctx.GetHeader("X-User-Id"),
		FullName:   ctx.GetHeader("X-User-Full-Name"),
		Email:      ctx.GetHeader("X-User-Email"),
		TenantId:   ctx.GetHeader("X-User-Tenant-Id"),
		CustomerId: ctx.GetHeader("X-User-Customer-Id"),
		Roles:      strings.Split(ctx.GetHeader("X-User-Roles"), ","),
	}
}

func (u AuthenticatedUser) IsAdmin() bool {
	for _, role := range u.Roles {
		if role == "admin" {
			return true
		}
	}
	return false
}

func (u AuthenticatedUser) IsCustomer() bool {
	for _, role := range u.Roles {
		if role == "customer" {
			return true
		}
	}
	return false
}

func AuthMiddleware() gin.HandlerFunc {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return func(ctx *gin.Context) {
		user := NewAuthenticatedUser(ctx)
		if err := validate.Struct(user); err != nil {
			ctx.JSON(400, gin.H{
				"status":  403,
				"message": "Forbidden",
			})
			ctx.Abort()
			return
		}
		ctx.Set("authenticated-user", user)
		ctx.Next()
	}
}

func GetAuthenticatedUser(ctx *gin.Context) AuthenticatedUser {
	return ctx.MustGet("authenticated-user").(AuthenticatedUser)
}
