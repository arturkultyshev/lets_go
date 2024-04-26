package model

import (
	"database/sql"
	"log"
	"time"
)

type Reviews struct {
	Id              int       `json:"id"`
	UserId          int       `json:"userid"`
	HotelId         int       `json:"hotelid"`
	RatingDecimal   float64   `json:"ratingdecimal"`
	PublicationDate time.Time `json:"publicationdate"`
	Comment         string    `json:"comment"`
}

type ReviewsModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}
