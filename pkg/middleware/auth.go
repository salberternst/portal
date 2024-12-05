package middleware

import (
	"os"
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
	getHeaderOrEnv := func(header, envVar string) string {
		value := ctx.GetHeader(header)
		if value == "" {
			value = os.Getenv(envVar)
		}
		return value
	}

	return AuthenticatedUser{
		Id:         getHeaderOrEnv("X-User-Id", "DEFAULT_USER_ID"),
		FullName:   getHeaderOrEnv("X-User-Full-Name", "DEFAULT_USER_FULL_NAME"),
		Email:      getHeaderOrEnv("X-User-Email", "DEFAULT_USER_EMAIL"),
		TenantId:   getHeaderOrEnv("X-User-Tenant-Id", "DEFAULT_USER_TENANT_ID"),
		CustomerId: getHeaderOrEnv("X-User-Customer-Id", "DEFAULT_USER_CUSTOMER_ID"),
		Roles:      strings.Split(getHeaderOrEnv("X-User-Roles", "DEFAULT_USER_ROLES"), ","),
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
