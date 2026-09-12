package utils

import (
	"encoding/json"
	"net/http"
)

// Info Required by frontend to map movies
type movieResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Rows        int    `json:"rows"`
	SeatsPerRow int    `json:"seats_per_row"`
}

// General Info About a Seat (Holding is Performed by Redis via TTL)
type seatInfo struct {
	SeatID string `json:"seat_id"`
	UserID string `json:"user_id"`
	Booked bool   `json:"booked"`
}

var movies = []movieResponse{
	{ID: "1", Title: "Spiderman", Rows: 2, SeatsPerRow: 5},
	{ID: "2", Title: "Batman", Rows: 6, SeatsPerRow: 7},
}

func ListMovies(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, http.StatusOK, movies)
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
