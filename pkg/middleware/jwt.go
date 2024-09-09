package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	TenantId    string `json:"tenant_id"`
	CustomerId  string `json:"customer_id"`
	RealmAccess struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
	jwt.RegisteredClaims
}

func TokenMiddleware() gin.HandlerFunc {
	parser := jwt.NewParser()

	return func(ctx *gin.Context) {
		idToken := ctx.GetHeader("X-Access-Token")
		token, _, err := parser.ParseUnverified(idToken, &Claims{})
		if err == nil {
			if claims, ok := token.Claims.(*Claims); ok {
				ctx.Set("access-token-claims", claims)
				ctx.Set("access-token", idToken)
				ctx.Next()
				return
			}
		}

		ctx.AbortWithError(http.StatusInternalServerError, errors.New("unable to decode X-Access-Token"))
	}
}

func GetAccessTokenClaims(ctx *gin.Context) *Claims {
	return ctx.MustGet("access-token-claims").(*Claims)
}

func GetAccessToken(ctx *gin.Context) string {
	return ctx.MustGet("access-token").(string)
}

func IsAdmin(ctx *gin.Context) bool {
	claims := GetAccessTokenClaims(ctx)
	for _, role := range claims.RealmAccess.Roles {
		if role == "admin" {
			return true
		}
	}
	return false
}

func IsCustomer(ctx *gin.Context) bool {
	claims := GetAccessTokenClaims(ctx)
	for _, role := range claims.RealmAccess.Roles {
		if role == "customer" {
			return true
		}
	}
	return false
}

func IsSysadmin(ctx *gin.Context) bool {
	claims := GetAccessTokenClaims(ctx)
	for _, role := range claims.RealmAccess.Roles {
		if role == "sysadmin" {
			return true
		}
	}
	return false
}
