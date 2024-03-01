package models

import (
	"database/sql"
	"log"
)

type Users struct {
	Id          int    `json:"id"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	PhoneNumber string `json:"phonenumber"`
	Email       string `json:"email"`
	IsAdmin     bool   `json:"isadmin"`
}

type UsersModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}
