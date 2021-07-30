package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"bitbucket.org/staydigital/truvest-identity-management/api/auth"
	"bitbucket.org/staydigital/truvest-identity-management/api/models"
	"bitbucket.org/staydigital/truvest-identity-management/api/responses"
	"bitbucket.org/staydigital/truvest-identity-management/api/utils/customErrorFormat"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// CreateRole godoc
// @Summary Create a role in the system
// @Description Create a role in the system. User must have "MANAGE_ROLES" permission tagged to its role in order to use this API.
// @Tags Role
// @Accept  json
// @Produce  json
// @Param user body models.CreateRole true "Create Role"
// @Success 201 {object} models.Role
// @Security ApiKeyAuth
// @Router /roles [post]
func (server *Server) CreateRole(w http.ResponseWriter, r *http.Request) {
	err := auth.CheckBlacklistedJWT(server.TTLCache, r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	check, err := server.HasPermission(r, []string{"MANAGE_ROLES"})
	if !check || err != nil {
		responses.ERROR(w, http.StatusForbidden, errors.New("Forbidden"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	role := models.Role{}
	err = json.Unmarshal(body, &role)
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
	role.Prepare(tokenID)
	err = role.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	roleCreated, err := role.SaveRole(server.DB)

	if err != nil {

		formattedError := customErrorFormat.FormatError(err.Error())

		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, roleCreated.ID))
	responses.JSON(w, http.StatusCreated, roleCreated)
}

// GetRoles godoc
// @Summary Get all Roles in the system
// @Description Get all roles in the system. In order to access this API, someone must have "ROLES_VIEW" Permission tagged to its role.
// @Tags Roles
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Role
// @Security ApiKeyAuth
// @Router /roles [get]
func (server *Server) GetRoles(w http.ResponseWriter, r *http.Request) {
	err := auth.CheckBlacklistedJWT(server.TTLCache, r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	check, err := server.HasPermission(r, []string{"ROLES_VIEW"})
	if !check || err != nil {
		responses.ERROR(w, http.StatusForbidden, errors.New("Forbidden"))
		return
	}

	role := models.Role{}

	roles, err := role.FindAllRoles(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, roles)
}

// GetRole godoc
// @Summary Get a Role by id in the system
// @Description Get a role by id in the system. In order to access this API, someone must have "ROLES_VIEW" Permission tagged to its role.
// @Tags Role
// @Accept  json
// @Produce  json
// @Param id path int true "ID of the role to get"
// @Success 200 {object} models.Role
// @Security ApiKeyAuth
// @Router /roles/{id} [get]
func (server *Server) GetRole(w http.ResponseWriter, r *http.Request) {
	err := auth.CheckBlacklistedJWT(server.TTLCache, r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	check, err := server.HasPermission(r, []string{"ROLES_VIEW"})
	if !check || err != nil {
		responses.ERROR(w, http.StatusForbidden, errors.New("Forbidden"))
		return
	}

	vars := mux.Vars(r)
	rid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	role := models.Role{}
	roleGotten, err := role.FindRoleByID(server.DB, uint32(rid))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, roleGotten)
}

// UpdateRole godoc
// @Summary Update a role by id in the system
// @Description Update a role by id in the system. In order to access this API, someone must have "MANAGE_ROLES" Permission tagged to its role.
// @Tags Role
// @Accept  json
// @Produce  json
// @Param id path int true "id of the Role"
// @Param user body models.CreateRole true "Create Role"
// @Success 200 {object} models.Role
// @Security ApiKeyAuth
// @Router /roles/{id} [put]
func (server *Server) UpdateRole(w http.ResponseWriter, r *http.Request) {
	err := auth.CheckBlacklistedJWT(server.TTLCache, r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	check, err := server.HasPermission(r, []string{"MANAGE_ROLES"})
	if !check || err != nil {
		responses.ERROR(w, http.StatusForbidden, errors.New("Forbidden"))
		return
	}

	vars := mux.Vars(r)
	rid, err := strconv.ParseUint(vars["id"], 10, 32)
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
	fetchRole := models.Role{}
	updateRole, err := fetchRole.FindRoleByID(server.DB, uint32(rid))
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = json.Unmarshal(body, &updateRole)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = updateRole.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	updatedRole, err := updateRole.UpdateARole(server.DB, uint32(rid), tokenID)
	if err != nil {
		formattedError := customErrorFormat.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, updatedRole)
}

// DeleteRole godoc
// @Summary Delete a role by id in the system
// @Description Delete a role by id in the system. In order to access this API, someone must have "MANAGE_ROLES" Permission tagged to its role.
// @Tags Role
// @Accept  json
// @Produce  json
// @Param id path int true "id of the Role"
// @Success 204 
// @Security ApiKeyAuth
// @Router /roles/{id} [delete]
func (server *Server) DeleteRole(w http.ResponseWriter, r *http.Request) {
	err := auth.CheckBlacklistedJWT(server.TTLCache, r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	check, err := server.HasPermission(r, []string{"MANAGE_ROLES"})
	if !check || err != nil {
		responses.ERROR(w, http.StatusForbidden, errors.New("Forbidden"))
		return
	}

	vars := mux.Vars(r)

	role := models.Role{}

	rid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	_, err = role.DeleteARole(server.DB, uint32(rid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", rid))
	responses.JSON(w, http.StatusNoContent, "")
}

// AddUsersToRole godoc
// @Summary Add users to role by id in the system
// @Description Add users to role by id in the system. In order to access this API, someone must have "USERS_ASSIGN_TO_ROLE" Permission tagged to its role.
// @Tags Role
// @Accept  json
// @Produce  json
// @Param id path int true "id of the Role"
// @Param permission body models.User_Role_Payload true "User Role Payload"
// @Success 200 
// @Security ApiKeyAuth
// @Router /roles/{id}/users [post]
func (server *Server) AddUsersToRole(w http.ResponseWriter, r *http.Request) {
	err := auth.CheckBlacklistedJWT(server.TTLCache, r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	check, err := server.HasPermission(r, []string{"USERS_ASSIGN_TO_ROLE"})
	if !check || err != nil {
		responses.ERROR(w, http.StatusForbidden, errors.New("Forbidden"))
		return
	}

	vars := mux.Vars(r)
	rid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	uids := models.User_Role_Payload{}
	err = json.Unmarshal(body, &uids)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	for i := range uids.Users {
		ur := models.User_Role{}
		ur.UserID = uids.Users[i]
		ur.RoleID = uint32(rid)

		ur.Prepare()
		err = ur.Validate()
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}
		err := ur.SaveUserToRole(server.DB)
		if err != nil {

			formattedError := customErrorFormat.FormatError(err.Error())
	
			responses.ERROR(w, http.StatusInternalServerError, formattedError)
			return
		}
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/", r.Host, r.RequestURI))
	responses.JSON(w, http.StatusCreated, "")
}

// DeleteUsersFromRole godoc
// @Summary Delete users from role by id in the system
// @Description Delete users from role by id in the system. In order to access this API, someone must have
// "USERS_ASSIGN_TO_ROLE" Permission tagged to its role.
// @Tags Role
// @Accept  json
// @Produce  json
// @Param id1 path int true "id of the Role"
// @Param id2 path int true "id of the User"
// @Success 200 
// @Security ApiKeyAuth
// @Router /roles/{id1}/users/{id2} [delete]
func (server *Server) DeleteUsersFromRole(w http.ResponseWriter, r *http.Request) {
	err := auth.CheckBlacklistedJWT(server.TTLCache, r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	check, err := server.HasPermission(r, []string{"USERS_ASSIGN_TO_ROLE"})
	if !check || err != nil {
		responses.ERROR(w, http.StatusForbidden, errors.New("Forbidden"))
		return
	}

	vars := mux.Vars(r)
	roleUsers := models.User_Role{}
	rid, err := strconv.ParseUint(vars["id1"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	uid, err := uuid.Parse(vars["id2"])
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	_, err = roleUsers.DeleteUsersFromRole(server.DB, uint32(rid), uid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d,%d", rid, uid))
	responses.JSON(w, http.StatusNoContent, "")
}

// AddPermissionsToRole godoc
// @Summary Add permissions to role by id in the system
// @Description Add permissions to role by id in the system. In order to access this API, someone must have "MANAGE_PERMISSION" and "PERMISSION_ASSIGN_TO_ROLE" Permission tagged to its role.
// @Tags Role
// @Accept  json
// @Produce  json
// @Param id path int true "id of the Role"
// @Param permission body models.Role_Permission_Payload true "Role Permission Payload"
// @Success 200 
// @Security ApiKeyAuth
// @Router /roles/{id}/permissions [post]
func (server *Server) AddPermissionsToRole(w http.ResponseWriter, r *http.Request) {
	err := auth.CheckBlacklistedJWT(server.TTLCache, r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	check, err := server.HasPermission(r, []string{"MANAGE_PERMISSION", "PERMISSION_ASSIGN_TO_ROLE"})
	if !check || err != nil {
		responses.ERROR(w, http.StatusForbidden, errors.New("Forbidden"))
		return
	}

	vars := mux.Vars(r)
	rid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	pids := models.Role_Permission_Payload{}
	err = json.Unmarshal(body, &pids)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	for i := range pids.Permissions {
		rp := models.Role_Permission{}
		rp.PermissionID = uint32(pids.Permissions[i])
		rp.RoleID = uint32(rid)

		rp.Prepare()
		err = rp.Validate()
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}
		err := rp.SavePermissionToRole(server.DB)
		if err != nil {

			formattedError := customErrorFormat.FormatError(err.Error())
	
			responses.ERROR(w, http.StatusInternalServerError, formattedError)
			return
		}
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/", r.Host, r.RequestURI))
	responses.JSON(w, http.StatusCreated, "")
}

// DeletePermissionsFromRole godoc
// @Summary Delete permissions from role by id in the system
// @Description Delete permissions from role by id in the system.
// @Tags Role
// @Accept  json
// @Produce  json
// @Param id1 path int true "id of the Role"
// @Param id2 path int true "id of the Permission"
// @Success 200 
// @Security ApiKeyAuth
// @Router /roles/{id1}/permissions/{id2} [delete]
func (server *Server) DeletePermissionsFromRole(w http.ResponseWriter, r *http.Request) {
	err := auth.CheckBlacklistedJWT(server.TTLCache, r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	check, err := server.HasPermission(r, []string{"MANAGE_PERMISSION", "PERMISSION_ASSIGN_TO_ROLE"})
	if !check || err != nil {
		responses.ERROR(w, http.StatusForbidden, errors.New("Forbidden"))
		return
	}

	vars := mux.Vars(r)
	rolePermission := models.Role_Permission{}
	rid, err := strconv.ParseUint(vars["id1"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	pid, err := strconv.ParseUint(vars["id2"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	_, err = rolePermission.DeleteRoleFromPermission(server.DB, uint32(rid), uint32(pid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d,%d", rid, pid))
	responses.JSON(w, http.StatusNoContent, "")
}
