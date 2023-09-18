package set

import "math/rand"

var empty struct{}

type StringSet struct {
	data map[string]struct{}
}

func InitStringSet() *StringSet {
	return &StringSet{data: make(map[string]struct{})}
}

func (s *StringSet) Add(v string) {
	s.data[v] = empty
}

func (s *StringSet) Remove(v string) {
	delete(s.data, v)
}

func (s *StringSet) size() int {
	return len(s.data)
}

// RandomN returns n random strings from the set
// The returned elements could repeat. if n > size
// of set then all elements are returned.
func (s *StringSet) RandomN(n int) []string {
	size := s.size()

	// This is inneficcient, but it's a start
	keys := make([]string, 0, size)
	for k := range s.data {
		keys = append(keys, k)
	}

	if n > size {
		return keys
	}

	randomKeys := make([]string, 0, size)
	for i := 0; i < n; i++ {
		ind := rand.Intn(size)
		randomKeys = append(randomKeys, keys[ind])
	}

	return randomKeys
}
