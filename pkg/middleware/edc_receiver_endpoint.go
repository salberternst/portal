package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/salberternst/portal/pkg/api"
)

func EdcReceiverEndpointMiddleware() gin.HandlerFunc {
	edcApi := api.NewEdcReceiverEndpointApi()

	return func(ctx *gin.Context) {
		ctx.Set("edc-endpoint-receiver-api", edcApi)
		ctx.Next()
	}
}

func GetEdcReceiverEndpointAPI(ctx *gin.Context) *api.EdcReceiverEndpointApi {
	return ctx.MustGet("edc-endpoint-receiver-api").(*api.EdcReceiverEndpointApi)
}
