package booking

type MemoryStore struct {
	// Dict of String to Booking
	bookings map[string]Booking
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		bookings: map[string]Booking{},
	}
}

func (s *MemoryStore) Book(b Booking) error {
	return nil
}

func (s *MemoryStore) ListBookings(movieID string) []string {
	return nil
}
