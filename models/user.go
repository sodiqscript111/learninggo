package models

import (
	"learninggo/db"
)

func AddUser(user db.User) error {
	result := db.DB.Create(&user)
	return result.Error
}
