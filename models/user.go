package models

import (
	"bytes"
	"crypto/sha512"
	"errors"
	"html"
	"log"
	"strings"

	"github.com/lovrog05/task-manager-backend/utils/token"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserID   uint   `gorm:"primaryKey"`
	Username string `gorm:"size:255;not null;unique" json:"username"`
	Password string `gorm:"size:255;not null;" json:"-"`
	FmcToken string `gorm:"size:255;not null;" json:"-"`
}

func GetUserByID(uid uint) (User, error) {

	var u User

	if err := DB.First(&u, uid).Error; err != nil {
		return u, errors.New("user not found")
	}

	u.PrepareGive()

	return u, nil

}

func (u *User) PrepareGive() {
	u.Password = ""
}

func VerifyPassword(password, hashedPassword string) error {
	password = PreparePasswordInput(password)
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(username string, password string) (string, error) {

	var err error

	u := User{}

	err = DB.Model(User{}).Where("username = ?", username).Take(&u).Error

	if err != nil {
		return "", err
	}

	err = VerifyPassword(password, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	token, err := token.GenerateToken(u.UserID)

	if err != nil {
		return "", err
	}

	return token, nil

}

func (u *User) UpdateFmcToken() error {

	if err := DB.Model(User{}).Where("user_id = ?", u.UserID).Update("fmc_token", u.FmcToken).Error; err != nil {
		return err
	}

	return nil
}

func (u *User) SaveUser() (*User, error) {

	var err error = nil
	err = DB.Create(&u).Error
	if err != nil {
		log.Println(err)
		return &User{}, err
	}
	return u, nil
}

func (u *User) BeforeSave() error {

	//turn password into hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	//remove spaces in username
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

	return nil

}

func PreparePasswordInput(plainText string) string {
	// Creates a SHA512 hash, trimmed to 64 characters, so that it fits in bcrypt
	hashedInput := sha512.Sum512_256([]byte(plainText))
	// Bcrypt terminates at NULL bytes, so we need to trim these away
	trimmedHash := bytes.Trim(hashedInput[:], "\x00")
	return string(trimmedHash)
}
