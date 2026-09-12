package booking

import (
	"cinema_booking_system/internal/utils"
	"net/http"
)

type handler struct {
	svc *Service
}

func NewHandler(svc *Service) *handler {
	return &handler{
		svc: svc,
	}
}

func (h *handler) ListSeats(w http.ResponseWriter, r *http.Request) {
	movieID := r.PathValue("movieID") // get via :movieID

	bookings, err := h.svc.ListBookings(movieID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, "Error Listing Booking")
	}

	// Return all bookings
	utils.WriteJSON(w, http.StatusOK, bookings)
}
