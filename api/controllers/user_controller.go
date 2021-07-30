package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"bitbucket.org/staydigital/truvest-identity-management/api/auth"
	"bitbucket.org/staydigital/truvest-identity-management/api/models"
	"bitbucket.org/staydigital/truvest-identity-management/api/responses"
	"bitbucket.org/staydigital/truvest-identity-management/api/utils"
	"bitbucket.org/staydigital/truvest-identity-management/api/utils/customErrorFormat"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// CreateUser godoc
// @Summary Create a user in the system
// @Description Create a user by ID in the system. In order to access this API, someone must have "USERS_CREATE" Permission tagged to its role. If Password is not passed while creating User, then it creates User with auto-generated password and sends it through Email.
// @Tags User
// @Accept  json
// @Produce  json
// @Param user body models.CreateUser true "Create User"
// @Success 201 {object} models.UserResponse
// @Security ApiKeyAuth
// @Router /users [post]
func (server *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	err := auth.CheckBlacklistedJWT(server.TTLCache, r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	check, err := server.HasPermission(r, []string{"USERS_CREATE"})
	if !check || err != nil {
		responses.ERROR(w, http.StatusForbidden, errors.New("Forbidden"))
		return
	}
	
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	var password string
	if user.Password == "" {
		password = utils.RandSeq()
		user.Password = password
	}
	fmt.Printf("New Auto-generated User Password: %s", password)
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
	user.Prepare(tokenID)
	err = user.Validate("")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	userCreated, err := user.SaveUser(server.DB)

	if err != nil {
		fmt.Println(err)
		formattedError := customErrorFormat.FormatError(err.Error())

		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	sm := models.SendMail{}
	sm.Email = userCreated.Email
	err = sm.SendEmail(password)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, userCreated.ID))
	responses.JSON(w, http.StatusCreated, models.PrepareResponse(userCreated))
}

