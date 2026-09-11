package booking

import "errors"

type Booking struct {
	ID      string
	MovieID string
	SeatID  string
	UserID  string
	Status  string
}

// Handler for Booking and Listing Bookings
type BookingStore interface {
	Book(b Booking) error
	ListBookings(movieID string) []Booking
}

var (
	ErrSeatAlreadyBooked = errors.New("Seat is already booked!")
)
