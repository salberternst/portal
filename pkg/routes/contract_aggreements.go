package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/salberternst/portal/pkg/api"
	"github.com/salberternst/portal/pkg/middleware"
)

func getContractAgreements(ctx *gin.Context) {
	query := QueryParams{}
	if err := ctx.BindQuery(&query); err != nil {
		RespondWithBadRequest(ctx, "Bad Request")
		return
	}

	querySpec, err := CreateQuerySpecWithoutFilterFromContext(ctx)
	if err != nil {
		RespondWithBadRequest(ctx, "Bad Request")
		return
	}

	if middleware.IsCustomer(ctx) {
		// create QuerySpec to filter assets that belong to the customer
		querySpec, err = CreateQuerySpecFromContext(ctx)
		if err != nil {
			RespondWithBadRequest(ctx, "Bad Request")
			return
		}

		// only return assets that belong to the customer
		assets, err := middleware.GetEdcAPI(ctx).GetAssets(querySpec)
		if err != nil {
			RespondWithInternalServerError(ctx)
			return
		}

		if len(assets) == 0 {
			ctx.JSON(http.StatusOK, []api.ContractAgreement{})
			return
		}

		assetIds := []string{}
		for _, asset := range assets {
			assetIds = append(assetIds, asset.Id)
		}

		querySpec.FilterExpression = []api.Criterion{
			{
				OperandLeft:  "assetId",
				Operator:     "in",
				OperandRight: assetIds,
			},
		}
	}

	contractAgreements, err := middleware.GetEdcAPI(ctx).GetContractAgreements(querySpec)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, contractAgreements)
}

func getContractAgreement(ctx *gin.Context) {
	id := ctx.Param("id")

	contractAgreement, err := middleware.GetEdcAPI(ctx).GetContractAgreement(id)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	if middleware.IsCustomer(ctx) {
		asset, err := middleware.GetEdcAPI(ctx).GetAsset(contractAgreement.AssetId)
		if err != nil {
			RespondWithResourceNotFound(ctx, id)
			return
		}

		if !CheckPrivateProperties(ctx, asset.PrivateProperties) {
			RespondWithResourceNotFound(ctx, id)
			return
		}
	}

	ctx.JSON(http.StatusOK, contractAgreement)
}

func getContractAgreementNegotiation(ctx *gin.Context) {
	id := ctx.Param("id")

	contractAgreementNegotiation, err := middleware.GetEdcAPI(ctx).GetContractAgreementNegotiation(id)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	if middleware.IsCustomer(ctx) {
		contractAgreement, err := middleware.GetEdcAPI(ctx).GetContractAgreement(contractAgreementNegotiation.ContractAgreementId)
		if err != nil {
			RespondWithResourceNotFound(ctx, id)
			return
		}

		asset, err := middleware.GetEdcAPI(ctx).GetAsset(contractAgreement.AssetId)
		if err != nil {
			RespondWithResourceNotFound(ctx, id)
			return
		}

		if !CheckPrivateProperties(ctx, asset.PrivateProperties) {
			RespondWithResourceNotFound(ctx, id)
			return
		}
	}

	ctx.JSON(http.StatusOK, contractAgreementNegotiation)
}

func addContractAggreementsRoutes(r *gin.RouterGroup) {
	contractAgreements := r.Group("/contractagreements")
	contractAgreements.GET("/", getContractAgreements)
	contractAgreements.GET("/:id", getContractAgreement)
	contractAgreements.GET("/:id/negotiation", getContractAgreementNegotiation)
}
