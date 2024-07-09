package routes

import (
	"errors"
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

func getCustomers(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		RespondWithForbidden(ctx)
		return
	}

	customerQuery := QueryParams{}
	if err := ctx.BindQuery(&customerQuery); err != nil {
		RespondWithBadRequest(ctx, "Bad Request")
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
		RespondWithInternalServerError(ctx)
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
		RespondWithForbidden(ctx)
		return
	}

	customerID := ctx.Param("id")

	group, err := middleware.GetKeycloakClient(ctx).GetGroup(ctx,
		middleware.GetKeycloakToken(ctx),
		middleware.GetKeycloakRealm(ctx),
		customerID,
	)
	if err != nil {
		RespondWithResourceNotFound(ctx, customerID)
		return
	}

	description := ""
	if group.Attributes != nil {
		tenantId, ok := (*group.Attributes)["tenant-id"]
		if !ok || tenantId[0] != middleware.GetAccessTokenClaims(ctx).TenantId {
			RespondWithResourceNotFound(ctx, customerID)
			return
		}

		if desc, ok := (*group.Attributes)["description"]; ok {
			description = desc[0]
		}
	} else {
		RespondWithResourceNotFound(ctx, customerID)
		return
	}

	members, err := middleware.GetKeycloakClient(ctx).GetGroupMembers(ctx,
		middleware.GetKeycloakToken(ctx),
		middleware.GetKeycloakRealm(ctx),
		*group.ID,
		gocloak.GetGroupsParams{},
	)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	customer := Customer{
		ID:          *group.ID,
		Name:        *group.Name,
		Description: description,
		Members:     members,
	}

	if group.Attributes == nil || (*group.Attributes)["customer-id"] == nil {
		RespondWithResourceNotFound(ctx, customerID)
		return
	}

	customerId := (*group.Attributes)["customer-id"][0]

	if utils.GetConfig().EnableThingsboard {
		thingsboardCustomer, err := middleware.GetThingsboardAPI(ctx).GetCustomer(
			middleware.GetAccessToken(ctx),
			customerId,
		)
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
		fusekiDataset, err := middleware.GetFusekiAPI(ctx).GetDataset(
			middleware.GetAccessTokenClaims(ctx).TenantId + "-" + customerId,
		)
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
		RespondWithForbidden(ctx)
		return
	}

	customer := Customer{}
	if err := ctx.BindJSON(&customer); err != nil {
		RespondWithBadRequest(ctx, "Bad Request")
		return
	}

	var customerId string

	if utils.GetConfig().EnableThingsboard {
		createdThingsboardCustomer, err := middleware.GetThingsboardAPI(ctx).CreateCustomer(
			middleware.GetAccessToken(ctx),
			api.ThingsboardCustomer{
				Name:  customer.Name,
				Title: customer.Name,
			},
		)
		if err != nil {
			RespondWithInternalServerError(ctx)
			return
		}

		customerId = createdThingsboardCustomer.Id.ID
	} else {
		customerId = uuid.New().String()
	}

	if utils.GetConfig().EnableFuseki {
		if err := middleware.GetFusekiAPI(ctx).CreateDataset(
			middleware.GetAccessTokenClaims(ctx).TenantId + "-" + customerId,
		); err != nil {
			RespondWithInternalServerError(ctx)
			return
		}
	}

	id, err := middleware.GetKeycloakClient(ctx).CreateGroup(ctx,
		middleware.GetKeycloakToken(ctx),
		middleware.GetKeycloakRealm(ctx),
		gocloak.Group{
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
		},
	)

	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	role, err := middleware.GetKeycloakClient(ctx).GetRealmRole(ctx,
		middleware.GetKeycloakToken(ctx),
		middleware.GetKeycloakRealm(ctx),
		"customer",
	)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	if err := middleware.GetKeycloakClient(ctx).AddRealmRoleToGroup(ctx,
		middleware.GetKeycloakToken(ctx),
		middleware.GetKeycloakRealm(ctx),
		id,
		[]gocloak.Role{
			*role,
		},
	); err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

func deleteCustomer(ctx *gin.Context) {
	customerID := ctx.Param("id")

	group, err := middleware.GetKeycloakClient(ctx).GetGroup(ctx,
		middleware.GetKeycloakToken(ctx),
		middleware.GetKeycloakRealm(ctx),
		customerID,
	)
	if err != nil {
		var apiErr *gocloak.APIError
		if errors.As(err, &apiErr) {
			if apiErr.Code == http.StatusNotFound {
				RespondWithResourceNotFound(ctx, customerID)
				return
			}
		}
		RespondWithInternalServerError(ctx)
		return
	}

	if group.Attributes == nil || (*group.Attributes)["customer-id"] == nil {
		RespondWithResourceNotFound(ctx, customerID)
		return
	}

	customerId := (*group.Attributes)["customer-id"][0]

	if utils.GetConfig().EnableThingsboard {
		if err := middleware.GetThingsboardAPI(ctx).DeleteCustomer(middleware.GetAccessToken(ctx), customerId); err != nil {
			RespondWithInternalServerError(ctx)
			return
		}
	}

	if utils.GetConfig().EnableFuseki {
		if err := middleware.GetFusekiAPI(ctx).DeleteDataset(middleware.GetAccessTokenClaims(ctx).TenantId + "-" + customerId); err != nil {
			RespondWithInternalServerError(ctx)
			return
		}
	}

	if err := middleware.GetKeycloakClient(ctx).DeleteGroup(
		ctx,
		middleware.GetKeycloakToken(ctx),
		middleware.GetKeycloakRealm(ctx),
		customerID,
	); err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": customerID,
	})
}

func updateCustomer(ctx *gin.Context) {
	customerID := ctx.Param("id")

	updateCustomer := CustomerUpdate{}
	if err := ctx.BindJSON(&updateCustomer); err != nil {
		RespondWithBadRequest(ctx, "Bad Request")
		return
	}

	group, err := middleware.GetKeycloakClient(ctx).GetGroup(
		ctx,
		middleware.GetKeycloakToken(ctx),
		middleware.GetKeycloakRealm(ctx),
		customerID,
	)
	if err != nil {
		var apiErr *gocloak.APIError
		if errors.As(err, &apiErr) {
			if apiErr.Code == http.StatusNotFound {
				RespondWithResourceNotFound(ctx, customerID)
				return
			}
		}
		RespondWithInternalServerError(ctx)
		return
	}

	if group.Attributes == nil || (*group.Attributes)["customer-id"] == nil {
		RespondWithResourceNotFound(ctx, customerID)
		return
	}

	if (*group.Attributes)["description"] == nil {
		(*group.Attributes)["description"] = []string{updateCustomer.Description}
	} else {
		(*group.Attributes)["description"][0] = updateCustomer.Description
	}

	if err := middleware.GetKeycloakClient(ctx).UpdateGroup(
		ctx,
		middleware.GetKeycloakToken(ctx),
		middleware.GetKeycloakRealm(ctx),
		*group,
	); err != nil {
		RespondWithInternalServerError(ctx)
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
