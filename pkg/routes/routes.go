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
	api.Use(middleware.JwtMiddleware())
	api.Use(middleware.AuthMiddleware())

	if utils.GetConfig().EnableUsersApi {
		api.Use(middleware.KeycloakMiddleware())
	}

	if utils.GetConfig().EnableDeviceApi {
		api.Use(middleware.ThingsboardMiddleware())
	}

	if utils.GetConfig().EnableFusekiBackend {
		api.Use(middleware.FusekiMiddleware())
	}

	api.Use(middleware.EdcMiddleware())
	api.Use(middleware.FederatedCatalogMiddleware())

	if utils.GetConfig().EnableUsersApi {
		addCustomersRoutes(api)
		addUsersRoute(api)
	}

	addAssetsRoutes(api)
	addPoliciesRoutes(api)
	addContractDefinitionsRoutes(api)
	addContractNegotationsRoutes(api)
	addContractAggreementsRoutes(api)
	addCatalogsRoutes(api)
	addTransferProcessesRoutes(api)
	addFederatedCatalogRoutes(api)
	addUserinfoRoutes(api)

	if utils.GetConfig().EnableDeviceApi {
		addDevicesRoutes(api)
	}
}
