package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/salberternst/portal/pkg/api"
	"github.com/salberternst/portal/pkg/middleware"
)

type ContractQuery struct {
	Page     uint `form:"page" binding:"required"`
	PageSize uint `form:"page_size"  binding:"required"`
}

func GetContractDefinitions(ctx *gin.Context) {
	query := ContractQuery{}
	if err := ctx.BindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	querySpec := api.QuerySpec{
		Context: map[string]string{
			"@vocab": "https://w3id.org/edc/v0.0.1/ns/",
		},
		Type:   "QuerySpec",
		Offset: (query.Page - 1) * query.PageSize,
		Limit:  query.PageSize,
		// SortOrder: "DESC",
		// SortField: "https://w3id.org/edc/v0.0.1/ns/id",
		FilterExpression: []api.Criterion{
			{
				OperandLeft:  "privateProperties.'https://w3id.org/edc/v0.0.1/ns/createdBy'",
				Operator:     "=",
				OperandRight: middleware.GetAccessTokenClaims(ctx).Subject,
			},
		},
	}

	contracts, err := middleware.GetEdcAPI(ctx).GetContractDefinitions(querySpec)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"error":   "unable_to_get_contract_definitions",
			"message": fmt.Sprintf("unable to get contract definitions: %v", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, contracts)
}

func GetContractDefinition(ctx *gin.Context) {
	contractId := ctx.Param("id")

	contract, err := middleware.GetEdcAPI(ctx).GetContractDefinition(contractId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"error":   "unable_to_get_contract_definition",
			"message": fmt.Sprintf("unable to get contract definition: %v", err),
		})
		return
	}

	if contract.PrivateProperties == nil || contract.PrivateProperties["createdBy"] != middleware.GetAccessTokenClaims(ctx).Subject {
		ctx.JSON(http.StatusForbidden, gin.H{
			"status":  http.StatusForbidden,
			"error":   "forbidden",
			"message": "You are not allowed to access this contract definition",
		})
		return
	}

	ctx.JSON(http.StatusOK, contract)
}

func DeleteContractDefinition(ctx *gin.Context) {
	contractId := ctx.Param("id")

	contract, err := middleware.GetEdcAPI(ctx).GetContractDefinition(contractId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"error":   "unable_to_get_contract_definition",
			"message": fmt.Sprintf("unable to get contract definition: %v", err),
		})
		return
	}

	if contract.PrivateProperties == nil || contract.PrivateProperties["createdBy"] != middleware.GetAccessTokenClaims(ctx).Subject {
		ctx.JSON(http.StatusForbidden, gin.H{
			"status":  http.StatusForbidden,
			"error":   "forbidden",
			"message": "You are not allowed to delete this contract definition",
		})
		return
	}

	err = middleware.GetEdcAPI(ctx).DeleteContractDefinition(contractId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"error":   "unable_to_delete_contract_definition",
			"message": fmt.Sprintf("unable to delete contract definition: %v", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": contractId,
	})
}

func CreateContractDefinition(ctx *gin.Context) {
	contractDefinition := api.ContractDefinition{}
	if err := ctx.BindJSON(&contractDefinition); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"error":   "bad_request",
			"message": fmt.Sprintf("unable to bind contract definition: %v", err),
		})
		return
	}

	contractDefinition.Id = uuid.New().String()
	if contractDefinition.PrivateProperties == nil {
		contractDefinition.PrivateProperties = map[string]string{}
	}
	contractDefinition.PrivateProperties["createdBy"] = middleware.GetAccessTokenClaims(ctx).Subject

	createdContractDefinition, err := middleware.GetEdcAPI(ctx).CreateContractDefinition(contractDefinition)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"error":   "unable_to_create_contract_definition",
			"message": fmt.Sprintf("unable to create contract definition: %v", err),
		})
		return
	}

	ctx.JSON(http.StatusCreated, createdContractDefinition)
}

func addContractDefinitionsRoutes(r *gin.RouterGroup) {
	contracts := r.Group("/contractdefinitions")
	contracts.GET("/", GetContractDefinitions)
	contracts.GET("/:id", GetContractDefinition)
	contracts.DELETE("/:id", DeleteContractDefinition)
	contracts.POST("/", CreateContractDefinition)
}
