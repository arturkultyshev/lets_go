package main

import (
	"booking/pkg/booking/models"
	"booking/pkg/booking/validator"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (app *application) respondWithError(w http.ResponseWriter, code int, message string) {
	app.respondWithJSON(w, code, map[string]string{"error": message})
}

func (app *application) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (app *application) createHotelHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name    string `json:"name"`
		Country string `json:"country"`
		City    string `json:"city"`
		Street  string `json:"street"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	hotels := &models.Hotels{
		Name:    input.Name,
		Country: input.Country,
		City:    input.City,
		Street:  input.Street,
	}

	err = app.models.Hotels.AddHotel(hotels)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, hotels)
}

func (app *application) getHotelHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["hotelId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid menu ID")
		return
	}

	hotel, err := app.models.Hotels.GetHotelById(id)
	fmt.Println(hotel)
	if err != nil {
		fmt.Println(err, hotel)
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, hotel)
}

func (app *application) getHotelsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string
		CostFrom int
		CostTo   int
		models.Filters
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

	//if models.ValidateFilters(v, input.Filters); !v.Valid() {
	//	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
	//	return
	//}
	hotels, _, err := app.models.Hotels.GetAll(input.Name, input.CostFrom, input.CostTo, input.Filters)

	if err != nil {
		fmt.Println(err)
		app.respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Respond with JSON
	app.respondWithJSON(w, http.StatusOK, hotels)
}

func (app *application) updateHotelHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["hotelId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid menu ID")
		return
	}

	hotel, err := app.models.Hotels.GetHotelById(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
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
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
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
		hotel.Rating = *input.Rating
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
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, hotel)
}

func (app *application) deleteHotelHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["hotelId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid menu ID")
		return
	}

	err = app.models.Hotels.DeleteHotel(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
