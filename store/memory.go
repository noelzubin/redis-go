package store

import (
	"fmt"
	"noelzubin/redis-go/set"
	"time"

	"github.com/wangjia184/sortedset"
)

// Value is the wrapper for value stored in the store
type Value struct {
	value  interface{}
	expiry *time.Time
}

// checks if the value is expired
func (v *Value) isExpired() bool {
	return v.expiry != nil && v.expiry.Before(time.Now())
}

// InMemStore is the main inmemory implementation of the store
type InMemStore struct {
	data           map[string]Value
	keysWithExpiry set.IStringSet
}

// InitStore initializes a new InMemStore
func InitStore(keysWithExpiry set.IStringSet) *InMemStore {
	return &InMemStore{
		data:           make(map[string]Value),
		keysWithExpiry: keysWithExpiry,
	}
}

func (s *InMemStore) Ping() *string {
	pong := "PONG"
	return &pong
}

func (s *InMemStore) Get(k string) *string {
	value, ok := s.data[k]

	if !ok {
		return nil
	}

	if value.isExpired() {
		delete(s.data, k)
		s.keysWithExpiry.Remove(k)
		return nil
	}

	res, ok := value.value.(string)

	if !ok {
		return nil
	}

	return &res
}

func (s *InMemStore) Set(k string, v string, e *time.Time) {
	value := Value{value: v, expiry: e}
	s.data[k] = value

	if e != nil {
		s.keysWithExpiry.Add(k)
	}
}

func (s *InMemStore) Del(keys ...string) int {
	delCount := 0
	for _, k := range keys {
		if (s.Get(k)) != nil {
			delCount++
			delete(s.data, k)
			s.keysWithExpiry.Remove(k)
		}
	}
	return delCount
}

func (s *InMemStore) Expire(k string, seconds int) int {
	value, ok := s.data[k]

	if !ok {
		fmt.Println("Key does not exist")
	}

	expiry := time.Now().Add(time.Duration(seconds) * time.Second)
	value.expiry = &expiry
	s.data[k] = value
	s.keysWithExpiry.Add(k)

	return 1
}

func (s *InMemStore) Keys(k string) []string {
	keys := make([]string, 0, len(s.data))

	for k, v := range s.data {
		if !v.isExpired() {
			keys = append(keys, k)
		} else {
			delete(s.data, k)
			s.keysWithExpiry.Remove(k)
		}
	}

	return keys
}

func (s *InMemStore) ZAdd(k string, scoreMembers []ScoreMember) int {
	value, ok := s.data[k]

	var set *sortedset.SortedSet

	if ok {
		set = value.value.(*sortedset.SortedSet)
	} else {
		set = sortedset.New()
	}

	for _, scoreMember := range scoreMembers {
		set.AddOrUpdate(scoreMember.member, sortedset.SCORE(scoreMember.score), nil)
	}

	s.data[k] = Value{value: set, expiry: nil}

	return len(scoreMembers)
}

func (s *InMemStore) ZRange(k string, start int, stop int, withScores bool) []string {
	values := make([]string, 0)

	value, ok := s.data[k]

	if !ok {
		return values
	}

	set := value.value.(*sortedset.SortedSet)

	fmt.Println("starst")
	for _, v := range set.GetByRankRange(start, stop, false) {
		fmt.Println(v.Key(), v.Score())
		values = append(values, v.Key())
		if withScores {
			values = append(values, fmt.Sprintf("%d", v.Score()))
		}
	}
	fmt.Println("end")
	fmt.Println(values)

	return values
}

func (s *InMemStore) CleanUp() {
	needsCleanup := true
	// Cleanup until there are very few expired keys
	for needsCleanup {
		needsCleanup = false
		randomKeys := s.keysWithExpiry.RandomN(20)
		expiredCount := 0

		// Delete the expired keys
		for _, k := range randomKeys {
			val, ok := s.data[k]
			if !ok {
				continue
			} else {
				if val.isExpired() {
					expiredCount++
					delete(s.data, k)
					s.keysWithExpiry.Remove(k)
				}
			}
		}

		// If more than 25% is expired
		if expiredCount > 5 {
			needsCleanup = true
		}
	}
}
