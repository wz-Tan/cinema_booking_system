package booking

import (
	"cinema_booking_system/internal/utils"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type handler struct {
	svc *Service
}

// Frontend Passes This In
type holdRequest struct {
	UserID string `json:"user_id"`
}

func NewHandler(svc *Service) *handler {
	return &handler{
		svc: svc,
	}
}

func (h *handler) HoldSeat(w http.ResponseWriter, r *http.Request) {
	movieID := r.PathValue("movieID")
	seatID := r.PathValue("seatID")

	var req holdRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, "Error Decoding Hold Request")
		log.Print("Error Decoding Hold Request")
		return
	}

	userBooking := Booking{
		UserID:  req.UserID,
		MovieID: movieID,
		SeatID:  seatID,
	}

	session, err := h.svc.Book(userBooking)

	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, "Error Booking Seat")
		log.Print("Error Booking Seat")
		return
	}

	// Hold Seat Response
	type holdResponse struct {
		SessionID string `json:"session_id"`
		MovieID   string `json:"movieID"`
		SeatID    string `json:"seat_id"`
		ExpiresAt string `json:"expires_at"`
	}

	// Session is An Unconfirmed Booking
	utils.WriteJSON(w, http.StatusOK, holdResponse{
		SessionID: session.ID,
		MovieID:   session.MovieID,
		SeatID:    seatID,
		ExpiresAt: session.ExpiresAt.Format(time.RFC3339),
	})
}

func (h *handler) ListSeats(w http.ResponseWriter, r *http.Request) {
	movieID := r.PathValue("movieID") // get via {movieID}

	bookings, err := h.svc.ListBookings(movieID)

	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, "Error Listing Booking")
		return
	}

	// List Out the Booked Seats for Frontend
	seats := make([]utils.SeatInfo, 0, len(bookings))
	for _, b := range bookings {
		seats = append(seats, utils.SeatInfo{
			SeatID:    b.SeatID,
			UserID:    b.UserID,
			Booked:    true,
			Confirmed: b.Status == "confirmed",
		})
	}

	// Return all bookings
	utils.WriteJSON(w, http.StatusOK, seats)
}

func (h *handler) ConfirmSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("sessionID")

	var req holdRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, "Error Decoding Hold Request")
		log.Print("Error Decoding Hold Request")
		return
	}

	// Non User
	if req.UserID == "" {
		utils.WriteJSON(w, http.StatusBadRequest, "No User ID Provided")
		return
	}

	if err := h.svc.ConfirmSession(r.Context(), sessionID, req.UserID); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, "Failed to Release Entry")
		log.Print("Error confirming session", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Successfully Released Session")
}

func (h *handler) ReleaseSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("sessionID")

	var req holdRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, "Error Decoding Hold Request")
		log.Print("Error Decoding Hold Request")
		return
	}

	// Non User
	if req.UserID == "" {
		utils.WriteJSON(w, http.StatusBadRequest, "No User ID Provided")
		return
	}

	if err := h.svc.ReleaseSession(r.Context(), sessionID, req.UserID); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, "Failed to Release Entry")
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Successfully Released Session")
}
