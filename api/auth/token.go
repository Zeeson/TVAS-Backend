package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/staydigital/truvest-identity-management/api/models"
	"github.com/ReneKroon/ttlcache/v2"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

func CreateToken(user_id uuid.UUID, username string, email string, role string, deviceID string) (models.LoginResponse, models.Refresh_Token, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user_id.String()
	claims["username"] = username
	claims["email"] = email
	claims["role"] = role
	tokenExpiry, _ := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRY_IN_MILLISECOND")) 
	expiry := time.Now().Add(time.Millisecond * time.Duration(tokenExpiry)).Unix() //Token expires after defined time interval
	claims["exp"] = expiry
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	var accessToken string
	var err error
	accessToken, err = token.SignedString([]byte(os.Getenv("API_SECRET")))
	if err != nil {
		return models.LoginResponse{}, models.Refresh_Token{}, err
	}

	refreshToken := models.Refresh_Token{}
	refreshToken.ID = uuid.New()
	refreshToken.DeviceID = deviceID
	refreshToken.RefreshCount = int64(0)
	refreshToken.Token = uuid.New()
	refreshToken.ExpiryDate = expiry
	refreshToken.CreatedAt = time.Now()
	refreshToken.CreatedBy = user_id
	refreshToken.UpdatedAt = time.Now()
	refreshToken.UpdatedBy = user_id

	loginResponse := models.LoginResponse{}	
	loginResponse.AccessToken = accessToken
	loginResponse.RefreshToken = refreshToken.Token
	loginResponse.ExpiryDate = expiry
	loginResponse.DeviceID = deviceID
	return loginResponse, refreshToken, nil
}

func TokenValid(r *http.Request) error {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		Pretty(claims)
	}
	return nil
}

func ExtractToken(r *http.Request) string {
	keys := r.URL.Query()
	token := keys.Get("token")
	if token != "" {
		return token
	}
	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func ExtractTokenID(r *http.Request) (string, error) {

	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid := fmt.Sprint(claims["user_id"])
		return uid, nil
	}
	return "", nil
}

func ExtractTokenExpiry(r *http.Request) (int64) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return int64(0)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		expiry := claims["exp"]
		return int64(expiry.(float64))
	}
	return int64(0)
}

func CheckBlacklistedJWT(cache *ttlcache.Cache, r *http.Request) (error) {
	tokenString := ExtractToken(r)
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid := fmt.Sprint(claims["user_id"])
		jwtCache, _ := cache.Get(uid)
		if jwtCache != nil {
			if strings.EqualFold(jwtCache.(string), tokenString) {
				return fmt.Errorf("User either already logged out of the device or accesstoken inactive: %v", uid)
			}
		}
	}
	return nil
}

//Pretty display the claims nicely in the terminal
func Pretty(data interface{}) {
	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(string(b))
}