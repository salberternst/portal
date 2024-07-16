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

	// currently no filter is applied
	// all users can create contract negotiations
	// todo: admins can create the negotiations but there is
	// currently no way to assign a user to a negotiation
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

	// currently no filter is applied
	// all users can see all contract negotiations

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

	// currently no filter is applied
	// all users can terminate all contract negotiations

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
