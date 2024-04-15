package main

import (
	"booking/pkg/booking/models"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type config struct {
	port string
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	models models.Models
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	var cfg config

	DATABASE_URL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"))

	flag.StringVar(&cfg.port, "port", ":"+os.Getenv("APPLICATION_PORT"), "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", DATABASE_URL, "PostgreSQL DSN")
	flag.Parse()

	// Connect to DB
	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	app := &application{
		config: cfg,
		models: models.NewModels(db),
	}

	app.run()
}

func (app *application) run() {
	r := mux.NewRouter()

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

	log.Printf("Starting server on %s\n", app.config.port)
	err := http.ListenAndServe(app.config.port, r)
	log.Fatal(err)
}

func openDB(cfg config) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config // struct.
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
