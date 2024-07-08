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

type UserInput struct {
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
		ctx.JSON(http.StatusForbidden, gin.H{
			"error":   "forbidden",
			"message": "not allowed to access this resource",
			"status":  http.StatusForbidden,
		})
		return
	}

	userId := ctx.Param("id")

	keycloakUser, err := middleware.GetKeycloakClient(ctx).GetUserByID(ctx,
		middleware.GetKeycloakToken(ctx),
		middleware.GetKeycloakRealm(ctx),
		userId,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"error":   "unable_to_get_user",
			"message": fmt.Sprintf("Failed to get user: %s", err.Error()),
		})
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
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"error":   "unable_to_get_user_groups",
			"message": fmt.Sprintf("Failed to get user groups: %s", err.Error()),
		})
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
		ctx.JSON(http.StatusForbidden, gin.H{
			"error":   "forbidden",
			"message": "not allowed to access this resource",
			"status":  http.StatusForbidden,
		})
		return
	}

	user := UserInput{}
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"error":   "bad_request",
			"message": fmt.Sprintf("Failed to bind JSON: %s", err.Error()),
		})
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
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "internal_server_error",
				"error":   err.Error(),
				"message": "Failed to get group",
			})
			return
		}

		if (*group.Attributes)["tenant-id"][0] != middleware.GetAccessTokenClaims(ctx).TenantId {
			ctx.JSON(http.StatusForbidden, gin.H{
				"status":  http.StatusForbidden,
				"error":   "forbidden",
				"message": "You are not allowed to create users for this customer",
			})
			return
		}

		keycloakUser.Groups = &[]string{*group.Path}
	} else if user.IsAdmin {
		keycloakUser.Groups = &[]string{middleware.GetAccessTokenClaims(ctx).TenantId}
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"error":   "bad_request",
			"message": "Customer or admin flag is required",
		})
		return
	}

	userId, err := middleware.GetKeycloakClient(ctx).CreateUser(ctx,
		middleware.GetKeycloakToken(ctx),
		middleware.GetKeycloakRealm(ctx),
		keycloakUser,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "internal_server_error",
			"error":   err.Error(),
			"message": "Failed to create user",
		})
		return
	}

	if utils.GetConfig().EnableThingsboard {
		if err := middleware.GetThingsboardAPI(ctx).CreateUser(middleware.GetKeycloakToken(ctx),
			user.Email,
			user.FirstName,
			user.LastName,
			middleware.GetAccessTokenClaims(ctx).TenantId,
			user.Group,
		); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "internal_server_error",
				"error":   err.Error(),
				"message": "Failed to create user in Thingsboard",
			})
			return
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":       userId,
		"password": keycloakUserPassword,
	})
}

func DeleteUser(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error":   "forbidden",
			"message": "not allowed to access this resource",
			"status":  http.StatusForbidden,
		})
		return
	}

	userId := ctx.Param("id")

	if err := middleware.GetKeycloakClient(ctx).DeleteUser(ctx,
		middleware.GetKeycloakToken(ctx),
		middleware.GetKeycloakRealm(ctx),
		userId,
	); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "internal_server_error",
			"error":   err.Error(),
			"message": "Failed to delete user",
		})
		return
	}

	if utils.GetConfig().EnableThingsboard {
		if err := middleware.GetThingsboardAPI(ctx).DeleteUser(middleware.GetKeycloakToken(ctx),
			userId,
		); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "internal_server_error",
				"error":   err.Error(),
				"message": "Failed to delete user in Thingsboard",
			})
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
