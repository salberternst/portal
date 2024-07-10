package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/salberternst/portal/pkg/api"
	"github.com/salberternst/portal/pkg/middleware"
)

func CreateContractNegotiation(ctx *gin.Context) {
	contractNegotiation := api.ContractRequest{}
	if err := ctx.BindJSON(&contractNegotiation); err != nil {
		RespondWithBadRequest(ctx, "Bad Request")
		return
	}

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

	ctx.JSON(http.StatusOK, contractNegotiation)
}

func TerminateContractNegotiation(ctx *gin.Context) {
	id := ctx.Param("id")

	terminateNegotiation := api.TerminateNegotiation{}
	if err := ctx.BindJSON(&terminateNegotiation); err != nil {
		RespondWithBadRequest(ctx, "Bad Request")
		return
	}

	err := middleware.GetEdcAPI(ctx).TerminateContractNegotiation(terminateNegotiation)
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
