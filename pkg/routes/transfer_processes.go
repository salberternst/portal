package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/salberternst/fsn_landing_page/pkg/api"
	"github.com/salberternst/fsn_landing_page/pkg/middleware"
)

func GetTransferProcesses(ctx *gin.Context) {
	querySpec := api.QuerySpec{
		Context: map[string]string{
			"@vocab": "https://w3id.org/edc/v0.0.1/ns/",
		},
		Type: "QuerySpec",
		FilterExpression: []api.Criterion{
			{
				OperandLeft:  "privateProperties.'https://w3id.org/edc/v0.0.1/ns/createdBy'",
				Operator:     "=",
				OperandRight: middleware.GetAccessTokenClaims(ctx).Subject,
			},
		},
	}

	transferProcesses, err := middleware.GetEdcAPI(ctx).GetTransferProcesses(querySpec)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"error":   "unable_to_get_transfer_processes",
			"message": fmt.Sprintf("unable to get transfer processes: %v", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, transferProcesses)
}

func GetTransferProcess(ctx *gin.Context) {
	transferProcessID := ctx.Param("id")

	transferProcess, err := middleware.GetEdcAPI(ctx).GetTransferProcess(transferProcessID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"error":   "unable_to_get_transfer_process",
			"message": fmt.Sprintf("unable to get transfer process: %v", err),
		})
		return
	}

	// if transferProcess.PrivateProperties["createdBy"] != claims.Subject {
	// 	ctx.JSON(http.StatusForbidden, gin.H{
	// 		"status":  http.StatusForbidden,
	// 		"error":   "forbidden",
	// 		"message": "you are not allowed to access this transfer process",
	// 	})
	// 	return
	// }

	ctx.JSON(http.StatusOK, transferProcess)
}

func CreateTransferProcess(ctx *gin.Context) {
	var transferRequest api.TransferRequest
	if err := ctx.BindJSON(&transferRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"error":   "invalid_transfer_process",
			"message": fmt.Sprintf("invalid transfer process: %v", err),
		})
		return
	}

	if transferRequest.PrivateProperties == nil {
		transferRequest.PrivateProperties = map[string]string{}
	}
	transferRequest.PrivateProperties["createdBy"] = middleware.GetAccessTokenClaims(ctx).Subject

	transferProcess, err := middleware.GetEdcAPI(ctx).CreateTransferProcess(transferRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"error":   "unable_to_create_transfer_process",
			"message": fmt.Sprintf("unable to create transfer process: %v", err),
		})
		return
	}

	ctx.JSON(http.StatusCreated, transferProcess)
}

func addTransferProcessesRoutes(r *gin.RouterGroup) {
	transferProcesses := r.Group("/transferprocesses")
	transferProcesses.GET("/", GetTransferProcesses)
	transferProcesses.GET("/:id", GetTransferProcess)
	transferProcesses.POST("/", CreateTransferProcess)
}
