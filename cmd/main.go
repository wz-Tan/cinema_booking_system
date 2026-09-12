package main

import (
	"cinema_booking_system/internal/adapters/ownRedis"
	"cinema_booking_system/internal/booking"
	"cinema_booking_system/internal/utils"
	"log"
	"net/http"
)

func main() {
	// Handler
	mux := http.NewServeMux()

	// Serve File
	mux.Handle("GET /", http.FileServer(http.Dir("static")))

	// Serve Functions
	mux.HandleFunc("GET /movies", utils.ListMovies)

	store := booking.NewRedisStore(ownRedis.NewClient("localhost:6379"))
	svc := booking.NewService(store)
	bookingHandler := booking.NewHandler(svc)

	mux.HandleFunc("GET /movies/:movieID/seats", bookingHandler.ListSeats)

	// Simple Server
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}

	log.Print("Server running on 8080")
}
