package models

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type User_Role struct {
	UserID	uuid.UUID	`gorm:"type:uuid"json:"userid"`
	RoleID	uint32		`json:"roleid"`
}

type User_Role_Payload struct {
	Users	[]uuid.UUID	`json:"users"`
}

func (ur *User_Role) Prepare() {
	ur.UserID = ur.UserID
	ur.RoleID = ur.RoleID
}

func (ur *User_Role) Validate() error {
	if len(ur.UserID) < 1 {
		return errors.New("Required UserID")
	}
	if ur.RoleID == 0 {
		return errors.New("Required RoleID")
	}
	return nil
}

func (ur *User_Role) SaveUserToRole(db *gorm.DB) (error) {
	var err error
	err = db.Debug().Model(&User_Role{}).Create(&ur).Error
	if err != nil {
		return err
	}
	return nil
}

func (ur *User_Role) DeleteUsersFromRole(db *gorm.DB, rid uint32, uid uuid.UUID) (int64, error) {

	db = db.Debug().Model(&User_Role{}).Where("role_id = ? and user_id = ?", rid, uid).Take(&User_Role{}).Delete(&User_Role{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Role/User not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}