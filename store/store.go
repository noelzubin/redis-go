package store

import "time"

// ScoreMember represents a member of a sorted set
type ScoreMember struct {
	score  int64
	member string
}

// NewScoreMember creates a new ScoreMember
func NewScoreMember(score int64, member string) ScoreMember {
	return ScoreMember{score: score, member: member}
}

// Store interface wraps the methods that a store must implement
type Store interface {
	// Ping the store
	Ping() *string
	// Get a value for key from store
	Get(k string) *string
	// Set a value for key with optional expiry time
	Set(k string, v string, e *time.Time)
	// Del deletes a Key from store
	Del(keys ...string) int
	// Expire updates the expiry time for a key
	Expire(k string, seconds int) int
	// Keys returns all keys
	Keys(k string) []string
	// ZAdd adds a member to a sorted set
	ZAdd(k string, s []ScoreMember) int
	// ZRange returns a range of members from a sorted set
	ZRange(k string, start int, stop int, withScores bool) []string
}
