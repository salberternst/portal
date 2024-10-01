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

	querySpec, err := CreateQuerySpecWithoutFilterFromContext(ctx)
	if err != nil {
		RespondWithBadRequest(ctx, "Bad Request")
		return
	}

	if middleware.IsCustomer(ctx) {
		// this works on the provider side as we get all
		querySpec, err = CreateQuerySpecFromContext(ctx)
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
			{
				OperandLeft:  "type",
				Operator:     "=",
				OperandRight: "PROVIDER",
			},
		}
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

	id := ctx.Param("id")

	transferProcess, err := middleware.GetEdcAPI(ctx).GetTransferProcess(id)
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
			OperandLeft:  "id",
			Operator:     "=",
			OperandRight: id,
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

func getTransferProcessDataRequest(ctx *gin.Context) {
	if !middleware.IsCustomer(ctx) && !middleware.IsAdmin(ctx) {
		RespondWithForbidden(ctx)
		return
	}

	id := ctx.Param("id")

	transferProcess, err := middleware.GetEdcAPI(ctx).GetTransferProcess(id)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	if transferProcess.Type == "CONSUMER" {
		dataRequest, err := middleware.GetEdcAPI(ctx).GetEdrDataAddress(id)
		if err != nil {
			RespondWithInternalServerError(ctx)
			return
		}

		ctx.JSON(http.StatusOK, dataRequest)
		return
	}

	RespondWithBadRequest(ctx, "Bad Request")
}

func terminateTransferProcess(ctx *gin.Context) {
	terminateTransferProcess := api.TerminateTransferProcess{}
	if err := ctx.BindJSON(&terminateTransferProcess); err != nil {
		RespondWithBadRequest(ctx, "Bad Request")
		return
	}

	id := ctx.Param("id")

	transferProcess, err := middleware.GetEdcAPI(ctx).GetTransferProcess(id)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	if transferProcess.Type == "CONSUMER" {
		if !middleware.IsAdmin(ctx) {
			RespondWithForbidden(ctx)
			return
		}
	} else {
		// provider contracts can only be terminated by the owner of the asset
		asset, err := middleware.GetEdcAPI(ctx).GetAsset(transferProcess.AssetId)
		if err != nil {
			RespondWithInternalServerError(ctx)
			return
		}

		if !CheckPrivateProperties(ctx, asset.PrivateProperties) {
			RespondWithForbidden(ctx)
			return
		}
	}

	err = middleware.GetEdcAPI(ctx).TerminateTransferProcess(terminateTransferProcess)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}
func getDataConsumerPull(ctx *gin.Context) {
	if !middleware.IsCustomer(ctx) && !middleware.IsAdmin(ctx) {
		RespondWithForbidden(ctx)
		return
	}

	id := ctx.Param("id")
	transferProcess, err := middleware.GetEdcAPI(ctx).GetTransferProcess(id)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	if transferProcess.Type == "CONSUMER" {
		DataAddress, err := middleware.GetEdcAPI(ctx).GetEdrDataAddress(id)
		if err != nil {
			RespondWithInternalServerError(ctx)
			return
		}
		DataConsumerPull, err := middleware.GetEdcAPI(ctx).GetDataConsumerPull(*DataAddress)
		if err != nil {
			RespondWithInternalServerError(ctx)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": DataConsumerPull})
		return
	}

	RespondWithBadRequest(ctx, "Bad Request")
}

func getRawDataConsumerPull(ctx *gin.Context) {
	if !middleware.IsCustomer(ctx) && !middleware.IsAdmin(ctx) {
		RespondWithForbidden(ctx)
		return
	}

	id := ctx.Param("id")
	transferProcess, err := middleware.GetEdcAPI(ctx).GetTransferProcess(id)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	if transferProcess.Type == "CONSUMER" {
		DataAddress, err := middleware.GetEdcAPI(ctx).GetEdrDataAddress(id)
		if err != nil {
			RespondWithInternalServerError(ctx)
			return
		}
		RawData, err := middleware.GetEdcAPI(ctx).GetRawDataConsumerPull(*DataAddress)
		if err != nil {
			RespondWithInternalServerError(ctx)
			return
		}
		ctx.JSON(http.StatusOK, RawData)
		return
	}

	RespondWithBadRequest(ctx, "Bad Request")
}

func addTransferProcessesRoutes(r *gin.RouterGroup) {
	transferProcesses := r.Group("/transferprocesses")
	transferProcesses.GET("/", getTransferProcesses)
	transferProcesses.GET("/:id", getTransferProcess)
	transferProcesses.GET("/:id/datarequest", getTransferProcessDataRequest)
	transferProcesses.GET("/:id/data", getDataConsumerPull)
	transferProcesses.GET("/:id/download", getRawDataConsumerPull)
	transferProcesses.POST("/:id/terminate", terminateTransferProcess)
	transferProcesses.POST("/", createTransferProcess)
}
