package routes

import (
	"maps"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/salberternst/portal/pkg/api"
	"github.com/salberternst/portal/pkg/middleware"
)

func getPolicies(ctx *gin.Context) {
	querySpec, err := CreateQuerySpecFromContext(ctx)
	if err != nil {
		RespondWithBadRequest(ctx, "Bad Request")
		return
	}

	policies, err := middleware.GetEdcAPI(ctx).GetPolicies(querySpec)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, policies)
}

func getPolicy(ctx *gin.Context) {
	id := ctx.Param("id")

	policyDefinition, err := middleware.GetEdcAPI(ctx).GetPolicy(id)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	if middleware.IsCustomer(ctx) {
		if policyDefinition.PrivateProperties == nil || !CheckPrivateProperties(ctx, policyDefinition.PrivateProperties) {
			RespondWithForbidden(ctx)
			return
		}
	}

	ctx.JSON(http.StatusOK, policyDefinition)
}

func createPolicy(ctx *gin.Context) {
	policyDefinition := api.PolicyDefinition{
		PrivateProperties: make(map[string]string),
	}

	if err := ctx.BindJSON(&policyDefinition); err != nil {
		RespondWithBadRequest(ctx, "Bad Request")
		return
	}

	policyDefinition.ID = uuid.New().String()
	privateProperties := BuildPrivatePropertiesFromContext(ctx)
	if privateProperties != nil {
		maps.Copy(policyDefinition.PrivateProperties, privateProperties)
	} else {
		policyDefinition.PrivateProperties = privateProperties
	}

	createdPolicy, err := middleware.GetEdcAPI(ctx).CreatePolicy(policyDefinition)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusCreated, createdPolicy)
}

func deletePolicy(ctx *gin.Context) {
	id := ctx.Param("id")

	policyDefinition, err := middleware.GetEdcAPI(ctx).GetPolicy(id)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	if middleware.IsCustomer(ctx) {
		if policyDefinition.PrivateProperties == nil || !CheckPrivateProperties(ctx, policyDefinition.PrivateProperties) {
			RespondWithForbidden(ctx)
			return
		}
	}

	if err = middleware.GetEdcAPI(ctx).DeletePolicy(id); err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func addPoliciesRoutes(r *gin.RouterGroup) {
	policiesGroup := r.Group("/policies")
	policiesGroup.GET("/", getPolicies)
	policiesGroup.GET("/:id", getPolicy)
	policiesGroup.POST("/", createPolicy)
	policiesGroup.DELETE("/:id", deletePolicy)
}
