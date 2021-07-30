package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"bitbucket.org/staydigital/truvest-identity-management/api/auth"
	"bitbucket.org/staydigital/truvest-identity-management/api/models"
	"bitbucket.org/staydigital/truvest-identity-management/api/responses"
	"bitbucket.org/staydigital/truvest-identity-management/api/utils"
	"bitbucket.org/staydigital/truvest-identity-management/api/utils/customErrorFormat"
	"github.com/ReneKroon/ttlcache/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Login struct {
	Email     	string    	`gorm:"size:100;not null;unique" json:"email"`
	Password  	string    	`gorm:"size:100;not null;" json:"password,omitempty"`
	DeviceID 	string 		`gorm:"size:255;not null;" json:"device_id"`
}

type Refresh struct {
	RefreshToken   	string    	`gorm:"size:100;not null" json:"refresh_token"`
}

type Logout struct {
	DeviceID 	string 		`gorm:"size:255;not null;" json:"device_id"`
}

// SignUp godoc
// @Summary Signup as a user in the system
// @Description Signup as a user in the system. If Password is not passed while creating User, then it creates User with auto-generated password. The Password is then sent through Email.
// @Tags SignUp
// @Accept  json
// @Produce  json
// @Param user body models.CreateUser true "Create User"
// @Success 201 {object} models.UserResponse
// @Router /signup [post]
func (server *Server) SignUp(w http.ResponseWriter, r *http.Request) {
	
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
	user.PrepareSignUp()
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

// Login godoc
// @Summary Login to the system
// @Description Login to the system. It returns the accessToken, refreshToken, expiry and device_id in a JSON format. You need to pass email and password along with a unique device_id. You can use device fingerprint or browser fingerprint to generate a unique device_id. The access and refresh token will be mapped to this id so that whenever the user logs out, it can detect the device from which the user is logging out in case of multi-session login and can expire the JWT for that specific device.
// @Tags Login
// @Accept  json
// @Produce  json
// @Param login body Login true "Login"
// @Success 200 {object} models.LoginResponse
// @Router /login [post]
func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	login := Login{}
	err = json.Unmarshal(body, &login)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	user.Email = login.Email
	user.Password = login.Password

	err = user.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	token, err := server.SignIn(user.Email, user.Password, login.DeviceID, "local")
	if err != nil {
		formattedError := customErrorFormat.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, token)
}

func (server *Server) SignIn(email, password, deviceID string, provider string) (models.LoginResponse, error) {

	var err error
	var role string

	user := models.User{}

	err = server.DB.Debug().Model(models.User{}).Where("email = ?", email).Preload("Roles").Take(&user).Error
	if err != nil {
		return models.LoginResponse{}, err
	}
	if !user.Enabled {
		return models.LoginResponse{}, errors.New("Account is locked!!")
	}
	if strings.EqualFold(provider, "local") {
		err = models.VerifyPassword(user.Password, password)
		if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
			return models.LoginResponse{}, err
		}
	}
	if len(user.Roles) > 0 {
		role = user.Roles[0].Name
	} else {
		role = ""
	}

	userDevice := models.User_Device{}
	userDevice.ID = uuid.New()
	userDevice.DeviceID = deviceID
	userDevice.IsRefreshActive = true
	userDevice.UserID = user.ID
	userDevice.CreatedAt = time.Now()
	userDevice.CreatedBy = user.ID
	userDevice.UpdatedAt = time.Now()
	userDevice.UpdatedBy = user.ID
	err = userDevice.SaveUserDevice(server.DB)
	if err != nil {
		return models.LoginResponse{}, err
	}

	loginResponse := models.LoginResponse{}
	refreshToken := models.Refresh_Token{}
	loginResponse, refreshToken, err = auth.CreateToken(user.ID, user.UserName, user.Email, role, deviceID)
	if err != nil {
		return models.LoginResponse{}, err
	}
	
	err = refreshToken.SaveRefreshToken(server.DB)
	if err != nil {
		return models.LoginResponse{}, err
	}

	return loginResponse, nil
}

