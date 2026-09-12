package ownRedis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func NewClient(addr string) *redis.Client {
	// Create Client
	rdb := redis.NewClient(&redis.Options{Addr: addr})

	// Test Client by Pinging
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Redis ping: %v", err)
	}

	log.Printf("Connected to redis at %s", addr)

	return rdb
}
