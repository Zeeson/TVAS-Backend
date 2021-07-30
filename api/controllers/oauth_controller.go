package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"bitbucket.org/staydigital/truvest-identity-management/api/models"
	"bitbucket.org/staydigital/truvest-identity-management/api/responses"
	"bitbucket.org/staydigital/truvest-identity-management/api/utils/customErrorFormat"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/markbates/goth/gothic"
)

func (server *Server) OauthSignIn(w http.ResponseWriter, r *http.Request) {
	
	if gothUser, err := gothic.CompleteUserAuth(w, r); err == nil {
		user := models.User{}
		fetchUser, _ := user.FindUserByEmail(server.DB, gothUser.Email)
		if len(fetchUser.Email) > 0 {
			token, err := server.SignIn(user.Email, user.Password, gothUser.AccessToken, "other")
			if err != nil {
				formattedError := customErrorFormat.FormatError(err.Error())
				responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
				return
			}
			responses.JSON(w, http.StatusOK, token)			
		} else {
			fetchUser.Email = gothUser.Email
			if len(gothUser.FirstName) > 0 {
				fetchUser.FirstName = gothUser.FirstName
			} else {
				fetchUser.FirstName = strings.Fields(gothUser.Name)[0]
			}
			if len(gothUser.LastName) > 0 {
				fetchUser.LastName = gothUser.LastName
			} else {
				fetchUser.LastName = strings.Fields(gothUser.Name)[1]
			}
			fetchUser.UserName = gothUser.Email
			userID := uuid.New()
			fetchUser.ID = userID
			fetchUser.CreatedBy = userID
			fetchUser.CreatedAt = time.Now()
			fetchUser.UpdatedBy = userID
			fetchUser.UpdatedAt = time.Now()
			userCreated, err := fetchUser.SaveUser(server.DB)
			if err != nil {
				fmt.Println(err)
				formattedError := customErrorFormat.FormatError(err.Error())		
				responses.ERROR(w, http.StatusInternalServerError, formattedError)
				return
			}
			token, err := server.SignIn(userCreated.Email, userCreated.Password, gothUser.AccessToken, "other")
			if err != nil {
				formattedError := customErrorFormat.FormatError(err.Error())
				responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
				return
			}
			responses.JSON(w, http.StatusOK, token)
		}
	} else {
		gothic.BeginAuthHandler(w, r)
	}
}

func (server *Server) OauthSuccessCallback(w http.ResponseWriter, r *http.Request) {	
	
	if gothUser, err := gothic.CompleteUserAuth(w, r); err == nil {
		user := models.User{}
		fetchUser, _ := user.FindUserByEmail(server.DB, gothUser.Email)
		if len(fetchUser.Email) > 0 {
			token, err := server.SignIn(user.Email, user.Password, gothUser.AccessToken, "other")
			if err != nil {
				formattedError := customErrorFormat.FormatError(err.Error())
				responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
				return
			}
			responses.JSON(w, http.StatusOK, token)			
		} else {
			vars := mux.Vars(r)
    		provider, ok := vars["provider"]
    		if !ok {
        		fmt.Println("provider is missing in parameters")
				responses.ERROR(w, http.StatusInternalServerError, errors.New("provider is missing in parameters"))
				return
    		}
			fetchUser.Email = gothUser.Email
			if len(gothUser.FirstName) > 0 {
				fetchUser.FirstName = gothUser.FirstName
			} else {
				fetchUser.FirstName = strings.Fields(gothUser.Name)[0]
			}
			if len(gothUser.LastName) > 0 {
				fetchUser.LastName = gothUser.LastName
			} else {
				fetchUser.LastName = strings.Fields(gothUser.Name)[1]
			}
			fetchUser.UserName = gothUser.Email
			userID := uuid.New()
			fetchUser.ID = userID
			fetchUser.CreatedBy = userID
			fetchUser.CreatedAt = time.Now()
			fetchUser.UpdatedBy = userID
			fetchUser.UpdatedAt = time.Now()
			fetchUser.Enabled = true
			fetchUser.Provider = provider
			userCreated, err := fetchUser.SaveUser(server.DB)
			if err != nil {
				fmt.Println(err)
				formattedError := customErrorFormat.FormatError(err.Error())		
				responses.ERROR(w, http.StatusInternalServerError, formattedError)
				return
			}
			token, err := server.SignIn(userCreated.Email, userCreated.Password, gothUser.AccessToken, provider)
			if err != nil {
				formattedError := customErrorFormat.FormatError(err.Error())
				responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
				return
			}
			responses.JSON(w, http.StatusOK, token)
		}
	}
}

// OAuthLogout godoc
// @Summary Logout from the system using OAuth provider
// @Description User can logout from the system and OAuth provider using this API. User need to pass the provider as param and last device_id with which user logged in.
// @Tags OAuth
// @Accept  json
// @Produce  json
// @Param provider path string true "OAuth provider"
// @Param logout body Logout true "Logout"
// @Success 200 {object} string
// @Security ApiKeyAuth
// @Router /logout/{provider} [post]
func (server *Server) OauthLogout(w http.ResponseWriter, r *http.Request) {	
	
	gothic.Logout(w, r)
	server.Logout(w, r)
}