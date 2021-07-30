package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"bitbucket.org/staydigital/truvest-identity-management/api/auth"
	"bitbucket.org/staydigital/truvest-identity-management/api/models"
	"bitbucket.org/staydigital/truvest-identity-management/api/responses"
	"bitbucket.org/staydigital/truvest-identity-management/api/utils"
	"bitbucket.org/staydigital/truvest-identity-management/api/utils/customErrorFormat"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// CreatePermission godoc
// @Summary Create a permission in the system
// @Description Create a permission in the system.User must have "MANAGE_PERMISSION" permission tagged to its role in order to access this API
// @Tags Permission
// @Accept  json
// @Produce  json
// @Param user body models.CreatePermission true "Create Permission"
// @Success 201 {object} models.Permission
// @Security ApiKeyAuth
// @Router /permissions [post]
func (server *Server) CreatePermission(w http.ResponseWriter, r *http.Request) {
	err := auth.CheckBlacklistedJWT(server.TTLCache, r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	check, err := server.HasPermission(r, []string{"MANAGE_PERMISSION"})
	if !check || err != nil {
		responses.ERROR(w, http.StatusForbidden, errors.New("Forbidden"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	permission := models.Permission{}
	err = json.Unmarshal(body, &permission)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	authID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	tokenID, err := uuid.Parse(authID)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	permission.Prepare(tokenID)
	err = permission.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	permissionCreated, err := permission.SavePermission(server.DB)

	if err != nil {

		formattedError := customErrorFormat.FormatError(err.Error())

		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, permissionCreated.ID))
	responses.JSON(w, http.StatusCreated, permissionCreated)
}

// GetPermissions godoc
// @Summary Get all Permissions in the system
// @Description Get all permissions in the system. In order to access this API, someone must have "VIEW_PERMISSION" Permission tagged to its role.
// @Tags Permissions
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Permission
// @Security ApiKeyAuth
// @Router /permissions [get]
func (server *Server) GetPermissions(w http.ResponseWriter, r *http.Request) {
	err := auth.CheckBlacklistedJWT(server.TTLCache, r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	check, err := server.HasPermission(r, []string{"VIEW_PERMISSION"})
	if !check || err != nil {
		responses.ERROR(w, http.StatusForbidden, errors.New("Forbidden"))
		return
	}

	permission := models.Permission{}

	permissions, err := permission.FindAllPermissions(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, permissions)
}

// GetPermission godoc
// @Summary Get a Permission by id in the system
// @Description Get details of a permission by ID in the system. In order to access this API, someone must have "VIEW_PERMISSION" Permission tagged to its role.
// @Tags Permission
// @Accept  json
// @Produce  json
// @Param id path int true "ID of the permission to get"
// @Success 200 {object} models.Permission
// @Security ApiKeyAuth
// @Router /permissions/{id} [get]
func (server *Server) GetPermission(w http.ResponseWriter, r *http.Request) {
	err := auth.CheckBlacklistedJWT(server.TTLCache, r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	check, err := server.HasPermission(r, []string{"VIEW_PERMISSION"})
	if !check || err != nil {
		responses.ERROR(w, http.StatusForbidden, errors.New("Forbidden"))
		return
	}

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	permission := models.Permission{}
	permissionGotten, err := permission.FindPermissionByID(server.DB, uint32(pid))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, permissionGotten)
}

// UpdatePermission godoc
// @Summary Update a permission by id in the system
// @Description Update a permission by id in the system. In order to access this API, someone must have "MANAGE_PERMISSION" Permission tagged to its role.
// @Tags Permission
// @Accept  json
// @Produce  json
// @Param id path int true "id of the Permission"
// @Param user body models.CreatePermission true "Create Permission"
// @Success 200 {object} models.Permission
// @Security ApiKeyAuth
// @Router /permissions/{id} [put]
func (server *Server) UpdatePermission(w http.ResponseWriter, r *http.Request) {
	err := auth.CheckBlacklistedJWT(server.TTLCache, r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	check, err := server.HasPermission(r, []string{"MANAGE_PERMISSION"})
	if !check || err != nil {
		responses.ERROR(w, http.StatusForbidden, errors.New("Forbidden"))
		return
	}

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	authID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	tokenID, err := uuid.Parse(authID)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	fetchPermission := models.Permission{}
	updatePermission, err := fetchPermission.FindPermissionByID(server.DB, uint32(pid))
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = json.Unmarshal(body, &updatePermission)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	updatePermission.Prepare(tokenID)
	err = updatePermission.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	updatedPermission, err := updatePermission.UpdateAPermission(server.DB, uint32(pid), tokenID)
	if err != nil {
		formattedError := customErrorFormat.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, updatedPermission)
}

// DeletePermission godoc
// @Summary Delete a permission by id in the system
// @Description Delete a permission by id in the system. In order to access this API, someone must have "MANAGE_PERMISSION" Permission tagged to its role.
// @Tags Permission
// @Accept  json
// @Produce  json
// @Param id path int true "id of the Permission"
// @Success 204 
// @Security ApiKeyAuth
// @Router /permissions/{id} [delete]
func (server *Server) DeletePermission(w http.ResponseWriter, r *http.Request) {
	err := auth.CheckBlacklistedJWT(server.TTLCache, r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	check, err := server.HasPermission(r, []string{"MANAGE_PERMISSION"})
	if !check || err != nil {
		responses.ERROR(w, http.StatusForbidden, errors.New("Forbidden"))
		return
	}

	vars := mux.Vars(r)

	permission := models.Permission{}

	pid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	_, err = permission.DeleteAPermission(server.DB, uint32(pid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")
}

// Permission Checks
func (server *Server) HasPermission(r *http.Request, p []string) (bool, error) {
	tokenString := auth.ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return false, err
	}
	
	var permissions []string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		uid, err := uuid.Parse(fmt.Sprint(claims["user_id"]))
		if err != nil {
			return false, err
		}
		user := models.User{}

		userGotten, err := user.FindUserByID(server.DB, uid)
		for i := range userGotten.Roles {
			for j := range userGotten.Roles[i].Permissions {
				permissions = append(permissions, userGotten.Roles[i].Permissions[j].Name)
			}
		}
		if err != nil {
			return false, err
		}
	}
	return utils.Contains(p, permissions), nil
}
