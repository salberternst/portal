package routes

import (
	"github.com/gin-gonic/gin"
)

func getHealth(ctx *gin.Context) {
	
}

func addHealthRoutes(r *gin.RouterGroup) {
	r.GET("/health", getHealth)
}
