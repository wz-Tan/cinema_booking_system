package booking

import (
	"context"
	"encoding/json"
	"errors"
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
func parseSessionKey(id string) string {
	return fmt.Sprintf("Session %s", id)
}

// Commit Booking
func (s *RedisStore) ConfirmSession(ctx context.Context, sessionID string, userID string) error {

	session, sessionKey, err := s.getSession(ctx, sessionID, userID)

	if err != nil {
		return errors.New(err.Error())
	}

	// Make Permanent (Extract Booking, Edit Status, Overwrite)
	// Change Status
	session.Status = "confirmed"

	// Turn Back into JSON
	val, _ := json.Marshal(session)

	// Overwrite in DB
	overwriteRes := s.rdb.Set(ctx, sessionKey, val, 0)
	if overwriteRes.Err() != nil {
		return errors.New("Error Overwriting Data")
	}

	// Persist Key
	if err := s.rdb.Persist(ctx, sessionID).Err(); err != nil {
		return errors.New("Error Confirming Key")
	}
	return nil
}

// Release Booking
func (s *RedisStore) ReleaseSession(ctx context.Context, sessionID string, userID string) error {

	_, sessionKey, err := s.getSession(ctx, sessionID, userID)

	if err != nil {
		return errors.New(err.Error())
	}

	// Remove from Bookings Table
	if err := s.rdb.Del(ctx, sessionKey).Err(); err != nil {
		return errors.New("Error Releasing Booking")
	}

	// Remove from Reverse Key Table
	if err := s.rdb.Del(ctx, parseSessionKey(sessionID)).Err(); err != nil {
		return errors.New("Error Releasing Key")
	}
	return nil
}

// Get Session / Booking from DB after Verifying (Return booking -> For editing, string -> booking key)
func (s *RedisStore) getSession(ctx context.Context, sessionID string, userID string) (Booking, string, error) {
	// Special Processing
	sessionID = parseSessionKey(sessionID)

	// Get the DB Key
	sessionKey, err := s.rdb.Get(ctx, sessionID).Result()

	if err != nil {
		return Booking{}, "", errors.New("Error Acquiring DB Key")
	}

	// Get Booking JSON object
	b, err := s.rdb.Get(ctx, sessionKey).Result()

	// Parse it
	session, err := parseBooking(b)
	if err != nil {
		return Booking{}, "", errors.New("Error converting DB field to session.")
	}

	// Verify Ownership
	if session.UserID != userID {
		return Booking{}, "", errors.New("Invalid Request from Non User")
	}

	return session, sessionKey, nil
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

	// Set Session
	s.rdb.Set(ctx, parseSessionKey(id), key, defaultHoldTTL)

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
func (s *RedisStore) Book(b Booking) (Booking, error) {
	session, err := s.hold(b)

	// Error Holding
	if err != nil {
		return Booking{}, err
	}

	// No Problem Holding
	log.Printf("Session created %v", session)

	return session, nil
}

func (s *RedisStore) ListBookings(movieID string) ([]Booking, error) {
	pattern := fmt.Sprintf("seat:%s:*", movieID)
	var sessions []Booking

	ctx := context.Background()

	// Array of Seats That Match the Pattern
	iter := s.rdb.Scan(ctx, 0, pattern, 0).Iterator()

	// Counter for Iterator
	for iter.Next(ctx) {
		// Returns JSON String
		val, err := s.rdb.Get(ctx, iter.Val()).Result()
		if err != nil {
			continue
		}

		// Parse String into Booking Struct
		session, err := parseBooking(val)
		if err != nil {
			continue
		}

		sessions = append(sessions, session)

	}

	return sessions, nil
}

func parseBooking(val string) (Booking, error) {
	booking := Booking{}

	// Parses JSON into Bytes then Write Into Struct
	if err := json.Unmarshal([]byte(val), &booking); err != nil {
		return Booking{}, err
	}

	// New Booking Object -> We Dont Need the Expiration Time for a Confirmed Session
	return Booking{
		ID:      booking.ID,
		MovieID: booking.MovieID,
		SeatID:  booking.SeatID,
		UserID:  booking.UserID,
		Status:  booking.Status,
	}, nil
}
