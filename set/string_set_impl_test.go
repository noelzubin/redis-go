package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Add(t *testing.T) {
	assert := assert.New(t)
	s := InitStringSet()
	s.Add("one")
	res := s.RandomN(1)
	assert.Equal("one", res[0])
}

func Test_Random2(t *testing.T) {
	assert := assert.New(t)
	s := InitStringSet()
	s.Add("one")
	s.Add("two")
	res := s.RandomN(3)
	assert.Equal(2, len(res))
	assert.True(res[0] == "one" || res[0] == "two")
	assert.True((res[1] == "one" || res[1] == "two") && res[1] != res[0])
}

func Test_Random_1(t *testing.T) {
	assert := assert.New(t)
	s := InitStringSet()
	s.Add("one")
	s.Add("two")
	s.Add("three")
	s.Add("four")
	res := s.RandomN(1)
	assert.Equal(1, len(res))
	assert.True(res[0] == "one" || res[0] == "two" || res[0] == "three" || res[0] == "four")
}

func Test_Random_Remove(t *testing.T) {
	assert := assert.New(t)
	s := InitStringSet()
	s.Add("one")
	res := s.RandomN(1)
	assert.Equal(1, len(res))

	s.Remove("one")
	res = s.RandomN(1)
	assert.Equal(0, len(res))
}
