package models

import (
	"database/sql"
	"log"
)

type Orders struct {
	Id             int    `json:"id"`
	UserId         int    `json:"userid"`
	HotelId        int    `json:"hotelid"`
	StartDate      string `json:"startdate"`
	EndDate        string `json:"enddate"`
	CreationDate   string `json:"creationdate"`
	AdditionalInfo string `json:"additionalinfo"`
}

type OrdersModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}
