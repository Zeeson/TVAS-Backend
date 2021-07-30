package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"bitbucket.org/staydigital/truvest-identity-management/api/models"
	"github.com/ReneKroon/ttlcache/v2"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/rs/cors"
)

type Server struct {
	DB     		*gorm.DB
	Router 		*mux.Router
	TTLCache	*ttlcache.Cache
}

func (server *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName, appProtocol, appHost, appPort string) {

	var err error

	if Dbdriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
		server.DB, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database", Dbdriver)
		}
	}

	server.TTLCache = ttlcache.NewCache()

	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), appProtocol + "://" + appHost + ":" + appPort + "/auth/google/callback"),
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), appProtocol + "://" + appHost + ":" + appPort + "/auth/github/callback"),
	)

	server.DB.Debug().AutoMigrate(&models.User{}, &models.Role{}, &models.User_Role{}, &models.Permission{}, &models.Role_Permission{}, models.User_Device{}, models.Refresh_Token{}) //database migration

	server.Router = mux.NewRouter()

	server.initializeRoutes()
}

func (server *Server) Run(addr string, allowedOrigins []string) {
	c := cors.New(cors.Options{
		AllowedOrigins: allowedOrigins,
	})

	fmt.Println("Listening to port "+addr)
	log.Fatal(http.ListenAndServe(addr, c.Handler(server.Router)))
}