// Refresh godoc
// @Summary Refresh JWT Token to Login to the system
// @Description Refresh JWT token to Login to the system. It returns the accessToken, refreshToken and expiry in a JSON format. You need to pass the refresh token from the previous call.
// @Tags Login
// @Accept  json
// @Produce  json
// @Param refresh body Refresh true "Refresh"
// @Success 200 {object} models.LoginResponse
// @Security ApiKeyAuth
// @Router /refresh [post]
func (server *Server) Refresh(w http.ResponseWriter, r *http.Request) {
	err := auth.CheckBlacklistedJWT(server.TTLCache, r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	refresh := Refresh{}
	err = json.Unmarshal(body, &refresh)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	var refreshToken uuid.UUID
	refreshToken, err = uuid.Parse(refresh.RefreshToken)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	var uid string
	uid, err = auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	var userID uuid.UUID
	userID, err = uuid.Parse(uid)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	token, err := server.Recreate(userID, refreshToken)
	if err != nil {
		formattedError := customErrorFormat.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}

	accessToken := auth.ExtractToken(r)
	expiry := auth.ExtractTokenExpiry(r)
	ttl := expiry - time.Now().Unix()
	server.TTLCache = ttlcache.NewCache()
	server.TTLCache.SetWithTTL(userID.String(), accessToken, time.Duration(ttl)*time.Second)
	fmt.Println(server.TTLCache.Get(userID.String()))

	responses.JSON(w, http.StatusOK, token)
}

func (server *Server) Recreate(userID uuid.UUID, token uuid.UUID) (models.LoginResponse, error) {

	var err error
	var role string

	user := models.User{}

	err = server.DB.Debug().Model(models.User{}).Where("id = ?", userID).Preload("Roles").Take(&user).Error
	if err != nil {
		return models.LoginResponse{}, err
	}
	if !user.Enabled {
		return models.LoginResponse{}, errors.New("Account is locked!!")
	}
	if len(user.Roles) > 0 {
		role = user.Roles[0].Name
	} else {
		role = ""
	}

	userDevice := models.User_Device{}
	device := &models.User_Device{}
	device, err = userDevice.FindUserDeviceByUserID(server.DB, userID)
	if err != nil {
		return models.LoginResponse{}, err
	}

	loginResponse := models.LoginResponse{}
	loginResponse, _, err = auth.CreateToken(user.ID, user.UserName, user.Email, role, device.DeviceID)
	if err != nil {
		return models.LoginResponse{}, err
	}

	refreshToken := models.Refresh_Token{}
	err = refreshToken.UpdateARefreshToken(server.DB, token, user.ID)
	if err != nil {
		return models.LoginResponse{}, err
	}

	return loginResponse, nil
}

// Logout godoc
// @Summary Logout from the system
// @Description User can logout from the system using this API. User need to pass device_id so that he can be logged out of the specific device during multi-session login. This will expire the JWT for the device from where he has logged in.
// @Tags Login
// @Accept  json
// @Produce  json
// @Param logout body Logout true "Logout"
// @Success 200 {object} string
// @Security ApiKeyAuth
// @Router /logout [post]
func (server *Server) Logout(w http.ResponseWriter, r *http.Request) {
	err := auth.CheckBlacklistedJWT(server.TTLCache, r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	logout := Logout{}
	err = json.Unmarshal(body, &logout)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	var token string
	token = auth.ExtractToken(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	var uid string
	uid, err = auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	var expiry int64
	expiry = auth.ExtractTokenExpiry(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	var userID uuid.UUID
	userID, err = uuid.Parse(uid)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = server.Destroy(userID, logout.DeviceID, token, expiry)
	if err != nil {
		formattedError := customErrorFormat.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, "User Logged out successfully!")
}

func (server *Server) Destroy(userID uuid.UUID, deviceID string, token string, expiry int64) (error) {

	var err error

	user := models.User{}
	err = server.DB.Debug().Model(models.User{}).Where("id = ?", userID).Preload("Roles").Take(&user).Error
	if err != nil {
		return err
	}
	if !user.Enabled {
		return errors.New("Account is locked!!")
	}

	userDevice := models.User_Device{}
	err = userDevice.DeleteAUserDevice(server.DB, userID)
	if err != nil {
		return err
	}

	ttl := expiry - time.Now().Unix()
	server.TTLCache = ttlcache.NewCache()
	server.TTLCache.SetWithTTL(user.ID.String(), token, time.Duration(ttl)*time.Second)
	fmt.Println(server.TTLCache.Get(user.ID.String()))

	refreshToken := models.Refresh_Token{}
	err = refreshToken.DeleteARefreshToken(server.DB, deviceID)
	if err != nil {
		return err
	}

	return nil
}