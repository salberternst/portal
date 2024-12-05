package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/salberternst/portal/pkg/middleware"
)

type UserInfo struct {
	FullName string   `json:"full_name"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
}

func getInfo(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, UserInfo{
		FullName: middleware.GetAuthenticatedUser(ctx).FullName,
		Email:    middleware.GetAuthenticatedUser(ctx).Email,
		Roles:    middleware.GetAuthenticatedUser(ctx).Roles,
	})
}

func addUserinfoRoutes(r *gin.RouterGroup) {
	r.GET("/userinfo", getInfo)
}
