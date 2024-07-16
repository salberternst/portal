package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/salberternst/portal/pkg/api"
	"github.com/salberternst/portal/pkg/middleware"
)

func getAssets(ctx *gin.Context) {
	querySpec, err := CreateQuerySpecFromContext(ctx)
	if err != nil {
		RespondWithBadRequest(ctx, "Bad Request")
		return
	}

	// only return assets that return to the customer
	assets, err := middleware.GetEdcAPI(ctx).GetAssets(querySpec)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, assets)
}

func getAsset(ctx *gin.Context) {
	id := ctx.Param("id")

	asset, err := middleware.GetEdcAPI(ctx).GetAsset(id)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	// check if current user has access to the asset
	if asset.PrivateProperties == nil || !CheckPrivateProperties(ctx, asset.PrivateProperties) {
		RespondWithResourceNotFound(ctx, id)
		return
	}

	ctx.JSON(http.StatusOK, asset)
}

func deleteAsset(ctx *gin.Context) {
	id := ctx.Param("id")

	asset, err := middleware.GetEdcAPI(ctx).GetAsset(id)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	// check if current user has access to the asset
	if asset.PrivateProperties == nil || !CheckPrivateProperties(ctx, asset.PrivateProperties) {
		// fixme: return 403 or 404
		RespondWithInternalServerError(ctx)
		return
	}

	if err := middleware.GetEdcAPI(ctx).DeleteAsset(id); err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func createAsset(ctx *gin.Context) {
	var asset api.Asset
	if err := ctx.BindJSON(&asset); err != nil {
		RespondWithBadRequest(ctx, "Bad Request")
		return
	}

	asset.Id = uuid.New().String()
	asset.PrivateProperties = BuildPrivatePropertiesFromContext(ctx)

	createdAsset, err := middleware.GetEdcAPI(ctx).CreateAsset(asset)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, createdAsset)
}

func addAssetsRoutes(r *gin.RouterGroup) {
	userGroup := r.Group("/assets")
	userGroup.GET("/", getAssets)
	userGroup.GET("/:id", getAsset)
	userGroup.DELETE("/:id", deleteAsset)
	userGroup.POST("/", createAsset)
}
