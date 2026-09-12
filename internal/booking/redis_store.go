package booking

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// Key Design for Redis
// seat:{movieID}:{seatID}
// session:{sessionID} -> Seat Key

// Time to Live for Hold
const defaultHoldTTL = 2 * time.Minute

type RedisStore struct {
	rdb *redis.Client
}

func NewRedisStore(rdb *redis.Client) *RedisStore {
	return &RedisStore{
		rdb: rdb,
	}
}

// String Parser for Session
func sessionKey(id string) string {
	return fmt.Sprintf("Session %s", id)
}

// Create Booking
func (s *RedisStore) hold(b Booking) (Booking, error) {
	id := uuid.New().String()
	now := time.Now()
	ctx := context.Background()
	key := fmt.Sprintf("seat:%s:%s", b.MovieID, b.SeatID)

	// Reassign New ID to Booking
	b.ID = id

	// Convert Booking into JSON object
	val, _ := json.Marshal(b)

	// Redis CLI (Key is the format above, value is our struct)
	res := s.rdb.SetArgs(ctx, key, val, redis.SetArgs{
		Mode: "NX", // run if not exists
		TTL:  defaultHoldTTL,
	})

	// Status of Creating Booking in Redis
	if res.Val() != "OK" {
		return Booking{}, ErrSeatAlreadyBooked
	}

	log.Printf("No problem creating booking")

	// Set Session
	s.rdb.Set(ctx, sessionKey(id), key, defaultHoldTTL)

	log.Printf("No problem creating sesion")

	return Booking{
		ID:        id,
		MovieID:   b.MovieID,
		SeatID:    b.SeatID,
		UserID:    b.UserID,
		Status:    "held",
		ExpiresAt: now.Add(defaultHoldTTL),
	}, nil
}

// Commit to Booking
func (s *RedisStore) Book(b Booking) error {
	session, err := s.hold(b)

	// Error Holding
	if err != nil {
		return err
	}

	// No Problem Holding
	log.Printf("Session booked %v", session)

	return nil
}

func (s *RedisStore) ListBookings(movieID string) []Booking {
	return []Booking{}
}
