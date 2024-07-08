package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/salberternst/portal/pkg/api"
	"github.com/salberternst/portal/pkg/middleware"
	"github.com/salberternst/portal/pkg/utils"
)

var TenantType = "tenant"
var CustomerType = "customer"

type FusekiDataset struct {
	Name  string `json:"name"`
	State bool   `json:"state,omitempty"`
	Error string `json:"error,omitempty"`
}

type ThingsboardCustomer struct {
	Id       string `json:"id,omitempty"`
	Title    string `json:"title,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Country  string `json:"country,omitempty"`
	State    string `json:"state,omitempty"`
	City     string `json:"city,omitempty"`
	Address  string `json:"address,omitempty"`
	Address2 string `json:"address2,omitempty"`
	Zip      string `json:"zip,omitempty"`
	Error    string `json:"error,omitempty"`
}

type Customer struct {
	ID          string               `json:"id,omitempty"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Thingsboard *ThingsboardCustomer `json:"thingsboard,omitempty"`
	Fuseki      *FusekiDataset       `json:"fuseki,omitempty"`
	Members     []*gocloak.User      `json:"members,omitempty"`
}

type CustomerUpdate struct {
	Description string `json:"description"`
}

type CustomerQuery struct {
	Page     uint `form:"page,default=1" binding:"required"`
	PageSize uint `form:"page_size,default=20" binding:"required"`
}

func getCustomers(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error":   "forbidden",
			"message": "not allowed to access this resource",
			"status":  http.StatusForbidden,
		})
		return
	}

	customerQuery := CustomerQuery{}
	if err := ctx.BindQuery(&customerQuery); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}

	// only show groups that have the tenant-id attribute set to the tenant-id of the user
	Q := fmt.Sprintf("tenant-id:%s type:customer", middleware.GetAccessTokenClaims(ctx).TenantId)
	briefRepresentation := false

	groups, err := middleware.GetKeycloakClient(ctx).GetGroups(ctx,
		middleware.GetKeycloakToken(ctx),
		middleware.GetKeycloakRealm(ctx),
		gocloak.GetGroupsParams{
			Q:                   &Q,
			BriefRepresentation: &briefRepresentation,
		},
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"error":   "unable_to_get_customers",
			"message": err.Error(),
		})
		return
	}

	customers := make([]Customer, len(groups))
	for i, group := range groups {
		description := ""
		customerId := ""

		if group.Attributes != nil {
			if desc, ok := (*group.Attributes)["description"]; ok {
				description = desc[0]
			}
			if id, ok := (*group.Attributes)["customer-id"]; ok {
				customerId = id[0]
			}
		}

		customers[i] = Customer{
			ID:          *group.ID,
			Name:        *group.Name,
			Description: description,
		}

		if utils.GetConfig().EnableThingsboard {
			customers[i].Thingsboard = &ThingsboardCustomer{
				Id: customerId,
			}
		}

		if utils.GetConfig().EnableFuseki {
			customers[i].Fuseki = &FusekiDataset{
				Name: middleware.GetAccessTokenClaims(ctx).TenantId + "-" + customerId,
			}
		}
	}

	ctx.JSON(http.StatusOK, customers)
}

