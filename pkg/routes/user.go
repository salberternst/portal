package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/salberternst/portal/pkg/middleware"
)

type UserInfo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func getInfo(ctx *gin.Context) {
	claims := middleware.GetAccessTokenClaims(ctx)
	ctx.JSON(http.StatusOK, UserInfo{
		Email: claims.Email,
		Name:  claims.Name,
	})
}

func addUserRoutes(r *gin.RouterGroup) {
	userGroup := r.Group("/user")
	userGroup.GET("/info", getInfo)
}
