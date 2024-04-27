package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

type Orders struct {
	Id             int       `json:"id"`
	UserId         int       `json:"userid"`
	HotelId        int       `json:"hotelid"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	CreationDate   time.Time `json:"created_at"`
	AdditionalInfo string    `json:"additional_info"`
}

type OrdersModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m OrdersModel) AddOrder(order *Orders) error {
	query := `
		INSERT INTO orders (hotel_id, user_id, start_date, end_date, additional_info, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, hotel_id, user_id, start_date, end_date, created_at, additional_info
		`
	args := []interface{}{order.HotelId, order.UserId, order.StartDate, order.EndDate, order.AdditionalInfo, order.CreationDate}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&order.Id, &order.HotelId, &order.UserId, &order.StartDate, &order.EndDate, &order.CreationDate, &order.AdditionalInfo)
}

func (m OrdersModel) GetOrders(userId int64) ([]*Orders, error) {
	query := `
		SELECT id, hotel_id, user_id, start_date, end_date, created_at, additional_info
		FROM orders
		WHERE user_id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*Orders
	for rows.Next() {
		var order Orders
		err := rows.Scan(&order.Id, &order.HotelId, &order.UserId, &order.StartDate, &order.EndDate, &order.CreationDate, &order.AdditionalInfo)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (m OrdersModel) GetOrder(orderId int) (*Orders, error) {
	query := `
		SELECT id, hotel_id, user_id, start_date, end_date, created_at, additional_info
		FROM orders
		WHERE id = $1
		`
	var order Orders
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, orderId)
	err := row.Scan(&order.Id, &order.HotelId, &order.UserId, &order.StartDate, &order.EndDate, &order.CreationDate, &order.AdditionalInfo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no order found with id %d", orderId)
		}
		return nil, err

	}

	return &order, nil
}

func (m OrdersModel) UpdateOrder(order *Orders) error {
	query := `
		UPDATE orders
		SET start_date = $1, end_date = $2, additional_info = $3
		WHERE id = $4
		RETURNING user_id, hotel_id, start_date, end_date, additional_info
		`
	args := []interface{}{order.StartDate, order.EndDate, order.AdditionalInfo, order.Id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&order.UserId, &order.HotelId, &order.StartDate, &order.EndDate, &order.AdditionalInfo)
}

func (m OrdersModel) DeleteOrder(orderId int) error {
	query := `
        DELETE FROM orders
        WHERE id = $1
        `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, orderId)
	return err
}
