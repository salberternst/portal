package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/salberternst/portal/pkg/api"
)

func FederatedCatalogMiddleware() gin.HandlerFunc {
	federatedCatalog := api.NewFederatedCatalogAPI()

	return func(ctx *gin.Context) {
		ctx.Set("federated-catalog-api", federatedCatalog)
		ctx.Next()
	}
}

func GetFederatedCatalogAPI(ctx *gin.Context) *api.FederatedCatalogAPI {
	return ctx.MustGet("federated-catalog-api").(*api.FederatedCatalogAPI)
}
