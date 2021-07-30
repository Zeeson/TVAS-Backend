package models

import (
	"errors"
	"html"
	"strings"
	"time"
	"unsafe"

	"github.com/badoux/checkmail"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        	uuid.UUID  	`gorm:"primary_key;type:uuid" json:"id"`
	UserName  	string    	`gorm:"size:255;not null;unique" json:"username"`
	FirstName  	string    	`gorm:"size:255;not null" json:"firstname"`
	LastName  	string    	`gorm:"size:255;not null" json:"lastname"`
	Email     	string    	`gorm:"size:100;not null;unique" json:"email"`
	Password  	string    	`gorm:"size:100;not null;" json:"password,omitempty"`
	CreatedAt 	time.Time 	`gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	CreatedBy 	uuid.UUID 	`gorm:"type:uuid;not null" json:"created_by"`
	UpdatedAt 	time.Time 	`gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	UpdatedBy 	uuid.UUID 	`gorm:"type:uuid;not null" json:"updated_by"`
	Enabled	  	bool		`gorm:"default:true" json:"enabled"`
	Provider	string		`gorm:"size:255;not null" json:"provider"`
	Roles	  	[]*Role		`gorm:"many2many:user_roles" json:"roles,omitempty"`
}

type CreateUser struct {
	UserName  	string    	`gorm:"size:255;not null;unique" json:"username"`
	FirstName  	string    	`gorm:"size:255;not null" json:"firstname"`
	LastName  	string    	`gorm:"size:255;not null" json:"lastname"`
	Email     	string    	`gorm:"size:100;not null;unique" json:"email"`
	Password  	string    	`gorm:"size:100;not null;" json:"password,omitempty"`
}

type EnableUser struct {
	Enabled	  	bool		`gorm:"default:true" json:"enabled"`
}

type UpdateUser struct {
	UserName  	string    	`gorm:"size:255;not null;unique" json:"username"`
	FirstName  	string    	`gorm:"size:255;not null" json:"firstname"`
	LastName  	string    	`gorm:"size:255;not null" json:"lastname"`
}

type UserResponse struct {
	ID        	uuid.UUID  	`gorm:"primary_key;type:uuid" json:"id"`
	UserName  	string    	`gorm:"size:255;not null;unique" json:"username"`
	FirstName  	string    	`gorm:"size:255;not null" json:"firstname"`
	LastName  	string    	`gorm:"size:255;not null" json:"lastname"`
	Email     	string    	`gorm:"size:100;not null;unique" json:"email"`
	CreatedAt 	time.Time 	`gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	CreatedBy 	uuid.UUID 	`gorm:"type:uuid;not null" json:"created_by"`
	UpdatedAt 	time.Time 	`gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	UpdatedBy 	uuid.UUID 	`gorm:"type:uuid;not null" json:"updated_by"`
	Enabled	  	bool		`gorm:"default:true" json:"enabled"`
	Provider	string		`gorm:"size:255;not null" json:"provider"`
	Roles	  	[]*Role		`gorm:"many2many:user_roles" json:"roles,omitempty"`
}

type LoginResponse struct {
	AccessToken		string    	`gorm:"size:255;not null;" json:"access_token"`
	RefreshToken	uuid.UUID  	`gorm:"size:100;not null;" json:"refresh_token"`
	ExpiryDate 		int64		`gorm:"default:0" json:"expiry_date"`
	DeviceID		string    	`gorm:"size:255;not null;" json:"device_id"`
}

type Set_User_Password_Payload struct {
	Password  	string    	`gorm:"size:100;not null;" json:"password,omitempty"`
}

