package models

import (
    "gorm.io/gorm"
    "golang.org/x/crypto/bcrypt"
)

type User struct {
    gorm.Model
    Username string `gorm:"unique;not null" json:"username"`
    Email    string `gorm:"unique;not null" json:"email"`
    Password string `gorm:"not null" json:"-"`
    Cars     []Car  `json:"cars"`
}

// HashPassword hashes the user's password before saving
func (u *User) HashPassword() error {
    bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
    if err != nil {
        return err
    }
    u.Password = string(bytes)
    return nil
}

// CheckPassword verifies the password
func (u *User) CheckPassword(password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
    return err == nil
}
