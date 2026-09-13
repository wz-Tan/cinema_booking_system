package booking

import "context"

type Service struct {
	store BookingStore
}

func NewService(store BookingStore) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) Book(booking Booking) (Booking, error) {
	return s.store.Book(booking)
}

func (s *Service) ListBookings(movieID string) ([]Booking, error) {
	return s.store.ListBookings(movieID)
}

func (s *Service) ReleaseSession(ctx context.Context, sessionID string, userID string) error {
	return s.store.ReleaseSession(ctx, sessionID, userID)
}

func (s *Service) ConfirmSession(ctx context.Context, sessionID string, userID string) error {
	return s.store.ConfirmSession(ctx, sessionID, userID)
}
