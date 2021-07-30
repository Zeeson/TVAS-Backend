package seed

import (
	"log"
	"os"

	"bitbucket.org/staydigital/truvest-identity-management/api/models"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

var ConstID, err = uuid.NewRandom()

var roles = []models.Role{
	{
		Name:   "System Admin",
		Description: "System Admin Role",
		CreatedBy: ConstID,
		UpdatedBy: ConstID,
	},
	{
		Name:   "Admin",
		Description: "Generic Admin Role",
		CreatedBy: ConstID,
		UpdatedBy: ConstID,
	},
	{
		Name:   "Customer",
		Description: "Customer Role",
		CreatedBy: ConstID,
		UpdatedBy: ConstID,
	},
	{
		Name:   "Guest",
		Description: "Guest Role",
		CreatedBy: ConstID,
		UpdatedBy: ConstID,
	},
}

var users_roles = []models.User_Role{
	{
		UserID: ConstID,
		RoleID:	1,
	},
}

var permissions = []models.Permission{
	{
		Name: 	"SYSTEM_ADMIN",
		CreatedBy: ConstID,
		UpdatedBy: ConstID,
	},
	{
		Name: 	"USERS_CREATE",
		CreatedBy: ConstID,
		UpdatedBy: ConstID,
	},
	{
		Name:   "USERS_VIEW",
		CreatedBy: ConstID,
		UpdatedBy: ConstID,
	},
	{
		Name:   "USERS_UNLOCK",
		CreatedBy: ConstID,
		UpdatedBy: ConstID,
	},
	{
		Name:   "USERS_ACTIVATE",
		CreatedBy: ConstID,
		UpdatedBy: ConstID,
	},
	{
		Name:   "USERS_CHANGE_PWD",
		CreatedBy: ConstID,
		UpdatedBy: ConstID,
	},
	{
		Name:   "USERS_ASSIGN_TO_ROLE",
		CreatedBy: ConstID,
		UpdatedBy: ConstID,
	},	
	{
		Name:   "ROLES_VIEW",
		CreatedBy: ConstID,
		UpdatedBy: ConstID,
	},	
	{
		Name:   "MANAGE_ROLES",
		CreatedBy: ConstID,
		UpdatedBy: ConstID,
	},	
	{
		Name:   "VIEW_PERMISSION",
		CreatedBy: ConstID,
		UpdatedBy: ConstID,
	},	
	{
		Name:   "MANAGE_PERMISSION",
		CreatedBy: ConstID,
		UpdatedBy: ConstID,
	},
	{
		Name:   "PERMISSION_ASSIGN_TO_ROLE",
		CreatedBy: ConstID,
		UpdatedBy: ConstID,
	},
}

var roles_permissions = []models.Role_Permission{
	{
		PermissionID: 1,
		RoleID:	1,
	},
	{
		PermissionID: 2,
		RoleID: 1,
	},
	{
		PermissionID: 3,
		RoleID: 1,
	},
	{
		PermissionID: 4,
		RoleID: 1,
	},
	{
		PermissionID: 5,
		RoleID: 1,
	},
	{
		PermissionID: 6,
		RoleID: 1,
	},
	{
		PermissionID: 7,
		RoleID: 1,
	},
	{
		PermissionID: 8,
		RoleID: 1,
	},
	{
		PermissionID: 9,
		RoleID: 1,
	},
	{
		PermissionID: 10,
		RoleID: 1,
	},
	{
		PermissionID: 11,
		RoleID: 1,
	},
	{
		PermissionID: 12,
		RoleID: 1,
	},
	{
		PermissionID: 2,
		RoleID: 2,
	},
	{
		PermissionID: 3,
		RoleID: 2,
	},
	{
		PermissionID: 4,
		RoleID: 2,
	},
	{
		PermissionID: 5,
		RoleID: 2,
	},
	{
		PermissionID: 6,
		RoleID: 2,
	},
	{
		PermissionID: 7,
		RoleID: 2,
	},
	{
		PermissionID: 8,
		RoleID: 2,
	},
	{
		PermissionID: 9,
		RoleID: 2,
	},
	{
		PermissionID: 10,
		RoleID: 2,
	},
	{
		PermissionID: 11,
		RoleID: 2,
	},
	{
		PermissionID: 12,
		RoleID: 2,
	},
}

// Load DB with seed data
func Load(db *gorm.DB, username string, firstname string, lastname string, email string, password string) {

	if os.Getenv("GORM_AUTOMIGRATE") == "true" {
		// Enable below if you want to drop the existing tables
		err := db.Debug().DropTableIfExists(&models.Role{}, &models.User{}, &models.User_Role{}, &models.Permission{}, &models.Role_Permission{}, models.User_Device{}, models.Refresh_Token{}).Error
		if err != nil {
			log.Fatalf("cannot drop table: %v", err)
		}
	
		err = db.Debug().AutoMigrate(&models.User{}, &models.Role{}, &models.User_Role{}, &models.Permission{}, &models.Role_Permission{}, models.User_Device{}, models.Refresh_Token{}).Error
		if err != nil {
			log.Fatalf("cannot migrate table: %v", err)
		}
	
		var users = []models.User{
			{
				ID:			 ConstID,
				UserName:	 username,
				FirstName:	 firstname,
				LastName:	 lastname,
				Email:		 email,
				Password:	 password,
				Provider: 	 "local",
				CreatedBy:	 ConstID,
				UpdatedBy:	 ConstID,
			},
		}
	
		for i := range users {
			err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
			if err != nil {
				log.Fatalf("cannot seed users table: %v", err)
			}
		}
	
		for i := range roles {
			err = db.Debug().Model(&models.Role{}).Create(&roles[i]).Error
			if err != nil {
				log.Fatalf("cannot seed roles table: %v", err)
			}
		}
	
		for i := range users_roles {
			err = db.Debug().Model(&models.User_Role{}).Create(&users_roles[i]).Error
			if err != nil {
				log.Fatalf("cannot seed users_roles table: %v", err)
			}
		}
	
		for i := range permissions {
			err = db.Debug().Model(&models.Permission{}).Create(&permissions[i]).Error
			if err != nil {
				log.Fatalf("cannot seed permissions table: %v", err)
			}
		}
	
		for i := range roles_permissions {
			err = db.Debug().Model(&models.Role_Permission{}).Create(&roles_permissions[i]).Error
			if err != nil {
				log.Fatalf("cannot seed roles_permissions table: %v", err)
			}
		}
	}
}