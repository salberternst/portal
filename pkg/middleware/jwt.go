package middleware

import (
	"errors"
	"net/http"
	"strings"

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

func JwtMiddleware() gin.HandlerFunc {
	parser := jwt.NewParser()

	return func(ctx *gin.Context) {
		idToken := ctx.GetHeader("X-Access-Token")
		if idToken != "" {
			token, _, err := parser.ParseUnverified(idToken, &Claims{})
			if err == nil {
				if claims, ok := token.Claims.(*Claims); ok {
					// extract claims from token and set them as headers
					ctx.Request.Header.Set("X-User-Id", claims.Subject)
					ctx.Request.Header.Set("X-User-Email", claims.Email)
					ctx.Request.Header.Set("X-User-Full-Name", claims.Name)
					ctx.Request.Header.Set("X-User-Roles", strings.Join(claims.RealmAccess.Roles, ","))
					ctx.Request.Header.Set("X-User-Tenant-Id", claims.TenantId)
					ctx.Request.Header.Set("X-User-Customer-Id", claims.CustomerId)
					ctx.Set("access-token-claims", claims)
					ctx.Set("access-token", idToken)
					ctx.Next()
					return
				}
			}
			ctx.AbortWithError(http.StatusInternalServerError, errors.New("unable to decode X-Access-Token"))
		}
	}
}

func GetAccessTokenClaims(ctx *gin.Context) *Claims {
	return ctx.MustGet("access-token-claims").(*Claims)
}

func GetAccessToken(ctx *gin.Context) string {
	return ctx.MustGet("access-token").(string)
}
