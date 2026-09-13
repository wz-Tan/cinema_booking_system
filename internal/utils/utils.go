package utils

import (
	"encoding/json"
	"net/http"
)

// Info Required by frontend to map movies
type movieResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Rows        int    `json:"Rows"`
	SeatsPerRow int    `json:"seats_per_row"`
}

// General Info About a Seat (Holding is Performed by Redis via TTL)
type SeatInfo struct {
	SeatID    string `json:"seat_id"`
	UserID    string `json:"user_id"`
	Booked    bool   `json:"booked"`
	Confirmed bool   `json:"confirmed"`
}

var movies = []movieResponse{
	{ID: "1", Title: "Spiderman", Rows: 10, SeatsPerRow: 10},
	{ID: "2", Title: "Batman", Rows: 10, SeatsPerRow: 10},
}

func ListMovies(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, http.StatusOK, movies)
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
