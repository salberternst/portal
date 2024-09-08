package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/salberternst/portal/pkg/api"
	"github.com/salberternst/portal/pkg/middleware"
)

func CreateContractNegotiation(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		RespondWithForbidden(ctx)
		return
	}

	contractNegotiation := api.ContractRequest{}
	if err := ctx.BindJSON(&contractNegotiation); err != nil {
		RespondWithBadRequest(ctx, "Bad Request")
		return
	}

	// currently no way to direclty assign a user to a negotiation
	// we could potentially use a callback to assign the user to the agreement

	createdContractDefinition, err := middleware.GetEdcAPI(ctx).CreateContractNegotiation(contractNegotiation)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusCreated, createdContractDefinition)
}

func GetContractNegotiation(ctx *gin.Context) {
	id := ctx.Param("id")

	contractNegotiation, err := middleware.GetEdcAPI(ctx).GetContractNegotiation(id)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	// consumer contracts can only be viewed by admins
	if contractNegotiation.Type == "CONSUMER" {
		if !middleware.IsAdmin(ctx) {
			RespondWithForbidden(ctx)
			return
		}
	} else {
		// provider contracts can only be viewed by the owner of the asset
		contractAggreement, err := middleware.GetEdcAPI(ctx).GetContractAgreement(contractNegotiation.ContractAgreementId)
		if err != nil {
			RespondWithInternalServerError(ctx)
			return
		}

		asset, err := middleware.GetEdcAPI(ctx).GetAsset(contractAggreement.AssetId)
		if err != nil {
			RespondWithInternalServerError(ctx)
			return
		}

		if !CheckPrivateProperties(ctx, asset.PrivateProperties) {
			RespondWithForbidden(ctx)
			return
		}
	}

	ctx.JSON(http.StatusOK, contractNegotiation)
}

func TerminateContractNegotiation(ctx *gin.Context) {
	id := ctx.Param("id")

	terminateNegotiation := api.TerminateNegotiation{}
	if err := ctx.BindJSON(&terminateNegotiation); err != nil {
		RespondWithBadRequest(ctx, "Bad Request")
		return
	}

	contractNegotiation, err := middleware.GetEdcAPI(ctx).GetContractNegotiation(id)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	if contractNegotiation.Type == "CONSUMER" {
		if !middleware.IsAdmin(ctx) {
			RespondWithForbidden(ctx)
			return
		}
	} else {
		// provider contracts can only be terminated by the owner of the asset
		contractAggreement, err := middleware.GetEdcAPI(ctx).GetContractAgreement(contractNegotiation.ContractAgreementId)
		if err != nil {
			RespondWithInternalServerError(ctx)
			return
		}

		asset, err := middleware.GetEdcAPI(ctx).GetAsset(contractAggreement.AssetId)
		if err != nil {
			RespondWithInternalServerError(ctx)
			return
		}

		if !CheckPrivateProperties(ctx, asset.PrivateProperties) {
			RespondWithForbidden(ctx)
			return
		}
	}

	err = middleware.GetEdcAPI(ctx).TerminateContractNegotiation(terminateNegotiation)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func addContractNegotationsRoutes(r *gin.RouterGroup) {
	contractNegotiations := r.Group("/contractnegotiations")
	contractNegotiations.POST("/", CreateContractNegotiation)
	contractNegotiations.GET("/:id", GetContractNegotiation)
	contractNegotiations.POST("/:id/terminate", TerminateContractNegotiation)
}
