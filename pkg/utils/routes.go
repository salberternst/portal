package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/salberternst/portal/pkg/api"
	"github.com/salberternst/portal/pkg/middleware"
)

func BuildFilterExpressionFromContext(ctx *gin.Context) []api.Criterion {
	filterExpressions := []api.Criterion{
		{
			OperandLeft:  "privateProperties.'https://w3id.org/edc/v0.0.1/ns/tenantId'",
			Operator:     "=",
			OperandRight: middleware.GetAccessTokenClaims(ctx).TenantId,
		},
	}

	if middleware.IsCustomer(ctx) {
		filterExpressions = append(filterExpressions, api.Criterion{
			OperandLeft:  "privateProperties.'https://w3id.org/edc/v0.0.1/ns/customerId'",
			Operator:     "=",
			OperandRight: middleware.GetAccessTokenClaims(ctx).CustomerId,
		})
	}

	return filterExpressions
}

func BuildPrivatePropertiesFromContext(ctx *gin.Context) map[string]string {
	privateProperties := map[string]string{
		"https://w3id.org/edc/v0.0.1/ns/tenantId": middleware.GetAccessTokenClaims(ctx).TenantId,
	}

	if middleware.IsCustomer(ctx) {
		privateProperties["https://w3id.org/edc/v0.0.1/ns/customerId"] = middleware.GetAccessTokenClaims(ctx).CustomerId
	}

	return privateProperties
}

func CheckPrivateProperties(ctx *gin.Context, privateProperties map[string]string) bool {
	if middleware.IsCustomer(ctx) {
		if privateProperties["https://w3id.org/edc/v0.0.1/ns/customerId"] != middleware.GetAccessTokenClaims(ctx).CustomerId &&
			privateProperties["customerId"] != middleware.GetAccessTokenClaims(ctx).CustomerId {
			return false
		}
	}

	if privateProperties["https://w3id.org/edc/v0.0.1/ns/tenantId"] != middleware.GetAccessTokenClaims(ctx).TenantId &&
		privateProperties["tenantId"] != middleware.GetAccessTokenClaims(ctx).TenantId {
		return false
	}

	return true
}
