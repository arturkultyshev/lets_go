package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

type Reviews struct {
	Id              int       `json:"id"`
	UserId          int       `json:"userid"`
	HotelId         int       `json:"hotelid"`
	Rating          float64   `json:"rating"`
	Comment         string    `json:"comment"`
	PublicationDate time.Time `json:"publication_date"`
}

type ReviewsModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m ReviewsModel) AddReview(review *Reviews) error {
	query := `
        INSERT INTO reviews (hotel_id, user_id, rating, comment, publication_date) 
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
        `
	args := []interface{}{review.HotelId, review.UserId, review.Rating, review.Comment, review.PublicationDate}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&review.Id)
}

func (m ReviewsModel) GetReviews(hotelId int) ([]*Reviews, error) {
	query := `
        SELECT id, hotel_id, user_id, rating, comment, publication_date
        FROM reviews
        WHERE hotel_id = $1
        `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, hotelId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*Reviews
	for rows.Next() {
		var review Reviews
		err := rows.Scan(&review.Id, &review.HotelId, &review.UserId, &review.Rating, &review.Comment, &review.PublicationDate)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, &review)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}

func (m ReviewsModel) GetReviewById(id int) (*Reviews, error) {
	query := `
        SELECT id, hotel_id, user_id, rating, comment
        FROM reviews
        WHERE id = $1
        `
	var review Reviews
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&review.Id, &review.HotelId, &review.UserId, &review.Rating, &review.Comment)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no review found with ID: %d", id)
		}
		return nil, err
	}
	return &review, nil
}

func (m ReviewsModel) UpdateReview(review *Reviews) error {
	query := `
        UPDATE reviews
        SET rating = $2, comment = $3
        WHERE id = $1
        RETURNING hotel_id, user_id
        `
	args := []interface{}{review.Id, review.Rating, review.Comment}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&review.HotelId, &review.UserId)
}

func (m ReviewsModel) DeleteReview(id int) error {
	query := `
        DELETE FROM reviews
        WHERE id = $1
        `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

func (m ReviewsModel) GetAverageRating(hotelId int) (float64, error) {
	query := `
		SELECT AVG(rating) FROM reviews
		WHERE hotel_id = $1
		`
	var avgRating float64
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, hotelId).Scan(&avgRating)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("no reviews found for hotel ID: %d", hotelId)
		}
		return 0, err
	}

	return avgRating, nil
}
