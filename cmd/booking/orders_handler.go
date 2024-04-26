package main

import (
	"booking/pkg/booking/model"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func (app *application) createOrderHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserId         int       `json:"userid"`
		HotelId        int       `json:"hotelid"`
		StartDate      time.Time `json:"start_date"`
		EndDate        time.Time `json:"end_date"`
		AdditionalInfo string    `json:"additional_info,omitempty"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	order := &model.Orders{
		UserId:         input.UserId,
		HotelId:        input.HotelId,
		StartDate:      input.StartDate,
		EndDate:        input.EndDate,
		AdditionalInfo: input.AdditionalInfo,
		CreationDate:   time.Now(),
	}

	err = app.models.Orders.AddOrder(order)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"order": order}, nil)
}

func (app *application) getOrdersHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	orders, err := app.models.Orders.GetOrders(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	app.writeJSON(w, http.StatusOK, envelope{"orders": orders}, nil)
}

func (app *application) updateOrderHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		StartDate      time.Time `json:"start_date"`
		EndDate        time.Time `json:"end_date"`
		AdditionalInfo string    `json:"additional_info"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		fmt.Println(err)
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}

	order := &model.Orders{
		Id:             id,
		StartDate:      input.StartDate,
		EndDate:        input.EndDate,
		AdditionalInfo: input.AdditionalInfo,
	}

	err = app.models.Orders.UpdateOrder(order)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"order": order}, nil)
}
func (app *application) deleteOrderHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Orders.DeleteOrder(id)

	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	app.writeJSON(w, http.StatusOK, envelope{"message": "order deleted"}, nil)
}