func getCustomer(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error":   "forbidden",
			"message": "not allowed to access this resource",
			"status":  http.StatusForbidden,
		})
		return
	}

	customerID := ctx.Param("id")

	accessToken := middleware.GetAccessToken(ctx)
	thingsboardApi := middleware.GetThingsboardAPI(ctx)
	fusekiApi := middleware.GetFusekiAPI(ctx)

	group, err := middleware.GetKeycloakClient(ctx).GetGroup(ctx,
		middleware.GetKeycloakToken(ctx),
		middleware.GetKeycloakRealm(ctx),
		customerID,
	)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "customer_not_found",
			"message": fmt.Sprintf("customer with id %s not found", customerID),
			"status":  http.StatusNotFound,
		})
		return
	}

	description := ""
	if group.Attributes != nil {
		tenantId, ok := (*group.Attributes)["tenant-id"]
		if !ok || tenantId[0] != middleware.GetAccessTokenClaims(ctx).TenantId {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "customer_not_found",
				"message": fmt.Sprintf("customer with id %s not found", customerID),
				"status":  http.StatusNotFound,
			})
			return
		}

		if desc, ok := (*group.Attributes)["description"]; ok {
			description = desc[0]
		}
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "customer_not_found",
			"message": fmt.Sprintf("customer with id %s not found", customerID),
			"status":  http.StatusNotFound,
		})
		return
	}

	members, err := middleware.GetKeycloakClient(ctx).GetGroupMembers(ctx,
		middleware.GetKeycloakToken(ctx),
		middleware.GetKeycloakRealm(ctx),
		*group.ID,
		gocloak.GetGroupsParams{},
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"error":   "unable_to_get_group_members",
			"message": err.Error(),
		})
		return
	}

	customer := Customer{
		ID:          *group.ID,
		Name:        *group.Name,
		Description: description,
		Members:     members,
	}

	if group.Attributes == nil || (*group.Attributes)["customer-id"] == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "customer_not_found",
			"message": fmt.Sprintf("customer with id %s not found", customerID),
			"status":  http.StatusNotFound,
		})
		return
	}

	customerId := (*group.Attributes)["customer-id"][0]

	if utils.GetConfig().EnableThingsboard {
		thingsboardCustomer, err := thingsboardApi.GetCustomer(accessToken, customerId)
		if err != nil {
			customer.Thingsboard = &ThingsboardCustomer{
				Error: err.Error(),
			}
		} else {
			customer.Thingsboard = &ThingsboardCustomer{
				Id:       thingsboardCustomer.Id.ID,
				Title:    thingsboardCustomer.Title,
				Name:     thingsboardCustomer.Name,
				Email:    thingsboardCustomer.Email,
				Phone:    thingsboardCustomer.Phone,
				Country:  thingsboardCustomer.Country,
				State:    thingsboardCustomer.State,
				City:     thingsboardCustomer.City,
				Address:  thingsboardCustomer.Address,
				Address2: thingsboardCustomer.Address2,
				Zip:      thingsboardCustomer.Zip,
			}
		}
	}

	if utils.GetConfig().EnableFuseki {
		fusekiDataset, err := fusekiApi.GetDataset(middleware.GetAccessTokenClaims(ctx).TenantId + "-" + customerId)
		if err != nil {
			customer.Fuseki = &FusekiDataset{
				Error: err.Error(),
			}
		} else {
			customer.Fuseki = &FusekiDataset{
				Name:  fusekiDataset.Name,
				State: fusekiDataset.State,
			}
		}
	}

	ctx.JSON(http.StatusOK, customer)
}

