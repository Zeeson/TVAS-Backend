package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Refresh_Token struct {
	ID        			uuid.UUID  	`gorm:"primary_key;type:uuid" json:"id"`
	DeviceID  			string    	`gorm:"size:255;not null;unique" json:"device_id"`
	RefreshCount	  	int64    	`gorm:"default:0" json:"refresh_count"`
	Token 				uuid.UUID  	`gorm:"type:uuid;not null;unique" json:"user_id"`
	ExpiryDate 			int64	 	`gorm:"default:0" json:"expiry_date"`
	CreatedAt 			time.Time 	`gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	CreatedBy 			uuid.UUID 	`gorm:"type:uuid;not null" json:"created_by"`
	UpdatedAt 			time.Time 	`gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	UpdatedBy 			uuid.UUID 	`gorm:"type:uuid;not null" json:"updated_by"`
}

func (rt *Refresh_Token) SaveRefreshToken(db *gorm.DB) (error) {
	var err error
	err = db.Debug().Create(&rt).Error
	if err != nil {
		return err
	}
	return nil
}

func (rt *Refresh_Token) FindRefreshTokenByDeviceID(db *gorm.DB, did uuid.UUID) (error) {
	var err error
	err = db.Debug().Model(Refresh_Token{}).Where("device_id = ?", did).Take(&rt).Error
	if err != nil {
		return err
	}
	if gorm.IsRecordNotFoundError(err) {
		return errors.New("User Device Not Found")
	}
	return nil
}

func (rt *Refresh_Token) UpdateARefreshToken(db *gorm.DB, token uuid.UUID, tuid uuid.UUID) (error) {

	db = db.Debug().Model(&Refresh_Token{}).Where("token = ?", token).Take(&rt).UpdateColumns(
		map[string]interface{}{
			"refresh_count": rt.RefreshCount + 1,
			"token": uuid.New(),
			"updated_at": time.Now(),
			"updated_by": tuid,
		},
	)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (rt *Refresh_Token) DeleteARefreshToken(db *gorm.DB, did string) (error) {

	db = db.Debug().Model(&Refresh_Token{}).Where("device_id = ?", did).Take(&Refresh_Token{}).Delete(&Refresh_Token{})

	if db.Error != nil {
		return db.Error
	}
	return nil
}
