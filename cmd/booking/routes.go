package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// routes is our main application's router.
func (app *application) routes() http.Handler {
	r := mux.NewRouter()
	// Convert the app.notFoundResponse helper to a http.Handler using the http.HandlerFunc()
	// adapter, and then set it as the custom error handler for 404 Not Found responses.
	r.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)

	// Convert app.methodNotAllowedResponse helper to a http.Handler and set it as the custom
	// error handler for 405 Method Not Allowed responses
	r.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)

	// healthcheck
	// TODO: Add a healthcheck endpoint to the router.
	// r.HandleFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	v1 := r.PathPrefix("/").Subrouter()

	// Menu Singleton
	// Create a new menu
	v1.HandleFunc("/hotels", app.createHotelHandler).Methods("POST")
	// Get hotels
	v1.HandleFunc("/hotels", app.getHotelsHandler).Methods("GET")
	// Get a specific hotel
	v1.HandleFunc("/hotels/{hotelId:[0-9]+}", app.getHotelHandler).Methods("GET")
	// Update a specific menu
	v1.HandleFunc("/hotels/{hotelId:[0-9]+}", app.updateHotelHandler).Methods("PUT")
	// Delete a specific menu
	v1.HandleFunc("/hotels/{hotelId:[0-9]+}", app.deleteHotelHandler).Methods("DELETE")

	users1 := r.PathPrefix("/users/").Subrouter()
	// User handlers with Authentication
	users1.HandleFunc("/register", app.registerUserHandler).Methods("POST")
	users1.HandleFunc("/activated", app.activateUserHandler).Methods("PUT")
	users1.HandleFunc("/login", app.createAuthenticationTokenHandler).Methods("POST")

	// Wrap the router with the panic recovery middleware and rate limit middleware.
	return app.authenticate(r)
}
