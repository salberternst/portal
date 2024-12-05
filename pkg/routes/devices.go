package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/salberternst/portal/pkg/api"
	"github.com/salberternst/portal/pkg/middleware"
)

type SyncStatus struct {
	Message   string  `json:"message"`
	Status    string  `json:"status"`
	Timestamp float64 `json:"ts"`
}

type DeviceCredentials struct {
	Id          string `json:"id"`
	Credentials string `json:"credentials"`
	Type        string `json:"type"`
}

type Device struct {
	Id           string             `json:"id"`
	Name         string             `json:"name,omitempty"`
	Gateway      bool               `json:"gateway,omitempty"`
	CustomerId   string             `json:"customerId,omitempty"`
	SyncStatus   *SyncStatus        `json:"syncStatus,omitempty"`
	CreatedAt    int64              `json:"createdAt,omitempty"`
	Credentials  *DeviceCredentials `json:"credentials"`
	ThingModel   string             `json:"thingModel,omitempty"`
	Title        string             `json:"title,omitempty"`
	Description  string             `json:"description,omitempty"`
	Category     string             `json:"category,omitempty"`
	Manufacturer string             `json:"manufacturer,omitempty"`
	Model        string             `json:"mode,omitempty"`
}

type CreateDevice struct {
	Name         string `json:"name"`
	Gateway      bool   `json:"gateway,omitempty"`
	CustomerId   string `json:"customerId,omitempty"`
	ThingModel   string `json:"thingModel,omitempty"`
	Title        string `json:"title,omitempty"`
	Description  string `json:"description,omitempty"`
	Category     string `json:"category,omitempty"`
	Manufacturer string `json:"manufacturer,omitempty"`
	Model        string `json:"model,omitempty"`
}

type UpdateDevice struct {
	Gateway      *bool   `json:"gateway,omitempty"`
	ThingModel   *string `json:"thingModel,omitempty"`
	CustomerId   *string `json:"customerId,omitempty"`
	Title        *string `json:"title,omitempty"`
	Description  *string `json:"description,omitempty"`
	Category     *string `json:"category,omitempty"`
	Manufacturer *string `json:"manufacturer,omitempty"`
	Model        *string `json:"model,omitempty"`
}

func getDevices(ctx *gin.Context) {
	queryParams := QueryParams{}
	if err := ctx.BindQuery(&queryParams); err != nil {
		RespondWithBadRequest(ctx, "Bad Request")
		return
	}

	var thingsboardDevices map[string]interface{}
	var err error

	// todo: implement query parameters

	if middleware.GetAuthenticatedUser(ctx).CustomerId == "" {
		thingsboardDevices, err = middleware.GetThingsboardAPI(ctx).GetTenantDevices(middleware.GetAccessToken(ctx), queryParams.Page, queryParams.PageSize)
	} else {
		thingsboardCustomerId, err := middleware.GetThingsboardAPI(ctx).GetCustomerId(middleware.GetAccessToken(ctx))
		if err != nil {
			RespondWithInternalServerError(ctx)
			return
		}
		thingsboardDevices, err = middleware.GetThingsboardAPI(ctx).GetCustomerDevices(
			middleware.GetAccessToken(ctx),
			thingsboardCustomerId,
			queryParams.Page,
			queryParams.PageSize,
		)
		if err != nil {
			RespondWithInternalServerError(ctx)
			return
		}
	}

	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	devices := []Device{}
	for _, item := range thingsboardDevices["data"].([]interface{}) {
		thingsboardDevice := item.(map[string]interface{})
		device := Device{
			Id:        thingsboardDevice["id"].(map[string]interface{})["id"].(string),
			Name:      thingsboardDevice["name"].(string),
			CreatedAt: int64(thingsboardDevice["createdTime"].(float64)),
		}

		additionalInfo := thingsboardDevice["additionalInfo"]
		if additionalInfo != nil {
			gateway := additionalInfo.(map[string]interface{})["gateway"]
			if gateway != nil {
				device.Gateway = gateway.(bool)
			}
		}

		thingsboardCustomerId := thingsboardDevice["customerId"].(map[string]interface{})["id"].(string)
		if thingsboardCustomerId != "" && thingsboardCustomerId != "13814000-1dd2-11b2-8080-808080808080" {
			device.CustomerId, _ = middleware.GetCustomerIdByThingsboardCustomerId(ctx, thingsboardCustomerId)
			// If the customer is not found, we will just not set the customer id (for now)
		}

		devices = append(devices, device)

	}

	thingsboardDevices["data"] = devices

	ctx.JSON(http.StatusOK, thingsboardDevices)
}

