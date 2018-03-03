package models

import (
//"encoding/json"
//"time"
//
//"github.com/go-sql-driver/mysql"
"golang.org/x/crypto/bcrypt"
//"github.com/jinzhu/gorm"
	"gotest/database"
	"time"
)

// User represents a user account for public visibility (used for public endpoints)
// Its MarshalJSON function wont expose its role.
type User struct {
	ID        uint `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time	`json:"updatedAt"`
	DeletedAt *time.Time	`json:"deletedAt"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	Password  string         `json:"password"`
	Admin     bool           `json:"admin"`
}

type UserHelper struct {
	db *database.MySQLDB
}

// AuthUser represents a user account for private visibility (used for login and update response)
// Its MarshalJSON function will expose its role.
type AuthUser struct {
	*User
	Admin bool `json:"admin"`
}

func (u *User) SetPassword(password string) {
	pwhash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	u.Password = string(pwhash)
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return false
	}

	return true
}

func NewUserHelper(db *database.MySQLDB) *UserHelper {
	return &UserHelper{db}
}

func (u *UserHelper) Exist(email string) bool {
	var user User
	if u.db.Where("email = ?", email).First(&user).RecordNotFound() {
		return false
	} else {
		return true
	}
}

func (u *User) IsAdmin() bool {
	return u.Admin == true
}

func (uh *UserHelper) FindByEmail(email string) (*User, error) {
	user := User{}
	err := uh.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user,nil
}

func (uh *UserHelper) FindById(id int) (*User, error) {
	user := User{}
	err := uh.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user,nil
}