type Forgot_User_Password_Payload struct {
	Email     	string    	`gorm:"size:100;not null;unique" json:"email"`
	Password  	string    	`gorm:"size:100;not null;" json:"password,omitempty"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *User) BeforeSave() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) Prepare(tuid uuid.UUID) {
	u.ID = uuid.New()
	u.UserName = html.EscapeString(strings.TrimSpace(u.UserName))
	u.FirstName = html.EscapeString(strings.TrimSpace(u.FirstName))
	u.LastName = html.EscapeString(strings.TrimSpace(u.LastName))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.Password = html.EscapeString(strings.TrimSpace(u.Password))
	u.CreatedAt = time.Now()
	u.CreatedBy = tuid
	u.UpdatedAt = time.Now()
	u.UpdatedBy = tuid
	u.Enabled = true
	u.Provider = "local"
}

func (u *User) PrepareSignUp() {
	id := uuid.New()
	u.ID = id
	u.UserName = html.EscapeString(strings.TrimSpace(u.UserName))
	u.FirstName = html.EscapeString(strings.TrimSpace(u.FirstName))
	u.LastName = html.EscapeString(strings.TrimSpace(u.LastName))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.Password = html.EscapeString(strings.TrimSpace(u.Password))
	u.CreatedAt = time.Now()
	u.CreatedBy = id
	u.UpdatedAt = time.Now()
	u.UpdatedBy = id
	u.Enabled = true
	u.Provider = "local"
}

func PrepareResponse(u *User) UserResponse {
	ur := UserResponse{}
	ur.ID = u.ID
	ur.UserName = u.UserName
	ur.FirstName = u.FirstName
	ur.LastName = u.LastName
	ur.Email = u.Email
	ur.CreatedAt = u.CreatedAt
	ur.CreatedBy = u.CreatedBy
	ur.UpdatedAt = u.UpdatedAt
	ur.UpdatedBy = u.UpdatedBy
	ur.Enabled = u.Enabled
	ur.Provider = u.Provider
	ur.Roles = u.Roles
	return ur
}

func PrepareResponses(u *[]User) []UserResponse {
	urs := []UserResponse{}
	for i := range *u {
		user := (*u)[i]
		urs = append(urs, PrepareResponse((*User)(unsafe.Pointer(&user))))
	}
	return urs
}

func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "login":
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil

	default:
		if u.UserName == "" {
			return errors.New("Required UserName")
		}
		if u.FirstName == "" {
			return errors.New("Required FirstName")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil
	}
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	var err error
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) FindAllUsers(db *gorm.DB) (*[]User, error) {
	var err error
	users := []User{}
	err = db.Debug().Model(&User{}).Preload("Roles").Preload("Roles.Permissions").Find(&users).Error
	if err != nil {
		return &[]User{}, err
	}
	return &users, err
}

func (u *User) FindUserByID(db *gorm.DB, uid uuid.UUID) (*User, error) {
	var err error
	err = db.Debug().Model(User{}).Where("id = ?", uid).Preload("Roles").Preload("Roles.Permissions").Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User Not Found")
	}
	return u, nil
}

func (u *User) FindUserByEmail(db *gorm.DB, email string) (*User, error) {
	var err error
	err = db.Debug().Model(User{}).Where("email = ?", email).Preload("Roles").Preload("Roles.Permissions").Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) UpdateAUser(db *gorm.DB, uid uuid.UUID, tuid uuid.UUID) (*User, error) {

	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"user_name":  u.UserName,
			"first_name":  u.FirstName,
			"last_name": u.LastName,
			"updated_at": time.Now(),
			"updated_by": tuid,
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}
	// This is the display the updated user
	err := db.Debug().Model(&User{}).Where("id = ?", uid).Preload("Roles").Preload("Roles.Permissions").Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) EnableDisableUser(db *gorm.DB, uid uuid.UUID, tuid uuid.UUID) (*User, error) {

	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"enabled":  u.Enabled,
			"updated_at": time.Now(),
			"updated_by": tuid,
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}
	// This will display the updated user
	err := db.Debug().Model(&User{}).Where("id = ?", uid).Preload("Roles").Preload("Roles.Permissions").Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) DeleteAUser(db *gorm.DB, uid uuid.UUID) (int64, error) {

	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (sup *Set_User_Password_Payload) ResetPassword(db *gorm.DB, uid uuid.UUID, tuid uuid.UUID) error {

	// To hash the password
	if len(sup.Password) < 1 {
		return errors.New("Required Password")
	}
	hashedPassword, err := Hash(sup.Password)
	if err != nil {
		return err
	}
	password := string(hashedPassword)

	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"password":  password,
			"updated_at": time.Now(),
			"updated_by": tuid,
		},
	)

	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (fup *Forgot_User_Password_Payload) ForgetPassword(db *gorm.DB) error {

	// To hash the password
	if len(fup.Password) < 1 {
		return errors.New("Required Password")
	}
	hashedPassword, err := Hash(fup.Password)
	if err != nil {
		return err
	}
	password := string(hashedPassword)

	if len(fup.Email) < 1 {
		return errors.New("Required Email")
	}
	if err := checkmail.ValidateFormat(fup.Email); err != nil {
		return errors.New("Invalid Email")
	}

	db = db.Debug().Model(&User{}).Where("email = ?", fup.Email).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"password":  password,
			"updated_at": time.Now(),
		},
	)

	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (u *User) WhoAmI(db *gorm.DB, uid uuid.UUID) (*User, error) {
	var err error
	err = db.Debug().Model(User{}).Where("id = ?", uid).Preload("Roles").Preload("Roles.Permissions").Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User Not Found")
	}
	return u, nil
}