package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

type Hotels struct {
	Id             int     `json:"id"`
	Name           string  `json:"name"`
	Country        string  `json:"country"`
	City           string  `json:"city"`
	Street         string  `json:"street"`
	Rating         float64 `json:"rating"`
	Capacity       int     `json:"capacity"`
	Cost           int     `json:"cost"`
	PhotoUrl       string  `json:"photo_url"`
	AdditionalInfo string  `json:"additional_info"`
}

type HotelsModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m HotelsModel) AddHotel(hotels *Hotels) error {
	// Insert a new menu item into the database.
	query := `
		INSERT INTO hotels (name, country, city, street) 
		VALUES ($1, $2, $3, $4)
		RETURNING id
		`
	args := []interface{}{hotels.Name, hotels.Country, hotels.City, hotels.Street}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&hotels.Id)
}

func (m HotelsModel) GetHotelById(id int) (*Hotels, error) {
	// Retrieve a specific menu item based on its ID.
	query := `
		SELECT id, name, country, city, street
		FROM hotels
		WHERE id = $1
		`
	var hotels Hotels
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&hotels.Id, &hotels.Name, &hotels.Country, &hotels.City, &hotels.Street)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no hotel found with ID: %d", id)
		}
		return nil, err
	}
	return &hotels, nil
}

func (m HotelsModel) UpdateHotel(hotels *Hotels) error {
	// Update a specific menu item in the database.
	query := `
		UPDATE hotels
		SET capacity = $2, cost = $3, additional_info = $4
		WHERE id = $1
		RETURNING name, city
		`
	args := []interface{}{hotels.Id, hotels.Capacity, hotels.Cost, hotels.AdditionalInfo}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&hotels.Name)
}

func (m HotelsModel) DeleteHotel(id int) error {
	// Delete a specific menu item from the database.
	query := `
		DELETE FROM hotels
		WHERE id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}
