package routes

import (
	"maps"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/salberternst/portal/pkg/api"
	"github.com/salberternst/portal/pkg/middleware"
)

func getContractDefinitions(ctx *gin.Context) {
	querySpec, err := CreateQuerySpecFromContext(ctx)
	if err != nil {
		RespondWithBadRequest(ctx, "Bad Request")
		return
	}

	contracts, err := middleware.GetEdcAPI(ctx).GetContractDefinitions(querySpec)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, contracts)
}

func getContractDefinition(ctx *gin.Context) {
	contractId := ctx.Param("id")

	contract, err := middleware.GetEdcAPI(ctx).GetContractDefinition(contractId)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	if middleware.GetAuthenticatedUser(ctx).IsCustomer() {
		if contract.PrivateProperties == nil || !CheckPrivateProperties(ctx, contract.PrivateProperties) {
			RespondWithForbidden(ctx)
			return
		}
	}

	ctx.JSON(http.StatusOK, contract)
}

func deleteContractDefinition(ctx *gin.Context) {
	contractId := ctx.Param("id")

	contract, err := middleware.GetEdcAPI(ctx).GetContractDefinition(contractId)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	if middleware.GetAuthenticatedUser(ctx).IsCustomer() {
		if contract.PrivateProperties == nil || !CheckPrivateProperties(ctx, contract.PrivateProperties) {
			RespondWithForbidden(ctx)
			return
		}
	}

	err = middleware.GetEdcAPI(ctx).DeleteContractDefinition(contractId)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": contractId,
	})
}

func createContractDefinition(ctx *gin.Context) {
	contractDefinition := api.ContractDefinition{}
	if err := ctx.BindJSON(&contractDefinition); err != nil {
		RespondWithBadRequest(ctx, "Bad Request")
		return
	}

	contractDefinition.Id = uuid.New().String()
	privateProperties := BuildPrivatePropertiesFromContext(ctx)
	if privateProperties != nil {
		maps.Copy(contractDefinition.PrivateProperties, privateProperties)
	} else {
		contractDefinition.PrivateProperties = privateProperties
	}

	createdContractDefinition, err := middleware.GetEdcAPI(ctx).CreateContractDefinition(contractDefinition)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusCreated, createdContractDefinition)
}

func addContractDefinitionsRoutes(r *gin.RouterGroup) {
	contracts := r.Group("/contractdefinitions")
	contracts.GET("/", getContractDefinitions)
	contracts.GET("/:id", getContractDefinition)
	contracts.DELETE("/:id", deleteContractDefinition)
	contracts.POST("/", createContractDefinition)
}
