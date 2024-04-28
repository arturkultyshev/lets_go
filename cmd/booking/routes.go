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

	hotels1 := r.PathPrefix("/hotels").Subrouter()

	// Create a new hotel
	hotels1.HandleFunc("", app.requirePermissions("hotel:write", app.createHotelHandler)).Methods("POST")
	// Get hotels
	hotels1.HandleFunc("", app.getHotelsHandler).Methods("GET")
	// Get a specific hotel, update a specific hotel, and delete a specific hotel
	hotels1.HandleFunc("/{id:[0-9]+}", app.getHotelHandler).Methods("GET")
	hotels1.HandleFunc("/{id:[0-9]+}", app.requirePermissions("hotel:update", app.updateHotelHandler)).Methods("PUT")
	hotels1.HandleFunc("/{id:[0-9]+}", app.requirePermissions("hotel:delete", app.deleteHotelHandler)).Methods("DELETE")
	// Create a new review
	hotels1.HandleFunc("/reviews", app.requirePermissions("review:write", app.createReviewHandler)).Methods("POST")
	// Get reviews
	hotels1.HandleFunc("/reviews/{id:[0-9]+}", app.getReviewsByHotelHandler).Methods("GET")
	// Update review
	hotels1.HandleFunc("/reviews/{id:[0-9]+}", app.requirePermissions("review:update", app.updateReviewHandler)).Methods("PUT")
	// Delete review
	hotels1.HandleFunc("/reviews/{id:[0-9]+}", app.requirePermissions("review:delete", app.deleteReviewHandler)).Methods("DELETE")

	orders1 := r.PathPrefix("/orders").Subrouter()
	// Create a new order
	orders1.HandleFunc("", app.requirePermissions("order:write", app.createOrderHandler)).Methods("POST")
	// Get orders by user
	orders1.HandleFunc("", app.requirePermissions("order:read", app.getOrdersHandler)).Methods("GET")
	// Get a specific order, update a specific order, and delete a specific order
	orders1.HandleFunc("/{id:[0-9]+}", app.requirePermissions("order:update", app.updateOrderHandler)).Methods("PUT")
	orders1.HandleFunc("/{id:[0-9]+}", app.requirePermissions("order:delete", app.deleteOrderHandler)).Methods("DELETE")

	users1 := r.PathPrefix("/users/").Subrouter()
	// User handlers with Authentication
	users1.HandleFunc("/register", app.registerUserHandler).Methods("POST")
	users1.HandleFunc("/activated", app.activateUserHandler).Methods("PUT")
	users1.HandleFunc("/login", app.createAuthenticationTokenHandler).Methods("POST")

	// Wrap the router with the panic recovery middleware and rate limit middleware.
	return app.authenticate(r)
}
