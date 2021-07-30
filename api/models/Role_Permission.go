package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type Role_Permission struct {
	PermissionID	uint32	`json:"permissionid"`
	RoleID			uint32	`json:"roleid"`
}

type Role_Permission_Payload struct {
	Permissions		[]uint32	`json:"permissions"`
}

func (rp *Role_Permission) Prepare() {
	rp.PermissionID = rp.PermissionID
	rp.RoleID = rp.RoleID
}

func (rp *Role_Permission) Validate() error {
	if rp.PermissionID == 0 {
		return errors.New("Required PermissionID")
	}
	if rp.RoleID == 0 {
		return errors.New("Required RoleID")
	}
	return nil
}

func (rp *Role_Permission) SavePermissionToRole(db *gorm.DB) (error) {
	var err error
	err = db.Debug().Model(&Role_Permission{}).Create(&rp).Error
	if err != nil {
		return err
	}
	return nil
}

func (rp *Role_Permission) DeleteRoleFromPermission(db *gorm.DB, rid uint32, pid uint32) (int64, error) {

	db = db.Debug().Model(&Role_Permission{}).Where("role_id = ? and permission_id = ?", rid, pid).Take(&Role_Permission{}).Delete(&Role_Permission{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Role/Permission not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}