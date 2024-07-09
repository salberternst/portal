package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/salberternst/portal/pkg/api"
	"github.com/salberternst/portal/pkg/middleware"
)

func getTransferProcesses(ctx *gin.Context) {
	if !middleware.IsCustomer(ctx) && !middleware.IsAdmin(ctx) {
		RespondWithForbidden(ctx)
		return
	}

	// this works on the provider side as we get all
	querySpec, err := CreateQuerySpecFromContext(ctx)
	if err != nil {
		RespondWithBadRequest(ctx, "Bad Request")
		return
	}

	assets, err := middleware.GetEdcAPI(ctx).GetAssets(querySpec)
	if err != nil {
		RespondWithInternalServerError(ctx)
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

	transferProcesses, err := middleware.GetEdcAPI(ctx).GetTransferProcesses(querySpec)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, transferProcesses)
}

func getTransferProcess(ctx *gin.Context) {
	if !middleware.IsCustomer(ctx) && !middleware.IsAdmin(ctx) {
		RespondWithForbidden(ctx)
		return
	}

	transferProcess, err := middleware.GetEdcAPI(ctx).GetTransferProcess(ctx.Param("id"))
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	if transferProcess.Type == "PROVIDER" {
		asset, err := middleware.GetEdcAPI(ctx).GetAsset(transferProcess.AssetId)
		if err != nil {
			RespondWithInternalServerError(ctx)
			return
		}

		if !CheckPrivateProperties(ctx, asset.PrivateProperties) {
			RespondWithForbidden(ctx)
			return
		}
	} else {
		// EDC currently does not return the privateProperties for the TransferProcess
		// Therefore, we need to check the privateProperties of the asset using the /request endpoint
		querySpec := CreateQuerySpec()
		querySpec.FilterExpression = BuildFilterExpressionFromContext(ctx)
		querySpec.FilterExpression = append(querySpec.FilterExpression, api.Criterion{
			OperandLeft:  "assetId",
			Operator:     "=",
			OperandRight: transferProcess.AssetId,
		})

		transferProcesses, err := middleware.GetEdcAPI(ctx).GetTransferProcesses(querySpec)
		if err != nil {
			RespondWithInternalServerError(ctx)
			return
		}

		if len(transferProcesses) == 0 {
			RespondWithForbidden(ctx)
			return
		}

		transferProcess = &transferProcesses[0]
	}

	ctx.JSON(http.StatusOK, transferProcess)
}

func createTransferProcess(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		RespondWithForbidden(ctx)
		return
	}

	var transferRequest api.TransferRequest
	if err := ctx.BindJSON(&transferRequest); err != nil {
		RespondWithBadRequest(ctx, "Bad Request")
		return
	}

	transferRequest.PrivateProperties = BuildPrivatePropertiesFromContext(ctx)

	transferProcess, err := middleware.GetEdcAPI(ctx).CreateTransferProcess(transferRequest)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusCreated, transferProcess)
}

func addTransferProcessesRoutes(r *gin.RouterGroup) {
	transferProcesses := r.Group("/transferprocesses")
	transferProcesses.GET("/", getTransferProcesses)
	transferProcesses.GET("/:id", getTransferProcess)
	transferProcesses.POST("/", createTransferProcess)
}
