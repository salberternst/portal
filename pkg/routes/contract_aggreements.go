package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/salberternst/portal/pkg/api"
	"github.com/salberternst/portal/pkg/middleware"
)

type ContractAgreementsQuery struct {
	Page     uint `form:"page" binding:"required"`
	PageSize uint `form:"page_size"  binding:"required"`
}

func GetContractAgreements(ctx *gin.Context) {
	query := ContractAgreementsQuery{}
	if err := ctx.BindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"error":   "bad_request",
			"message": fmt.Sprintf("unable to bind query: %v", err),
		})
		return
	}

	querySpec := api.QuerySpec{
		Context: map[string]string{
			"@vocab": "https://w3id.org/edc/v0.0.1/ns/",
		},
		Type: "QuerySpec",
		// FilterExpression: utils.BuildFilterExpressionFromContext(ctx),
	}

	assets, err := middleware.GetEdcAPI(ctx).GetAssets(querySpec)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"error":   "unable_to_get_assets",
			"message": fmt.Sprintf("unable to get assets: %v", err),
		})
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

	querySpec = api.QuerySpec{
		Context: map[string]string{
			"@vocab": "https://w3id.org/edc/v0.0.1/ns/",
		},
		Type: "QuerySpec",
		FilterExpression: []api.Criterion{
			{
				OperandLeft:  "assetId",
				Operator:     "in",
				OperandRight: assetIds,
			},
		},
	}

	contractAgreements, err := middleware.GetEdcAPI(ctx).GetContractAgreements(querySpec)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"error":   "unable_to_get_contract_agreements",
			"message": fmt.Sprintf("unable to get contract agreements: %v", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, contractAgreements)
}

func GetContractAgreement(ctx *gin.Context) {
	id := ctx.Param("id")

	contractAgreement, err := middleware.GetEdcAPI(ctx).GetContractAgreement(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"error":   "unable_to_get_contract_agreement",
			"message": fmt.Sprintf("unable to get contract agreement: %v", err),
		})
		return
	}

	// if contractAgreement.PrivateProperties == nil || !utils.CheckPrivateProperties(ctx, contractAgreement.PrivateProperties) {
	// 	ctx.JSON(http.StatusForbidden, gin.H{
	// 		"status":  http.StatusForbidden,
	// 		"error":   "forbidden",
	// 		"message": "You are not allowed to access this contract agreement",
	// 	})
	// 	return
	// }

	// todo: check if user is allowed to access this contract agreement
	// todo: admin can see all contract agreements

	ctx.JSON(http.StatusOK, contractAgreement)
}

func GetContractAgreementNegotiation(ctx *gin.Context) {
	id := ctx.Param("id")

	contractAgreementNegotiation, err := middleware.GetEdcAPI(ctx).GetContractAgreementNegotiation(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"error":   "unable_to_get_contract_agreement_negotiation",
			"message": fmt.Sprintf("unable to get contract agreement negotiation: %v", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, contractAgreementNegotiation)
}

func addContractAggreementsRoutes(r *gin.RouterGroup) {
	contractAgreements := r.Group("/contractagreements")
	contractAgreements.GET("/", GetContractAgreements)
	contractAgreements.GET("/:id", GetContractAgreement)
	contractAgreements.GET("/:id/negotiation", GetContractAgreementNegotiation)
}
