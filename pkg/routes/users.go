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
	Group     string `json:"group,omitempty"`
}

type User struct {
	Id            string  `json:"id,omitempty"`
	Email         string  `json:"email"`
	EmailVerified bool    `json:"emailVerified,omitempty"`
	FirstName     string  `json:"firstName"`
	LastName      string  `json:"lastName"`
	Groups        []Group `json:"groups,omitempty"`
	Password      string  `json:"password,omitempty"`
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

	userGroups := []Group{}
	for _, group := range groups {
		// query doesn't seem to work, so we need to filter the groups manually
		tenantId := (*group.Attributes)["tenant-id"]
		customerId := (*group.Attributes)["customer-id"]
		if tenantId != nil && tenantId[0] == middleware.GetAccessTokenClaims(ctx).TenantId && customerId != nil {
			userGroups = append(userGroups, Group{
				Id:   *group.ID,
				Name: *group.Name,
			})
		}
	}

	user := User{
		Id:            *keycloakUser.ID,
		Email:         *keycloakUser.Email,
		EmailVerified: *keycloakUser.EmailVerified,
		Groups:        userGroups,
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
	}

	if user.Group != "" {
		group, err := middleware.GetKeycloakClient(ctx).GetGroup(ctx,
			middleware.GetKeycloakToken(ctx),
			middleware.GetKeycloakRealm(ctx),
			user.Group,
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

			if err := middleware.GetThingsboardAPI(ctx).CreateUser(middleware.GetAccessToken(ctx),
				user.Email,
				user.FirstName,
				user.LastName,
				thingsboardTenantId,
				(*group.Attributes)["customer-id"][0],
			); err != nil {
				RespondWithInternalServerError(ctx)
				return
			}
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

	if err := middleware.GetKeycloakClient(ctx).DeleteUser(ctx,
		middleware.GetKeycloakToken(ctx),
		middleware.GetKeycloakRealm(ctx),
		userId,
	); err != nil {
		RespondWithInternalServerError(ctx)
		return
	}

	if utils.GetConfig().EnableDeviceApi {
		if err := middleware.GetThingsboardAPI(ctx).DeleteUser(middleware.GetKeycloakToken(ctx),
			userId,
		); err != nil {
			RespondWithInternalServerError(ctx)
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func addUsersRoute(r *gin.RouterGroup) {
	usersGroup := r.Group("/users")
	usersGroup.POST("/", createUser)
	usersGroup.GET("/:id", getUser)
	usersGroup.DELETE("/:id", DeleteUser)
}
