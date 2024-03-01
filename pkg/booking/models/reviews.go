package models

import (
	"database/sql"
	"log"
)

type Reviews struct {
	Id              int     `json:"id"`
	UserId          int  `json:"userid"`
	HotelId         int  `json:"hotelid"`
	RatingDecimal   float64 `json:"ratingdecimal"`
	PublicationDate string     `json:"publicationdate"`
	Comment         string     `json:"comment"`
}

type ReviewsModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}