func getDevice(ctx *gin.Context) {
	deviceId := ctx.Param("id")

	thingsboardDevice, err := middleware.GetThingsboardAPI(ctx).GetDevice(middleware.GetAccessToken(ctx), deviceId)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	thingsboardDeviceAttributes, err := middleware.GetThingsboardAPI(ctx).GetDeviceAttributes(middleware.GetAccessToken(ctx), deviceId)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	thingsboardDeviceCredentials, err := middleware.GetThingsboardAPI(ctx).GetDeviceCredentials(middleware.GetAccessToken(ctx), deviceId)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	device := &Device{
		Id:        thingsboardDevice["id"].(map[string]interface{})["id"].(string),
		Name:      thingsboardDevice["name"].(string),
		CreatedAt: int64(thingsboardDevice["createdTime"].(float64)),
		Credentials: &DeviceCredentials{
			Credentials: thingsboardDeviceCredentials["credentialsId"].(string),
			Type:        thingsboardDeviceCredentials["credentialsType"].(string),
		},
	}

	thingsboardCustomerId := thingsboardDevice["customerId"].(map[string]interface{})["id"].(string)
	if thingsboardCustomerId != "" && thingsboardCustomerId != "13814000-1dd2-11b2-8080-808080808080" {
		device.CustomerId, _ = middleware.GetCustomerIdByThingsboardCustomerId(ctx, thingsboardCustomerId)
		// If the customer is not found, we will just not set the customer id (for now)
	}

	additionalInfo := thingsboardDevice["additionalInfo"].(map[string]interface{})
	if additionalInfo != nil {
		if additionalInfo["gateway"] != nil {
			device.Gateway = additionalInfo["gateway"].(bool)
		}
	}

	for _, deviceAttribute := range thingsboardDeviceAttributes {
		if deviceAttribute["key"] == "thing-metadata" {
			thingMetadata := deviceAttribute["value"].(map[string]interface{})
			if thingMetadata != nil {
				if thingMetadata["title"] != nil {
					device.Title = thingMetadata["title"].(string)
				}
				if thingMetadata["description"] != nil {
					device.Description = thingMetadata["description"].(string)
				}
				if thingMetadata["category"] != nil {
					device.Category = thingMetadata["category"].(string)
				}
				if thingMetadata["manufacturer"] != nil {
					device.Manufacturer = thingMetadata["manufacturer"].(string)
				}
				if thingMetadata["model"] != nil {
					device.Model = thingMetadata["model"].(string)
				}
			}
		} else if deviceAttribute["key"] == "thing-model" {
			if deviceAttribute["value"] != nil {
				device.ThingModel = deviceAttribute["value"].(string)
			}
		} else if deviceAttribute["key"] == "thing-registry-sync-status" {
			if deviceAttribute["value"] != nil {
				thingRegistrySyncStatus := deviceAttribute["value"].(map[string]interface{})
				device.SyncStatus = &SyncStatus{
					Message:   thingRegistrySyncStatus["message"].(string),
					Status:    thingRegistrySyncStatus["status"].(string),
					Timestamp: thingRegistrySyncStatus["ts"].(float64),
				}
			}
		}
	}

	ctx.JSON(http.StatusOK, device)
}

func createDevice(ctx *gin.Context) {
	createDevice := CreateDevice{}
	if err := ctx.BindJSON(&createDevice); err != nil {
		RespondWithBadRequest(ctx, "Bad Request")
		return
	}

	if createDevice.CustomerId == "" {
		createDevice.CustomerId = "13814000-1dd2-11b2-8080-808080808080"
	} else {
		thingsboardCustomerId, err := middleware.GetThingsboardCustomerIdByCustomerId(ctx, createDevice.CustomerId)
		if err != nil {
			RespondWithInternalServerError(ctx)
			return
		}
		createDevice.CustomerId = thingsboardCustomerId
	}

	thingsboardCreateDevice := api.ThingsboardDevice{
		Name: createDevice.Name,
		Customer: api.ThingsboardId{
			Id:         createDevice.CustomerId,
			EntityType: "CUSTOMER",
		},
		AdditionalInfo: map[string]interface{}{
			"gateway":     createDevice.Gateway,
			"description": createDevice.Description,
		},
	}

	createdDevice, err := middleware.GetThingsboardAPI(ctx).CreateDevice(middleware.GetAccessToken(ctx), thingsboardCreateDevice)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	attributes := make(map[string]interface{})
	attributes["thing-model"] = createDevice.ThingModel
	attributes["thing-metadata"] = map[string]interface{}{
		"title":        createDevice.Title,
		"category":     createDevice.Category,
		"manufacturer": createDevice.Manufacturer,
		"model":        createDevice.Model,
		"description":  createDevice.Description,
	}

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

	if updateDevice.CustomerId == nil {
		customerId := "13814000-1dd2-11b2-8080-808080808080"
		updateDevice.CustomerId = &customerId
	} else {
		thingsboardCustomerId, err := middleware.GetThingsboardCustomerIdByCustomerId(ctx, *updateDevice.CustomerId)
		if err != nil {
			RespondWithInternalServerError(ctx)
			return
		}
		updateDevice.CustomerId = &thingsboardCustomerId
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

	if updateDevice.CustomerId != nil {
		device["customerId"] = map[string]string{
			"id":         *updateDevice.CustomerId,
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
		attributes["thing-metadata"] = map[string]interface{}{
			"title":        updateDevice.Title,
			"category":     updateDevice.Category,
			"manufacturer": updateDevice.Manufacturer,
			"model":        updateDevice.Model,
			"description":  updateDevice.Description,
		}

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
