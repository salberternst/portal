package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/salberternst/fsn_landing_page/pkg/api"
	"github.com/salberternst/fsn_landing_page/pkg/middleware"
	"github.com/salberternst/fsn_landing_page/pkg/utils"
)

type Asset struct {
	ID         any    `json:"@id"`
	Type       string `json:"@type"`
	Properties struct {
		Name        string `json:"name"`
		ContentType string `json:"contenttype"`
	} `json:"properties"`
	DataAddress map[string]string `json:"dataAddress"`
}

type AssetQuery struct {
	Page     uint `form:"page" binding:"required"`
	PageSize uint `form:"page_size"  binding:"required"`
}

func getAssets(ctx *gin.Context) {
	assetQuery := AssetQuery{}
	if err := ctx.BindQuery(&assetQuery); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"error":   "bad_request",
			"message": fmt.Sprintf("unable to bind asset query: %v", err),
		})
		return
	}

	querySpec := api.QuerySpec{
		Context: map[string]string{
			"@vocab": "https://w3id.org/edc/v0.0.1/ns/",
		},
		Type:             "QuerySpec",
		Offset:           (assetQuery.Page - 1) * assetQuery.PageSize,
		Limit:            assetQuery.PageSize,
		FilterExpression: utils.BuildFilterExpressionFromContext(ctx),
	}

	assets, err := middleware.GetEdcAPI(ctx).GetAssets(querySpec)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"error":   "unable_to_get_assets",
			"message": fmt.Sprintf("unable to get assets: %v", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, assets)
}

func getAsset(ctx *gin.Context) {
	id := ctx.Param("id")

	asset, err := middleware.GetEdcAPI(ctx).GetAsset(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"error":   "unable_to_get_asset",
			"message": fmt.Sprintf("unable to get asset: %v", err),
		})
		return
	}

	if asset.PrivateProperties == nil || !utils.CheckPrivateProperties(ctx, asset.PrivateProperties) {
		ctx.JSON(http.StatusForbidden, gin.H{
			"status":  http.StatusForbidden,
			"error":   "forbidden",
			"message": "You are not allowed to view this asset",
		})
		return
	}

	ctx.JSON(http.StatusOK, asset)
}

func deleteAsset(ctx *gin.Context) {
	id := ctx.Param("id")

	asset, err := middleware.GetEdcAPI(ctx).GetAsset(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"error":   "unable_to_get_asset",
			"message": fmt.Sprintf("unable to get asset: %v", err),
		})
		return
	}

	if asset.PrivateProperties == nil || !utils.CheckPrivateProperties(ctx, asset.PrivateProperties) {
		ctx.JSON(http.StatusForbidden, gin.H{
			"status":  http.StatusForbidden,
			"error":   "forbidden",
			"message": "You are not allowed to delete this asset",
		})
		return
	}

	err = middleware.GetEdcAPI(ctx).DeleteAsset(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"error":   "unable_to_delete_asset",
			"message": fmt.Sprintf("unable to delete asset: %v", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func createAsset(ctx *gin.Context) {
	var asset api.Asset
	if err := ctx.BindJSON(&asset); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"error":   "bad_request",
			"message": fmt.Sprintf("unable to bind asset: %v", err),
		})
		return
	}

	asset.Id = uuid.New().String()
	asset.PrivateProperties = utils.BuildPrivatePropertiesFromContext(ctx)

	createdAsset, err := middleware.GetEdcAPI(ctx).CreateAsset(asset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"error":   "unable_to_create_asset",
			"message": fmt.Sprintf("unable to create asset: %v", err),
		})
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
