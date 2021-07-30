package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Permission struct {
	ID          	uint32    		`gorm:"primary_key;auto_increment" json:"id"`
	Name     		string    		`gorm:"size:255;not null;unique" json:"name"`
	CreatedAt 		time.Time		`gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	CreatedBy	 	uuid.UUID		`gorm:"type:uuid;not null" json:"created_by"`
	UpdatedAt 		time.Time		`gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	UpdatedBy	 	uuid.UUID		`gorm:"type:uuid;not null" json:"updated_by"`
}

type CreatePermission struct {
	Name     		string    		`gorm:"size:255;not null;unique" json:"name"`
}

func (p *Permission) Prepare(tuid uuid.UUID) {
	p.ID = 0
	p.Name = html.EscapeString(strings.TrimSpace(p.Name))
	p.CreatedAt = time.Now()
	p.CreatedBy = tuid
	p.UpdatedAt = time.Now()
	p.UpdatedBy = tuid
}

func (p *Permission) Validate() error {
	if p.Name == "" {
		return errors.New("Required Name")
	}
	return nil
}

func (p *Permission) SavePermission(db *gorm.DB) (*Permission, error) {
	var err error
	err = db.Debug().Model(&Permission{}).Create(&p).Error
	if err != nil {
		return &Permission{}, err
	}
	return p, nil
}

func (p *Permission) FindAllPermissions(db *gorm.DB) (*[]Permission, error) {
	var err error
	permissions := []Permission{}
	err = db.Debug().Model(&Permission{}).Find(&permissions).Error
	if err != nil {
		return &[]Permission{}, err
	}
	return &permissions, nil
}

func (p *Permission) FindPermissionByID(db *gorm.DB, pid uint32) (*Permission, error) {
	var err error
	err = db.Debug().Model(&Permission{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Permission{}, err
	}
	return p, nil
}

func (p *Permission) UpdateAPermission(db *gorm.DB, pid uint32, tuid uuid.UUID) (*Permission, error) {

	var err error
	err = db.Debug().Model(&Permission{}).Where("id = ?", pid).Updates(
		Permission{
				Name: p.Name,
				UpdatedAt: time.Now(),
				UpdatedBy: tuid,
			}).Error
	if err != nil {
		return &Permission{}, err
	}
	p.ID = pid
	return p, nil
}

func (p *Permission) DeleteAPermission(db *gorm.DB, pid uint32) (int64, error) {

	db = db.Debug().Model(&Permission{}).Where("id = ?", pid).Take(&Permission{}).Delete(&Permission{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Role not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}