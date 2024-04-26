package main

import (
	"booking/pkg/booking/model"
	"booking/pkg/booking/validator"
	"errors"
	"net/http"
)

func (app *application) createHotelHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name           string `json:"name"`
		Country        string `json:"country"`
		City           string `json:"city"`
		Street         string `json:"street"`
		Capacity       int    `json:"capacity,omitempty"`
		Cost           int    `json:"cost,omitempty"`
		PhotoUrl       string `json:"photo_url,omitempty"`
		AdditionalInfo string `json:"additional_info,omitempty"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}

	hotels := &model.Hotels{
		Name:           input.Name,
		Country:        input.Country,
		City:           input.City,
		Street:         input.Street,
		Capacity:       input.Capacity,
		Cost:           input.Cost,
		PhotoUrl:       input.PhotoUrl,
		AdditionalInfo: input.AdditionalInfo,
	}

	err = app.models.Hotels.AddHotel(hotels)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"hotel": hotels}, nil)
}

func (app *application) getHotelHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	hotel, err := app.models.Hotels.GetHotelById(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"hotel": hotel}, nil)
}

func (app *application) getHotelsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string
		CostFrom int
		CostTo   int
		model.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	input.Name = app.readStrings(qs, "name", "")
	input.CostFrom = app.readInt(qs, "costFrom", 0, v)
	input.CostTo = app.readInt(qs, "costTo", 0, v)

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 5, v)

	input.Filters.Sort = app.readStrings(qs, "sort", "id")

	input.Filters.SortSafeList = []string{
		// ascending sort values
		"id", "name", "cost",
		// descending sort values
		"-id", "-name", "-cost",
	}

	if model.ValidateFilters(v, input.Filters); !v.Valid() {
		app.errorResponse(w, r, http.StatusUnprocessableEntity, v.Errors)
		return
	}

	hotels, metadata, err := app.models.Hotels.GetAll(input.Name, input.CostFrom, input.CostTo, input.Filters)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"hotels": hotels, "metadata": metadata}, nil)
}

func (app *application) updateHotelHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	hotel, err := app.models.Hotels.GetHotelById(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Name           *string  `json:"name"`
		Country        *string  `json:"country"`
		City           *string  `json:"city"`
		Street         *string  `json:"street"`
		Rating         *float64 `json:"nutritionValue"`
		Capacity       *int     `json:"capacity"`
		Cost           *int     `json:"cost"`
		PhotoUrl       *string  `json:"photo_url"`
		AdditionalInfo *string  `json:"additional_info"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		hotel.Name = *input.Name
	}

	if input.Country != nil {
		hotel.Country = *input.Country
	}

	if input.City != nil {
		hotel.City = *input.City
	}

	if input.Street != nil {
		hotel.Street = *input.Street
	}

	if input.Rating != nil {
		avgRating, err := app.models.Reviews.GetAverageRating(hotel.Id)
		if err != nil {
			hotel.Rating = 0
		} else {
			hotel.Rating = avgRating
		}
	}

	if input.Capacity != nil {
		hotel.Capacity = *input.Capacity
	}

	if input.Cost != nil {
		hotel.Cost = *input.Cost
	}

	if input.PhotoUrl != nil {
		hotel.PhotoUrl = *input.PhotoUrl
	}

	if input.AdditionalInfo != nil {
		hotel.AdditionalInfo = *input.AdditionalInfo
	}

	err = app.models.Hotels.UpdateHotel(hotel)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"hotel": hotel}, nil)
}

func (app *application) deleteHotelHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Hotels.DeleteHotel(id)

	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	app.writeJSON(w, http.StatusOK, envelope{"message": "hotel deleted"}, nil)
}
