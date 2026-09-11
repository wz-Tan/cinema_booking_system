package booking

type MemoryStore struct {
	// Dict of Seat ID to Booking (Assume One Hall)
	bookings map[string]Booking
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		bookings: map[string]Booking{},
	}
}

func (s *MemoryStore) Book(b Booking) error {
	// Seat Already Booked
	if _, exists := s.bookings[b.SeatID]; exists {
		return ErrSeatAlreadyBooked
	}

	// Seat Not Booked
	s.bookings[b.SeatID] = b
	return nil
}

func (s *MemoryStore) ListBookings(movieID string) []Booking {
	var movieBookings []Booking

	for _, booking := range s.bookings {
		if booking.MovieID == movieID {
			movieBookings = append(movieBookings, booking)
		}
	}

	return movieBookings
}
