package middleware

import (
	"net/http"
	"os"

	"github.com/Clarilab/gocloaksession"
	"github.com/Nerzal/gocloak/v13"
	"github.com/gin-gonic/gin"
)

func KeycloakMiddleware() gin.HandlerFunc {
	keycloakClientId := os.Getenv("KEYCLOAK_CLIENT_ID")
	keycloakClientSecret := os.Getenv("KEYCLOAK_CLIENT_SECRET")
	keycloakRealm := os.Getenv("KEYCLOAK_REALM")
	keycloakUrl := os.Getenv("KEYCLOAK_URL")

	session, err := gocloaksession.NewSession(keycloakClientId, keycloakClientSecret, keycloakRealm, keycloakUrl)
	if err != nil {
		panic("Something wrong with the credentials or url")
	}

	return func(ctx *gin.Context) {
		client := session.GetGoCloakInstance()

		token, err := session.GetKeycloakAuthToken()
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized",
				"status":  http.StatusUnauthorized,
			})
			ctx.Abort()
			return
		}

		ctx.Set("keycloak-client", client)
		ctx.Set("keycloak-token", token.AccessToken)
		ctx.Set("keycloak-realm", keycloakRealm)
		ctx.Next()
	}
}

func GetKeycloakClient(ctx *gin.Context) *gocloak.GoCloak {
	return ctx.MustGet("keycloak-client").(*gocloak.GoCloak)
}

func GetKeycloakToken(ctx *gin.Context) string {
	return ctx.MustGet("keycloak-token").(string)
}

func GetKeycloakRealm(ctx *gin.Context) string {
	return ctx.MustGet("keycloak-realm").(string)
}
