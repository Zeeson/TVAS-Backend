package api

import (
	"fmt"
	"log"
	"os"

	"bitbucket.org/staydigital/truvest-identity-management/api/controllers"
	"bitbucket.org/staydigital/truvest-identity-management/api/seed"
	"github.com/joho/godotenv"
)

var server = controllers.Server{}

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("sad .env file found")
	}
}

func Run() {

	var err error
	err = godotenv.Load()
	if err != nil {
		log.Printf("Error getting env, %v. But lets move on with the System defined variables", err)
	} else {
		fmt.Println("We are getting the env values")
	}

	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("APP_PROTOCOL"), os.Getenv("APP_HOST"), os.Getenv("APP_PORT"))

	seed.Load(server.DB, os.Getenv("SYSTEM_ADMIN_USERNAME"), os.Getenv("SYSTEM_ADMIN_FIRSTNAME"), os.Getenv("SYSTEM_ADMIN_LASTNAME"), os.Getenv("SYSTEM_EMAIL"), os.Getenv("SYSTEM_ADMIN_PASSWORD"))

	server.Run(":"+os.Getenv("APP_PORT"), []string{os.Getenv("ALLOWED_ORIGINS")})

}