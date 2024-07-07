package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/salberternst/portal/pkg/api"
)

func ThingsboardMiddleware() gin.HandlerFunc {
	thingsboardApi := api.NewThingsboardAPI()

	return func(ctx *gin.Context) {
		ctx.Set("thingsboard-api", thingsboardApi)
		ctx.Next()
	}
}

func GetThingsboardAPI(ctx *gin.Context) *api.ThingsboardAPI {
	return ctx.MustGet("thingsboard-api").(*api.ThingsboardAPI)
}