func createCustomer(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error":   "forbidden",
			"message": "not allowed to access this resource",
			"status":  http.StatusForbidden,
		})
		return
	}

	customer := Customer{}
	if err := ctx.BindJSON(&customer); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}

	client := ctx.MustGet("keycloak-client").(*gocloak.GoCloak)
	thingsboardApi := ctx.MustGet("thingsboard-api").(*api.ThingsboardAPI)
	fusekiApi := ctx.MustGet("fuseki-api").(*api.FusekiAPI)
	accessToken := ctx.MustGet("access-token").(string)
	keycloakToken := ctx.MustGet("keycloak-token").(string)

	var customerId string

	if utils.GetConfig().EnableThingsboard {
		createdThingsboardCustomer, err := thingsboardApi.CreateCustomer(accessToken, api.ThingsboardCustomer{
			Name:  customer.Name,
			Title: customer.Name,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"error":   "unable_to_create_thingsboard_customer",
				"message": err.Error(),
			})
			return
		}

		customerId = createdThingsboardCustomer.Id.ID
	} else {
		customerId = uuid.New().String()
	}

	if utils.GetConfig().EnableFuseki {
		if err := fusekiApi.CreateDataset(middleware.GetAccessTokenClaims(ctx).TenantId + "-" + customerId); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"error":   "unable_to_create_fuseki_dataset",
				"message": err.Error(),
			})
			return
		}
	}

	id, err := client.CreateGroup(ctx, keycloakToken, middleware.GetKeycloakRealm(ctx), gocloak.Group{
		Name: &customer.Name,
		Attributes: &map[string][]string{
			"tenant-id":             {middleware.GetAccessTokenClaims(ctx).TenantId},
			"customer-id":           {customerId},
			"thingsboard-tenant-id": {middleware.GetAccessTokenClaims(ctx).ThingsboardTenantId},
			"description":           {customer.Description},
			"created-by":            {middleware.GetAccessTokenClaims(ctx).Email},
			"created-at":            {time.Now().String()},
			"type":                  {CustomerType},
		},
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"error":   "unable_to_create_keycloak_group",
			"message": err.Error(),
		})
		return
	}

	role, err := client.GetRealmRole(ctx,
		keycloakToken,
		middleware.GetKeycloakRealm(ctx),
		"customer",
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"error":   "unable_to_get_keycloak_role",
			"message": err.Error(),
		})
		return
	}

	err = client.AddRealmRoleToGroup(ctx,
		keycloakToken,
		middleware.GetKeycloakRealm(ctx),
		id,
		[]gocloak.Role{
			*role,
		},
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"error":   "unable_to_add_keycloak_role_to_group",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

func deleteCustomer(ctx *gin.Context) {
	customerID := ctx.Param("id")

	client := ctx.MustGet("keycloak-client").(*gocloak.GoCloak)
	thingsboardApi := ctx.MustGet("thingsboard-api").(*api.ThingsboardAPI)
	fusekiApi := ctx.MustGet("fuseki-api").(*api.FusekiAPI)
	accessToken := ctx.MustGet("access-token").(string)
	keycloakToken := ctx.MustGet("keycloak-token").(string)

	group, err := client.GetGroup(ctx,
		middleware.GetKeycloakToken(ctx),
		middleware.GetKeycloakRealm(ctx),
		customerID,
	)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "customer not found",
			"status":  http.StatusNotFound,
		})
		return
	}

	if group.Attributes == nil || (*group.Attributes)["customer-id"] == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "customer not found",
			"status":  http.StatusNotFound,
		})
		return
	}

	customerId := (*group.Attributes)["customer-id"][0]

	if utils.GetConfig().EnableThingsboard {
		err = thingsboardApi.DeleteCustomer(accessToken, customerId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"error":   "unable_to_delete_thingsboard_customer",
				"message": err.Error(),
			})
			return
		}
	}

	if utils.GetConfig().EnableFuseki {
		err = fusekiApi.DeleteDataset(middleware.GetAccessTokenClaims(ctx).TenantId + "-" + customerId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"error":   "unable_to_delete_fuseki_dataset",
				"message": err.Error(),
			})
			return
		}
	}

	if err := client.DeleteGroup(ctx,
		keycloakToken,
		middleware.GetKeycloakRealm(ctx),
		customerID,
	); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"error":   "unable_to_delete_keycloak_group",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "customer deleted",
		"status":  http.StatusOK,
	})
}

func updateCustomer(ctx *gin.Context) {
	customerID := ctx.Param("id")

	updateCustomer := CustomerUpdate{}
	if err := ctx.BindJSON(&updateCustomer); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := ctx.MustGet("keycloak-client").(*gocloak.GoCloak)
	keycloakToken := ctx.MustGet("keycloak-token").(string)

	group, err := client.GetGroup(ctx,
		keycloakToken,
		middleware.GetKeycloakRealm(ctx),
		customerID,
	)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"error":   "customer_not_found",
			"message": fmt.Sprintf("customer with id %s not found", customerID),
		})
		return
	}

	if group.Attributes == nil {
		group.Attributes = &map[string][]string{
			"description": {updateCustomer.Description},
		}
	} else if (*group.Attributes)["description"] == nil {
		(*group.Attributes)["description"] = []string{updateCustomer.Description}
	} else {
		(*group.Attributes)["description"][0] = updateCustomer.Description
	}

	err = client.UpdateGroup(ctx, keycloakToken, middleware.GetKeycloakRealm(ctx), *group)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"error":   "unable_to_update_customer",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": customerID,
	})
}

func addCustomersRoutes(r *gin.RouterGroup) {
	customersGroup := r.Group("/customers")
	customersGroup.GET("/", getCustomers)
	customersGroup.GET("/:id", getCustomer)
	customersGroup.POST("/", createCustomer)
	customersGroup.DELETE("/:id", deleteCustomer)
	customersGroup.PUT("/:id", updateCustomer)
}
