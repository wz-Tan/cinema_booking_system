package booking

import (
	"context"
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
	Book(b Booking) (Booking, error)
	ListBookings(movieID string) ([]Booking, error)
	ReleaseSession(ctx context.Context, sessionID string, userID string) error
	ConfirmSession(ctx context.Context, sessionID string, userID string) error
}

var (
	ErrSeatAlreadyBooked = errors.New("Seat is already booked!")
)
