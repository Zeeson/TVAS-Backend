package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Role struct {
	ID          	uint32		    			`gorm:"primary_key;auto_increment" json:"id"`
	Name     		string		    			`gorm:"size:255;not null;unique" json:"name"`
	Description 	string		    			`gorm:"size:255;not null;" json:"description"`
	CreatedAt 		time.Time		 			`gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	CreatedBy	 	uuid.UUID				 	`gorm:"type:uuid;not null" json:"created_by"`
	UpdatedAt 		time.Time		 			`gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	UpdatedBy	 	uuid.UUID				 	`gorm:"type:uuid;not null" json:"updated_by"`
	Permissions		[]*Permission				`gorm:"many2many:role_permissions" json:"permissions"`
}

type CreateRole struct {
	Name     		string    		`gorm:"size:255;not null;unique" json:"name"`
	Description 	string    		`gorm:"size:255;not null;" json:"description"`
}

func (r *Role) Prepare(tuid uuid.UUID) {
	r.ID = 0
	r.Name = html.EscapeString(strings.TrimSpace(r.Name))
	r.Description = html.EscapeString(strings.TrimSpace(r.Description))
	r.CreatedAt = time.Now()
	r.CreatedBy = tuid
	r.UpdatedAt = time.Now()
	r.UpdatedBy = tuid
}

func (r *Role) Validate() error {
	if r.Name == "" {
		return errors.New("Required Name")
	}
	if r.Description == "" {
		return errors.New("Required Description")
	}
	return nil
}

func (r *Role) SaveRole(db *gorm.DB) (*Role, error) {
	var err error
	err = db.Debug().Model(&Role{}).Create(&r).Error
	if err != nil {
		return &Role{}, err
	}
	return r, nil
}

func (r *Role) FindAllRoles(db *gorm.DB) (*[]Role, error) {
	var err error
	roles := []Role{}
	err = db.Debug().Model(&Role{}).Preload("Permissions").Find(&roles).Error
	if err != nil {
		return &[]Role{}, err
	}
	return &roles, nil
}

func (r *Role) FindRoleByID(db *gorm.DB, rid uint32) (*Role, error) {
	var err error
	err = db.Debug().Model(&Role{}).Where("id = ?", rid).Preload("Permissions").Take(&r).Error
	if err != nil {
		return &Role{}, err
	}
	return r, nil
}

func (r *Role) UpdateARole(db *gorm.DB, rid uint32, tuid uuid.UUID) (*Role, error) {

	var err error
	err = db.Debug().Model(&Role{}).Where("id = ?", r.ID).Updates(
		Role{
				Name: r.Name,
				Description: r.Description,
				UpdatedAt: time.Now(),
				UpdatedBy: tuid,
			}).Error
	if err != nil {
		return &Role{}, err
	}
	return r, nil
}

func (r *Role) DeleteARole(db *gorm.DB, rid uint32) (int64, error) {

	db = db.Debug().Model(&Role{}).Where("id = ?", rid).Take(&Role{}).Delete(&Role{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Role not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}