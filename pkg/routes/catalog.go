package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/salberternst/portal/pkg/api"
	"github.com/salberternst/portal/pkg/middleware"
)

func getCatalog(ctx *gin.Context) {
	catalogRequest := api.CatalogRequest{}
	if err := ctx.BindJSON(&catalogRequest); err != nil {
		RespondWithBadRequest(ctx, "Bad Request")
		return
	}

	catalog, err := middleware.GetEdcAPI(ctx).GetCatalog(catalogRequest)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, catalog)
}

func getCatalogDataset(ctx *gin.Context) {
	var datasetRequest api.DatasetRequest
	if err := ctx.BindJSON(&datasetRequest); err != nil {
		RespondWithBadRequest(ctx, "Bad Request")
		return
	}

	catalogDataset, err := middleware.GetEdcAPI(ctx).GetCatalogDataset(datasetRequest)
	if err != nil {
		fmt.Println(err)
		RespondWithInternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, catalogDataset)
}

func addCatalogsRoutes(r *gin.RouterGroup) {
	contractAgreements := r.Group("/catalog")
	contractAgreements.POST("/", getCatalog)
	contractAgreements.POST("/dataset", getCatalogDataset)
}
