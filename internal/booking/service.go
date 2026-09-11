package booking

type Service struct {
	store BookingStore
}

func NewService(store BookingStore) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) Book(booking Booking) error {
	return s.store.Book(booking)
}
