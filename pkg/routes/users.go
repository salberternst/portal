package routes

import (
	"fmt"
	"net/http"

	"github.com/Nerzal/gocloak/v13"
	"github.com/gin-gonic/gin"
	"github.com/salberternst/portal/pkg/middleware"
	"github.com/salberternst/portal/pkg/utils"
)

type RealmRoleMapping struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Composite   bool   `json:"composite"`
	ClientRole  bool   `json:"clientRole"`
	ContainerId string `json:"containerId"`
}

type Group struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type CreateUser struct {
	Id        string `json:"id,omitempty"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	IsAdmin   bool   `json:"isAdmin,omitempty"`
	Customer  string `json:"customer,omitempty"`
}

type User struct {
	Id            string     `json:"id,omitempty"`
	Email         string     `json:"email"`
	EmailVerified bool       `json:"emailVerified,omitempty"`
	FirstName     string     `json:"firstName"`
	LastName      string     `json:"lastName"`
	Customers     []Customer `json:"customers,omitempty"`
	Password      string     `json:"password,omitempty"`
}

func getUser(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		RespondWithForbidden(ctx)
		return
	}

	userId := ctx.Param("id")

	keycloakUser, err := middleware.GetKeycloakClient(ctx).GetUserByID(ctx,
		middleware.GetKeycloakToken(ctx),
		middleware.GetKeycloakRealm(ctx),
		userId,
	)
	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	Q := fmt.Sprintf("tenant-id:%s customer-id:", middleware.GetAccessTokenClaims(ctx).TenantId)
	briefRepresentation := false
	groups, err := middleware.GetKeycloakClient(ctx).GetUserGroups(
		ctx,
		middleware.GetKeycloakToken(ctx),
		middleware.GetKeycloakRealm(ctx),
		userId,
		gocloak.GetGroupsParams{
			BriefRepresentation: &briefRepresentation,
			Q:                   &Q,
		},
	)

	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	customers := []Customer{}
	for _, group := range groups {
		if group.Attributes == nil || (*group.Attributes)["tenant-id"] == nil || (*group.Attributes)["customer-id"] == nil {
			continue
		}

		tenantId := (*group.Attributes)["tenant-id"]
		customerId := (*group.Attributes)["customer-id"]

		if tenantId != nil && tenantId[0] == middleware.GetAccessTokenClaims(ctx).TenantId && customerId != nil {
			customers = append(customers, Customer{
				ID:   *group.ID,
				Name: *group.Name,
			})
		}
	}

	user := User{
		Id:            *keycloakUser.ID,
		Email:         *keycloakUser.Email,
		EmailVerified: *keycloakUser.EmailVerified,
		Customers:     customers,
	}

	if keycloakUser.FirstName != nil {
		user.FirstName = *keycloakUser.FirstName
	}

	if keycloakUser.LastName != nil {
		user.LastName = *keycloakUser.LastName
	}

	ctx.JSON(http.StatusOK, user)
}

func createUser(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		RespondWithForbidden(ctx)
		return
	}

	user := CreateUser{}
	if err := ctx.BindJSON(&user); err != nil {
		RespondWithBadRequest(ctx, "Bad Request")
		return
	}

	keycloakUserPassword := utils.GenerateRandomPassword(16)

	keycloakUser := gocloak.User{
		Username:      &user.Email,
		Email:         &user.Email,
		FirstName:     &user.FirstName,
		LastName:      &user.LastName,
		Enabled:       gocloak.BoolP(true),
		EmailVerified: gocloak.BoolP(true),
		Credentials: &[]gocloak.CredentialRepresentation{
			{
				Temporary: gocloak.BoolP(true),
				Type:      gocloak.StringP("password"),
				Value:     &keycloakUserPassword,
			},
		},
		Attributes: &map[string][]string{},
	}

	if user.Customer != "" {
		group, err := middleware.GetKeycloakClient(ctx).GetGroup(ctx,
			middleware.GetKeycloakToken(ctx),
			middleware.GetKeycloakRealm(ctx),
			user.Customer,
		)

		if err != nil {
			RespondWithInternalServerError(ctx)
			return
		}

		if (*group.Attributes)["tenant-id"][0] != middleware.GetAccessTokenClaims(ctx).TenantId {
			RespondWithForbidden(ctx)
			return
		}

		keycloakUser.Groups = &[]string{*group.Path}

		if utils.GetConfig().EnableDeviceApi {
			thingsboardTenantId, err := middleware.GetThingsboardAPI(ctx).GetTenantId(middleware.GetAccessToken(ctx))
			if err != nil {
				RespondWithInternalServerError(ctx)
				return
			}

			if (*group.Attributes)["thingsboard-customer-id"] == nil {
				RespondWithInternalServerError(ctx)
				return
			}

			thingsboardUserId, err := middleware.GetThingsboardAPI(ctx).CreateUser(middleware.GetAccessToken(ctx),
				user.Email,
				user.FirstName,
				user.LastName,
				thingsboardTenantId,
				(*group.Attributes)["thingsboard-customer-id"][0],
			)

			if err != nil {
				RespondWithInternalServerError(ctx)
				return
			}

			(*keycloakUser.Attributes)["thingsboard-user-id"] = []string{thingsboardUserId}
		}
	} else if user.IsAdmin {
		keycloakUser.Groups = &[]string{middleware.GetAccessTokenClaims(ctx).TenantId}
	} else {
		RespondWithBadRequest(ctx, "Bad Request")
		return
	}

	userId, err := middleware.GetKeycloakClient(ctx).CreateUser(ctx,
		middleware.GetKeycloakToken(ctx),
		middleware.GetKeycloakRealm(ctx),
		keycloakUser,
	)

	if err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":       userId,
		"password": keycloakUserPassword,
	})
}

func DeleteUser(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		RespondWithForbidden(ctx)
		return
	}

	userId := ctx.Param("id")

	if utils.GetConfig().EnableDeviceApi {
		users, err := middleware.GetKeycloakClient(ctx).GetUserByID(ctx,
			middleware.GetKeycloakToken(ctx),
			middleware.GetKeycloakRealm(ctx),
			userId,
		)
		if err != nil {
			RespondWithInternalServerError(ctx)
			return
		}

		if users.Attributes == nil || (*users.Attributes)["thingsboard-user-id"] == nil {
			RespondWithInternalServerError(ctx)
			return
		}

		thingsboardUserId := (*users.Attributes)["thingsboard-user-id"][0]

		if err := middleware.GetThingsboardAPI(ctx).DeleteUser(middleware.GetAccessToken(ctx),
			thingsboardUserId,
		); err != nil {
			RespondWithInternalServerError(ctx)
			return
		}
	}

	if err := middleware.GetKeycloakClient(ctx).DeleteUser(ctx,
		middleware.GetKeycloakToken(ctx),
		middleware.GetKeycloakRealm(ctx),
		userId,
	); err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": userId,
	})
}

func addUsersRoute(r *gin.RouterGroup) {
	usersGroup := r.Group("/users")
	usersGroup.POST("/", createUser)
	usersGroup.GET("/:id", getUser)
	usersGroup.DELETE("/:id", DeleteUser)
}
