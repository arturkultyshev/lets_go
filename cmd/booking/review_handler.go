package main

import (
	"booking/pkg/booking/model"
	"errors"
	"net/http"
	"time"
)

func (app *application) createReviewHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserId  int     `json:"userid"`
		HotelId int     `json:"hotelid"`
		Rating  float64 `json:"rating"`
		Comment string  `json:"comment"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}

	review := &model.Reviews{
		UserId:          input.UserId,
		HotelId:         input.HotelId,
		Rating:          input.Rating,
		Comment:         input.Comment,
		PublicationDate: time.Now(),
	}

	err = app.models.Reviews.AddReview(review)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"review": review}, nil)
}

func (app *application) getReviewsByHotelHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	review, err := app.models.Reviews.GetReviews(id)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"reviews": review}, nil)
}

func (app *application) updateReviewHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	review, err := app.models.Reviews.GetReviewById(id)
	if err != nil {
		app.notFoundResponse(w, r)
	}

	var input struct {
		Rating  *float64 `json:"rating"`
		Comment *string  `json:"comment"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Rating != nil {
		review.Rating = *input.Rating
	}

	if input.Comment != nil {
		review.Comment = *input.Comment
	}

	err = app.models.Reviews.UpdateReview(review)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"Review:": review}, nil)
}

func (app *application) deleteReviewHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Reviews.DeleteReview(id)

	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	app.writeJSON(w, http.StatusOK, envelope{"message": "review deleted"}, nil)
}
