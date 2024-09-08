package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/salberternst/portal/pkg/api"
	"github.com/salberternst/portal/pkg/middleware"
)

func getFederatedCatalog(ctx *gin.Context) {
	querySpec := api.QuerySpec{}
	if err := ctx.BindJSON(&querySpec); err != nil {
		RespondWithBadRequest(ctx, "Bad Request")
		return
	}

	catalog, err := middleware.GetFederatedCatalogAPI(ctx).Query(querySpec)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, catalog)
}

func addFederatedCatalogRoutes(r *gin.RouterGroup) {
	contractAgreements := r.Group("/federated-catalog")
	contractAgreements.POST("/", getFederatedCatalog)
}
