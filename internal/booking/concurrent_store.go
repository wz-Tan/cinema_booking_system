package booking

import "sync"

type ConcurrentStore struct {
	bookings     map[string]Booking
	sync.RWMutex // for mutex
}

func NewConcurrentStore() *ConcurrentStore {
	return &ConcurrentStore{
		bookings: map[string]Booking{},
	}
}

func (s *ConcurrentStore) Book(b Booking) error {
	s.Lock()
	defer s.Unlock()

	// Seat Already Booked
	if _, exists := s.bookings[b.SeatID]; exists {
		return ErrSeatAlreadyBooked
	}

	// Seat Not Booked
	s.bookings[b.SeatID] = b
	return nil
}

func (s *ConcurrentStore) ListBookings(movieID string) []Booking {
	s.RLock()
	defer s.RUnlock()

	var movieBookings []Booking

	for _, booking := range s.bookings {
		if booking.MovieID == movieID {
			movieBookings = append(movieBookings, booking)
		}
	}

	return movieBookings
}