// GetAllUser godoc
// @Summary Get all users in the system
// @Description Get details of all users in the system. In order to access this API, someone must have "SYSTEM_ADMIN" and "USERS_VIEW" Permission tagged to its role.
// @Tags Users
// @Accept  json
// @Produce  json
// @Success 200 {object} models.User
// @Security ApiKeyAuth
// @Router /users [get]
func (server *Server) GetUsers(w http.ResponseWriter, r *http.Request) {
	err := auth.CheckBlacklistedJWT(server.TTLCache, r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	check, err := server.HasPermission(r, []string{"SYSTEM_ADMIN", "USERS_VIEW"})
	if !check || err != nil {
		responses.ERROR(w, http.StatusForbidden, errors.New("Forbidden"))
		return
	}

	user := models.User{}

	users, err := user.FindAllUsers(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, models.PrepareResponses(users))
}

// GetUser godoc
// @Summary Get a user by ID in the system
// @Description Get details of a user by ID in the system. In order to access this API, someone must have "USERS_VIEW" Permission tagged to its role.
// @Tags User
// @Accept  json
// @Produce  json
// @Param id path string true "ID of the user to get"
// @Success 200 {object} models.UserResponse
// @Security ApiKeyAuth
// @Router /users/{id} [get]
func (server *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	err := auth.CheckBlacklistedJWT(server.TTLCache, r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	check, err := server.HasPermission(r, []string{"USERS_VIEW"})
	if !check || err != nil {
		responses.ERROR(w, http.StatusForbidden, errors.New("Forbidden"))
		return
	}

	vars := mux.Vars(r)
	uid, err := uuid.Parse(vars["id"])
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	user := models.User{}
	userGotten, err := user.FindUserByID(server.DB, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, models.PrepareResponse(userGotten))
}

// UpdateUser godoc
// @Summary Update a user in the system
// @Description Update a user by ID in the system. In order to access this API, someone must have "USERS_CREATE" Permission tagged to its role.
// @Tags User
// @Accept  json
// @Produce  json
// @Param id path int true "ID of the user to be updated"
// @Param user body models.CreateUser true "Update User"
// @Success 200 {object} models.UserResponse
// @Security ApiKeyAuth
// @Router /users/{id} [put]
func (server *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	err := auth.CheckBlacklistedJWT(server.TTLCache, r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	check, err := server.HasPermission(r, []string{"USERS_CREATE"})
	if !check || err != nil {
		responses.ERROR(w, http.StatusForbidden, errors.New("Forbidden"))
		return
	}

	vars := mux.Vars(r)
	uid, err := uuid.Parse(vars["id"])
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	fetchUser := models.User{}
	updateUser, err := fetchUser.FindUserByID(server.DB, uid)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = json.Unmarshal(body, &updateUser)
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
	updatedUser, err := updateUser.UpdateAUser(server.DB, uid, tokenID)
	if err != nil {
		formattedError := customErrorFormat.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, models.PrepareResponse(updatedUser))
}

// DeleteUser godoc
// @Summary Delete a user in the system
// @Description Delete a user by ID in the system. In order to access this API, someone must have "USERS_CREATE" Permission tagged to its role.
// @Tags User
// @Accept  json
// @Produce  json
// @Param id path int true "ID of the user to be deleted"
// @Success 204
// @Security ApiKeyAuth
// @Router /users/{id} [delete]
func (server *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	err := auth.CheckBlacklistedJWT(server.TTLCache, r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	check, err := server.HasPermission(r, []string{"USERS_CREATE"})
	if !check || err != nil {
		responses.ERROR(w, http.StatusForbidden, errors.New("Forbidden"))
		return
	}

	vars := mux.Vars(r)

	user := models.User{}

	uid, err := uuid.Parse(vars["id"])
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	_, err = user.DeleteAUser(server.DB, uid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	responses.JSON(w, http.StatusNoContent, "")
}

// GetLoggedInUser godoc
// @Summary Get the Logged in User details
// @Description Get the logged in user details by Authorization header being passed. This is useful by other microservices to check for authorization based upon the roles and permissions listed.
// @Tags User
// @Accept  json
// @Produce  json
// @Success 204 {object} models.UserResponse
// @Security ApiKeyAuth
// @Router /user/me [get]
func (server *Server) GetLoggedInUser(w http.ResponseWriter, r *http.Request) {
	err := auth.CheckBlacklistedJWT(server.TTLCache, r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
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
	user := models.User{}
	userGotten, err := user.WhoAmI(server.DB, tokenID)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, models.PrepareResponse(userGotten))
}

// SetPassword godoc
// @Summary Set the Password of the user.
// @Description Update a user's password by ID in the system. The User himself can only change the password through this endpoint.
// @Tags User
// @Accept  json
// @Produce  json
// @Param id path string true "ID of the user to be updated"
// @Param user body models.Set_User_Password_Payload true "Update User"
// @Success 200 
// @Security ApiKeyAuth
// @Router /users/{id}/setPassword [post]
func (server *Server) SetPassword(w http.ResponseWriter, r *http.Request) {
	err := auth.CheckBlacklistedJWT(server.TTLCache, r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	vars := mux.Vars(r)

	uid, err := uuid.Parse(vars["id"])
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
	if tokenID != uid {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	setPassword := models.Set_User_Password_Payload{}
	err = json.Unmarshal(body, &setPassword)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = setPassword.ResetPassword(server.DB, uid, tokenID)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	responses.JSON(w, http.StatusNoContent, "")
}

// ForgotPassword godoc
// @Summary Set the Password of the user if he has forgotten
// @Description Set the Password of the user if he has forgotten.
// @Tags User
// @Accept  json
// @Produce  json
// @Param user body models.Forgot_User_Password_Payload true "Forgot Password"
// @Success 200 
// @Router /users/forgotPassword [post]
func (server *Server) ForgotPassword(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	forgotPassword := models.Forgot_User_Password_Payload{}
	err = json.Unmarshal(body, &forgotPassword)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = forgotPassword.ForgetPassword(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusNoContent, "")
}

// SendMail godoc
// @Summary Send a mail with a link for the user to reset password
// @Description Send a mail with a link for the user to reset password. This can be used in combination with Forget Password API to make sure a legitimate user is trying to reset his password.
// @Tags User
// @Accept  json
// @Produce  json
// @Param user body models.SendMail true "Send Email"
// @Success 200 
// @Router /users/sendMail [post]
func (server *Server) SendMail(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	sendMail := models.SendMail{}
	err = json.Unmarshal(body, &sendMail)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	// err = sm.SendGridMail()
	err = sendMail.SendEmail("")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	responses.JSON(w, http.StatusNoContent, "")
}

// EnableUser godoc
// @Summary Enable/Disable a user in the system
// @Description Enable/Disable a user by ID in the system. In order to access this API, someone must have "USERS_CREATE" Permission tagged to its role. Set Enabled field in payload to true to enable a user/Set Enabled field in payload to false to disable a user 
// @Tags User
// @Accept  json
// @Produce  json
// @Param id path int true "ID of the user to be updated"
// @Param user body models.EnableUser true "Enable User"
// @Success 200 {object} models.UserResponse
// @Security ApiKeyAuth
// @Router /users/{id}/enableUser [put]
func (server *Server) EnableUser(w http.ResponseWriter, r *http.Request) {
	err := auth.CheckBlacklistedJWT(server.TTLCache, r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	check, err := server.HasPermission(r, []string{"USERS_CREATE"})
	if !check || err != nil {
		responses.ERROR(w, http.StatusForbidden, errors.New("Forbidden"))
		return
	}

	vars := mux.Vars(r)
	uid, err := uuid.Parse(vars["id"])
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	fetchUser := models.User{}
	updateUser, err := fetchUser.FindUserByID(server.DB, uid)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = json.Unmarshal(body, &updateUser)
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
	updatedUser, err := updateUser.EnableDisableUser(server.DB, uid, tokenID)
	if err != nil {
		formattedError := customErrorFormat.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, models.PrepareResponse(updatedUser))
}
