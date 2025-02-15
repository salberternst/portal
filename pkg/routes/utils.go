package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/salberternst/portal/pkg/api"
	"github.com/salberternst/portal/pkg/middleware"
)

func RespondWithError(ctx *gin.Context, status int, errorCode, message string) {
	ctx.JSON(status, gin.H{
		"status":  status,
		"error":   errorCode,
		"message": message,
	})
}

func RespondWithInternalServerError(ctx *gin.Context) {
	RespondWithError(ctx, 500, "internal_server_error", "Internal Server Error")
}

func RespondWithBadRequest(ctx *gin.Context, message string) {
	RespondWithError(ctx, 400, "bad_request", message)
}

func RespondWithResourceNotFound(ctx *gin.Context, id string) {
	RespondWithError(ctx, 404, "not_found", "Resource "+id+" not found")
}

func RespondWithForbidden(ctx *gin.Context) {
	RespondWithError(ctx, 403, "forbidden", "Forbidden")
}

func CreateQuerySpec() api.QuerySpec {
	return api.QuerySpec{
		Context: map[string]string{
			"@vocab": "https://w3id.org/edc/v0.0.1/ns/",
		},
		Type: "QuerySpec",
	}
}

func CreateQuerySpecFromContext(ctx *gin.Context) (api.QuerySpec, error) {
	queryParams := QueryParams{}
	if err := ctx.BindQuery(&queryParams); err != nil {
		return api.QuerySpec{}, err
	}

	return api.QuerySpec{
		Context: map[string]string{
			"@vocab": "https://w3id.org/edc/v0.0.1/ns/",
		},
		Type:             "QuerySpec",
		Offset:           (queryParams.Page - 1) * queryParams.PageSize,
		Limit:            queryParams.PageSize,
		FilterExpression: BuildFilterExpressionFromContext(ctx),
	}, nil
}

func CreateQuerySpecWithoutFilterFromContext(ctx *gin.Context) (api.QuerySpec, error) {
	queryParams := QueryParams{}
	if err := ctx.BindQuery(&queryParams); err != nil {
		return api.QuerySpec{}, err
	}

	return api.QuerySpec{
		Context: map[string]string{
			"@vocab": "https://w3id.org/edc/v0.0.1/ns/",
		},
		Type:      "QuerySpec",
		Offset:    (queryParams.Page - 1) * queryParams.PageSize,
		Limit:     queryParams.PageSize,
		SortOrder: queryParams.SortOrder,
		SortField: queryParams.SortField,
	}, nil
}

func BuildFilterExpressionFromContext(ctx *gin.Context) []api.Criterion {
	filterExpressions := []api.Criterion{}

	if middleware.GetAuthenticatedUser(ctx).IsCustomer() {
		filterExpressions = append(filterExpressions, api.Criterion{
			OperandLeft:  "privateProperties.'https://w3id.org/edc/v0.0.1/ns/tenantId'",
			Operator:     "=",
			OperandRight: middleware.GetAuthenticatedUser(ctx).TenantId,
		})

		filterExpressions = append(filterExpressions, api.Criterion{
			OperandLeft:  "privateProperties.'https://w3id.org/edc/v0.0.1/ns/customerId'",
			Operator:     "=",
			OperandRight: middleware.GetAuthenticatedUser(ctx).CustomerId,
		})
	}

	return filterExpressions
}

func BuildPrivatePropertiesFromContext(ctx *gin.Context) map[string]string {
	privateProperties := make(map[string]string)
	privateProperties["https://w3id.org/edc/v0.0.1/ns/tenantId"] = middleware.GetAuthenticatedUser(ctx).TenantId

	if middleware.GetAuthenticatedUser(ctx).IsCustomer() {
		privateProperties["https://w3id.org/edc/v0.0.1/ns/customerId"] = middleware.GetAuthenticatedUser(ctx).CustomerId
	}

	return privateProperties
}

func CheckPrivateProperties(ctx *gin.Context, privateProperties map[string]string) bool {
	if middleware.GetAuthenticatedUser(ctx).IsCustomer() {
		if privateProperties["https://w3id.org/edc/v0.0.1/ns/customerId"] != middleware.GetAuthenticatedUser(ctx).CustomerId &&
			privateProperties["customerId"] != middleware.GetAuthenticatedUser(ctx).CustomerId {
			return false
		}

		if privateProperties["https://w3id.org/edc/v0.0.1/ns/tenantId"] != middleware.GetAuthenticatedUser(ctx).TenantId &&
			privateProperties["tenantId"] != middleware.GetAuthenticatedUser(ctx).TenantId {
			return false
		}
	}

	return true
}
