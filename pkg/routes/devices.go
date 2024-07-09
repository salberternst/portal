package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/salberternst/portal/pkg/api"
	"github.com/salberternst/portal/pkg/middleware"
)

type CreateDevice struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Gateway     bool   `json:"gateway,omitempty"`
	ThingModel  string `json:"thingModel,omitempty"`
	CustomerId  string `json:"customerId,omitempty"`
}

type UpdateDevice struct {
	Description *string `json:"description,omitempty"`
	Gateway     *bool   `json:"gateway,omitempty"`
	ThingModel  *string `json:"thingModel,omitempty"`
	Customer    *string `json:"customer,omitempty"`
}

func getDevices(ctx *gin.Context) {
	var devices map[string]interface{}
	var err error

	// todo: implement query parameters

	if middleware.GetAccessTokenClaims(ctx).CustomerId == "" {
		devices, err = middleware.GetThingsboardAPI(ctx).GetTenantDevices(middleware.GetAccessToken(ctx))
	} else {
		devices, err = middleware.GetThingsboardAPI(ctx).GetCustomerDevices(
			middleware.GetAccessToken(ctx),
			middleware.GetAccessTokenClaims(ctx).CustomerId,
		)
	}

	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, devices)
}

func getDevice(ctx *gin.Context) {
	deviceId := ctx.Param("id")

	device, err := middleware.GetThingsboardAPI(ctx).GetDevice(middleware.GetAccessToken(ctx), deviceId)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	deviceAttributes, err := middleware.GetThingsboardAPI(ctx).GetDeviceAttributes(middleware.GetAccessToken(ctx), deviceId)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	deviceCredentials, err := middleware.GetThingsboardAPI(ctx).GetDeviceCredentials(middleware.GetAccessToken(ctx), deviceId)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	device["attributes"] = deviceAttributes
	device["credentials"] = deviceCredentials

	ctx.JSON(http.StatusOK, device)
}

func createDevice(ctx *gin.Context) {
	createDevice := CreateDevice{}
	if err := ctx.BindJSON(&createDevice); err != nil {
		RespondWithBadRequest(ctx, "Bad Request")
		return
	}

	thingsboardCreateDevice := api.ThingsboardDevice{
		Name: createDevice.Name,
		AdditionalInfo: map[string]interface{}{
			"gateway":     createDevice.Gateway,
			"description": createDevice.Description,
			"customerId": map[string]string{
				"id":         createDevice.CustomerId,
				"entityType": "CUSTOMER",
			},
		},
	}

	createdDevice, err := middleware.GetThingsboardAPI(ctx).CreateDevice(middleware.GetAccessToken(ctx), thingsboardCreateDevice)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	attributes := make(map[string]interface{})
	attributes["thing-model"] = createDevice.ThingModel

	deviceId := createdDevice["id"].(map[string]interface{})["id"].(string)

	if err = middleware.GetThingsboardAPI(ctx).CreateDeviceAttributes(middleware.GetAccessToken(ctx), deviceId, attributes); err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusCreated, createdDevice)
}

func deleteDevice(ctx *gin.Context) {
	deviceId := ctx.Param("id")

	if err := middleware.GetThingsboardAPI(ctx).DeleteDevice(middleware.GetAccessToken(ctx), deviceId); err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": deviceId,
	})
}

func updateDevice(ctx *gin.Context) {
	deviceId := ctx.Param("id")

	updateDevice := UpdateDevice{}
	if err := ctx.BindJSON(&updateDevice); err != nil {
		RespondWithBadRequest(ctx, "Bad Request")
		return
	}

	device, err := middleware.GetThingsboardAPI(ctx).GetDevice(middleware.GetAccessToken(ctx), deviceId)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	if updateDevice.Description != nil {
		device["additionalInfo"].(map[string]interface{})["description"] = updateDevice.Description
	}

	if updateDevice.Gateway != nil {
		device["additionalInfo"].(map[string]interface{})["gateway"] = updateDevice.Gateway
	}

	if updateDevice.Customer != nil {
		device["customerId"] = map[string]string{
			"id":         *updateDevice.Customer,
			"entityType": "CUSTOMER",
		}
	}

	if err := middleware.GetThingsboardAPI(ctx).UpdateDevice(middleware.GetAccessToken(ctx), deviceId, device); err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	if updateDevice.ThingModel != nil {
		attributes := make(map[string]interface{})
		attributes["thing-model"] = updateDevice.ThingModel

		if err := middleware.GetThingsboardAPI(ctx).CreateDeviceAttributes(middleware.GetAccessToken(ctx), deviceId, attributes); err != nil {
			RespondWithInternalServerError(ctx)
			return
		}
	} else {
		if err := middleware.GetThingsboardAPI(ctx).DeleteDeviceAttribute(middleware.GetAccessToken(ctx), deviceId, "thing-model"); err != nil {
			RespondWithInternalServerError(ctx)
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"id": deviceId})
}

func addDevicesRoutes(r *gin.RouterGroup) {
	devicesGroup := r.Group("/devices")
	devicesGroup.GET("/", getDevices)
	devicesGroup.GET("/:id", getDevice)
	devicesGroup.POST("/", createDevice)
	devicesGroup.DELETE("/:id", deleteDevice)
	devicesGroup.PUT("/:id", updateDevice)
}
