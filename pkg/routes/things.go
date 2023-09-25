package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/salberternst/fsn_landing_page/pkg/models"
)

func getThings(ctx *gin.Context) {
	accessToken := ctx.GetHeader("X-Access-Token")

	credentials, err := models.GetThingsboardToken(accessToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	things, err := models.GetThingsboardThings(credentials.Token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	ctx.JSON(200, things)
}

func getThingCredentials(ctx *gin.Context) {
	credentials, err := models.GetThingsboardToken(ctx.GetHeader("X-Access-Token"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	thingsboardCredentials, err := models.GetThingsboardCredentials(ctx.Param("id"), credentials.Token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, thingsboardCredentials)
}

func createThing(ctx *gin.Context) {
	var thing models.ThingsboardCreateThing
	if err := ctx.ShouldBindJSON(&thing); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	credentials, err := models.GetThingsboardToken(ctx.GetHeader("X-Access-Token"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	err = models.CreateThingsboardThing(thing, credentials.Token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

}

func deleteThing(ctx *gin.Context) {
	credentials, err := models.GetThingsboardToken(ctx.GetHeader("X-Access-Token"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = models.DeleteThingsboardThing(ctx.Param("id"), credentials.Token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

func getThing(ctx *gin.Context) {
	credentials, err := models.GetThingsboardToken(ctx.GetHeader("X-Access-Token"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	thingsboardThing, err := models.GetThingsboardThing(ctx.Param("id"), credentials.Token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, thingsboardThing)
}

func addThingRoutes(r *gin.RouterGroup) {
	thingsGroup := r.Group("/things")
	thingsGroup.GET("/", getThings)
	thingsGroup.POST("/", createThing)
	thingsGroup.GET("/:id", getThing)
	thingsGroup.GET("/:id/credentials", getThingCredentials)
	thingsGroup.DELETE("/:id", deleteThing)
}
