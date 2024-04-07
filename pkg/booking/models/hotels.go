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

func (m HotelsModel) GetAll(name string, from, to int, filters Filters) ([]*Hotels, Metadata, error) {

	// Retrieve all menu items from the database.
	query := fmt.Sprintf(
		`
		SELECT count(*) OVER(), id, name, country, city, street, capacity, cost
		FROM hotels
		WHERE (LOWER(name) = LOWER($1) OR $1 = '')
		AND (cost >= $2 OR $2 = 0)
		AND (cost <= $3 OR $3 = 0)
		ORDER BY %s %s, id ASC
		LIMIT $4 OFFSET $5
		`,
		filters.sortColumn(), filters.sortDirection())

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Organize our four placeholder parameter values in a slice.
	args := []interface{}{name, from, to, filters.limit(), filters.offset()}

	// log.Println(query, title, from, to, filters.limit(), filters.offset())
	// Use QueryContext to execute the query. This returns a sql.Rows result set containing
	// the result.
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	// Importantly, defer a call to rows.Close() to ensure that the result set is closed
	// before GetAll returns.
	defer func() {
		if err := rows.Close(); err != nil {
			m.ErrorLog.Println(err)
		}
	}()

	// Declare a totalRecords variable
	totalRecords := 0

	var hotels []*Hotels
	for rows.Next() {
		var hotel Hotels
		err := rows.Scan(&totalRecords, &hotel.Id, &hotel.Name, &hotel.Country, &hotel.City, &hotel.Street, &hotel.Capacity, &hotel.Cost)
		if err != nil {
			return nil, Metadata{}, err
		}

		// Add the Movie struct to the slice
		hotels = append(hotels, &hotel)
	}

	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	// Generate a Metadata struct, passing in the total record count and pagination parameters
	// from the client.
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	// If everything went OK, then return the slice of the movies and metadata.
	return hotels, metadata, nil
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

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&hotels.Name, &hotels.City)
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
