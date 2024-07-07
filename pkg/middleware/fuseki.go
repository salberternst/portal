package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/salberternst/portal/pkg/api"
)

func FusekiMiddleware() gin.HandlerFunc {
	fusekiApi := api.NewFusekiAPI()

	return func(ctx *gin.Context) {
		ctx.Set("fuseki-api", fusekiApi)
		ctx.Next()
	}
}

func GetFusekiAPI(ctx *gin.Context) *api.FusekiAPI {
	return ctx.MustGet("fuseki-api").(*api.FusekiAPI)
}
