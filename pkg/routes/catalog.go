package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/salberternst/fsn_landing_page/pkg/api"
	"github.com/salberternst/fsn_landing_page/pkg/middleware"
)

func GetCatalog(ctx *gin.Context) {
	catalogRequest := api.CatalogRequest{}
	if err := ctx.BindJSON(&catalogRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"error":   "bad_request",
			"message": fmt.Sprintf("unable to bind catalog request: %v", err),
		})
		return
	}

	catalog, err := middleware.GetEdcAPI(ctx).GetCatalog(catalogRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, catalog)
}

func GetCatalogDataset(ctx *gin.Context) {
	var datasetRequest api.DatasetRequest
	if err := ctx.BindJSON(&datasetRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"error":   "bad_request",
			"message": fmt.Sprintf("unable to bind dataset request: %v", err),
		})
		return
	}

	catalogDataset, err := middleware.GetEdcAPI(ctx).GetCatalogDataset(datasetRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"error":   "unable_to_get_catalog_dataset",
			"message": fmt.Sprintf("unable to get catalog dataset: %v", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, catalogDataset)
}

func addCatalogsRoutes(r *gin.RouterGroup) {
	contractAgreements := r.Group("/catalog")
	contractAgreements.POST("/", GetCatalog)
	contractAgreements.POST("/dataset", GetCatalogDataset)
}
