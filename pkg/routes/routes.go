package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/salberternst/portal/pkg/middleware"
	"github.com/salberternst/portal/pkg/utils"
)

func AddRoutes(r *gin.Engine) {
	api := r.Group("/api/portal")

	addHealthRoutes(api)

	// only use the middlewares in the api group
	api.Use(middleware.TokenMiddleware())
	// api.Use(middleware.KeycloakMiddleware())
	// api.Use(middleware.ThingsboardMiddleware())
	// api.Use(middleware.FusekiMiddleware())
	api.Use(middleware.EdcMiddleware())
	api.Use(middleware.FederatedCatalogMiddleware())

	addCustomersRoutes(api)
	addUsersRoute(api)
	addAssetsRoutes(api)
	addPoliciesRoutes(api)
	addContractDefinitionsRoutes(api)
	addContractNegotationsRoutes(api)
	addContractAggreementsRoutes(api)
	addCatalogsRoutes(api)
	addTransferProcessesRoutes(api)
	addFederatedCatalogRoutes(api)
	if utils.GetConfig().EnableDeviceApi {
		addDevicesRoutes(api)
	}
}
