package booking

import (
	"errors"
	"time"
)

type Booking struct {
	ID        string
	MovieID   string
	SeatID    string
	UserID    string
	Status    string
	ExpiresAt time.Time
}

// Handler for Booking and Listing Bookings
type BookingStore interface {
	Book(b Booking) error
	ListBookings(movieID string) ([]Booking, error)
}

var (
	ErrSeatAlreadyBooked = errors.New("Seat is already booked!")
)
