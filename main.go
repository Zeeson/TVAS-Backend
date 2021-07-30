package main

import (
	"os"
	"strings"

	"bitbucket.org/staydigital/truvest-identity-management/api"
	"bitbucket.org/staydigital/truvest-identity-management/docs"
)

// @title Truvest Identity Management Service APIs
// @version 1.0
// @description This is an RBAC based full fledge API serice for managing users, roles and permissions in the system
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	var sb strings.Builder
	// sb.WriteString(os.Getenv("APP_PROTOCOL"))
	// sb.WriteString("://")
	sb.WriteString(os.Getenv("APP_HOST"))
	sb.WriteString(":")
	sb.WriteString(os.Getenv("APP_PORT"))
	// Programmatically set swagger info
	docs.SwaggerInfo.Host = sb.String()

	api.Run()
}