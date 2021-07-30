package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type User_Device struct {
	ID        			uuid.UUID  	`gorm:"primary_key;type:uuid" json:"id"`
	DeviceID  			string    	`gorm:"size:255;not null;unique" json:"device_id"`
	IsRefreshActive  	bool    	`gorm:"default:true" json:"is_refresh_active"`
	UserID 				uuid.UUID  	`gorm:"type:uuid;not null" json:"user_id"`
	CreatedAt 			time.Time 	`gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	CreatedBy 			uuid.UUID 	`gorm:"type:uuid;not null" json:"created_by"`
	UpdatedAt 			time.Time 	`gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	UpdatedBy 			uuid.UUID 	`gorm:"type:uuid;not null" json:"updated_by"`
}

func (ud *User_Device) SaveUserDevice(db *gorm.DB) (error) {
	var err error
	err = db.Debug().Create(&ud).Error
	if err != nil {
		return err
	}
	return nil
}

func (ud *User_Device) FindUserDeviceByUserID(db *gorm.DB, uid uuid.UUID) (*User_Device, error) {
	var err error
	err = db.Debug().Model(User_Device{}).Where("user_id = ?", uid).Take(&ud).Error
	if err != nil {
		return &User_Device{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User_Device{}, errors.New("User Device Not Found")
	}
	return ud, nil
}

func (ud *User_Device) UpdateAUserDevice(db *gorm.DB, uid uuid.UUID, tuid uuid.UUID, isRefreshActive bool) (*User_Device, error) {

	db = db.Debug().Model(&User_Device{}).Where("user_id = ?", uid).Take(&User_Device{}).UpdateColumns(
		map[string]interface{}{
			"is_refresh_active": isRefreshActive,
			"updated_at": time.Now(),
			"updated_by": tuid,
		},
	)
	if db.Error != nil {
		return &User_Device{}, db.Error
	}
	err := db.Debug().Model(&User_Device{}).Where("user_id = ?", uid).Take(&ud).Error
	if err != nil {
		return &User_Device{}, err
	}
	return ud, nil
}

func (ud *User_Device) DeleteAUserDevice(db *gorm.DB, uid uuid.UUID) (error) {

	db = db.Debug().Model(&User_Device{}).Where("user_id = ?", uid).Take(&User_Device{}).Delete(&User_Device{})

	if db.Error != nil {
		return db.Error
	}
	return nil
}