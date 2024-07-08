package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/salberternst/portal/pkg/api"
	"github.com/salberternst/portal/pkg/middleware"
)

type TransferProcessQuery struct {
	Page     uint `form:"page" binding:"required"`
	PageSize uint `form:"page_size"  binding:"required"`
}

func GetTransferProcesses(ctx *gin.Context) {
	if !middleware.IsCustomer(ctx) && !middleware.IsAdmin(ctx) {
		ctx.JSON(http.StatusForbidden, gin.H{
			"status":  http.StatusForbidden,
			"error":   "forbidden",
			"message": "you are not allowed to access transfer processes",
		})
		return
	}

	transferProcessQuery := TransferProcessQuery{}
	if err := ctx.BindQuery(&transferProcessQuery); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"error":   "bad_request",
			"message": fmt.Sprintf("invalid query parameters: %v", err),
		})
		return
	}

	querySpec := api.QuerySpec{
		Context: map[string]string{
			"@vocab": "https://w3id.org/edc/v0.0.1/ns/",
		},
		Type:   "QuerySpec",
		Offset: (transferProcessQuery.Page - 1) * transferProcessQuery.PageSize,
		Limit:  transferProcessQuery.PageSize,
	}

	if middleware.IsCustomer(ctx) {
		// the customer sees only that transfer processes that are related to his assets
		assetQuery := api.QuerySpec{
			Context: map[string]string{
				"@vocab": "https://w3id.org/edc/v0.0.1/ns/",
			},
			Type:   "QuerySpec",
			Offset: 0,
			Limit:  1000,
			FilterExpression: []api.Criterion{
				{
					OperandLeft:  "privateProperties.'https://w3id.org/edc/v0.0.1/ns/tenantId'",
					Operator:     "=",
					OperandRight: middleware.GetAccessTokenClaims(ctx).TenantId,
				},
				{
					OperandLeft:  "privateProperties.'https://w3id.org/edc/v0.0.1/ns/customerId'",
					Operator:     "=",
					OperandRight: middleware.GetAccessTokenClaims(ctx).CustomerId,
				},
			},
		}

		assets, err := middleware.GetEdcAPI(ctx).GetAssets(assetQuery)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"error":   "unable_to_get_assets",
				"message": fmt.Sprintf("unable to get assets: %v", err),
			})
			return
		}

		assetIds := make([]string, len(assets))
		for i, asset := range assets {
			assetIds[i] = asset.Id
		}

		if len(assets) == 0 {
			ctx.JSON(http.StatusOK, []api.TransferProcess{})
			return
		}

		querySpec.FilterExpression = []api.Criterion{
			{
				OperandLeft:  "assetId",
				Operator:     "in",
				OperandRight: assetIds,
			},
		}
	}

	transferProcesses, err := middleware.GetEdcAPI(ctx).GetTransferProcesses(querySpec)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"error":   "unable_to_get_transfer_processes",
			"message": fmt.Sprintf("unable to get transfer processes: %v", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, transferProcesses)
}

func GetTransferProcess(ctx *gin.Context) {
	if !middleware.IsCustomer(ctx) && !middleware.IsAdmin(ctx) {
		ctx.JSON(http.StatusForbidden, gin.H{
			"status":  http.StatusForbidden,
			"error":   "forbidden",
			"message": "you are not allowed to access transfer processes",
		})
		return
	}

	transferProcess, err := middleware.GetEdcAPI(ctx).GetTransferProcess(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"error":   "unable_to_get_transfer_process",
			"message": fmt.Sprintf("unable to get transfer process: %v", err),
		})
		return
	}

	if transferProcess.PrivateProperties == nil ||
		transferProcess.PrivateProperties["https://w3id.org/edc/v0.0.1/ns/tenantId"] != middleware.GetAccessTokenClaims(ctx).TenantId {
		ctx.JSON(http.StatusForbidden, gin.H{
			"status":  http.StatusForbidden,
			"error":   "forbidden",
			"message": "you are not allowed to access this transfer process",
		})
		return
	}

	if middleware.IsCustomer(ctx) {
		contractAgreement, err := middleware.GetEdcAPI(ctx).GetContractAgreement(transferProcess.ContractAgreementId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"error":   "unable_to_get_contract_agreement",
				"message": fmt.Sprintf("unable to get contract agreement: %v", err),
			})
			return
		}

		asset, err := middleware.GetEdcAPI(ctx).GetAsset(contractAgreement.AssetId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"error":   "unable_to_get_asset",
				"message": fmt.Sprintf("unable to get asset: %v", err),
			})
			return
		}

		if asset.PrivateProperties["https://w3id.org/edc/v0.0.1/ns/tenantId"] != middleware.GetAccessTokenClaims(ctx).TenantId ||
			asset.PrivateProperties["https://w3id.org/edc/v0.0.1/ns/customerId"] != middleware.GetAccessTokenClaims(ctx).CustomerId {
			ctx.JSON(http.StatusForbidden, gin.H{
				"status":  http.StatusForbidden,
				"error":   "forbidden",
				"message": "you are not allowed to access this transfer process",
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, transferProcess)
}

func CreateTransferProcess(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		ctx.JSON(http.StatusForbidden, gin.H{
			"status":  http.StatusForbidden,
			"error":   "forbidden",
			"message": "you are not allowed to create transfer processes",
		})
		return
	}

	var transferRequest api.TransferRequest
	if err := ctx.BindJSON(&transferRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"error":   "invalid_transfer_process",
			"message": fmt.Sprintf("invalid transfer process: %v", err),
		})
		return
	}

	if transferRequest.PrivateProperties == nil {
		transferRequest.PrivateProperties = map[string]string{}
	}
	transferRequest.PrivateProperties["tenantId"] = middleware.GetAccessTokenClaims(ctx).TenantId

	transferProcess, err := middleware.GetEdcAPI(ctx).CreateTransferProcess(transferRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"error":   "unable_to_create_transfer_process",
			"message": fmt.Sprintf("unable to create transfer process: %v", err),
		})
		return
	}

	ctx.JSON(http.StatusCreated, transferProcess)
}

func addTransferProcessesRoutes(r *gin.RouterGroup) {
	transferProcesses := r.Group("/transferprocesses")
	transferProcesses.GET("/", GetTransferProcesses)
	transferProcesses.GET("/:id", GetTransferProcess)
	transferProcesses.POST("/", CreateTransferProcess)
}
