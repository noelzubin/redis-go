package utils

import (
	"noelzubin/redis-go/store"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetScoreMemberPairs(t *testing.T) {
	assert := assert.New(t)

	res, err := GetScoreMemberPairs([]string{"1", "one"})
	assert.Nil(err)
	assert.Equal(store.NewScoreMember(int64(1), "one"), res[0])
}

func Test_GetScoreMemberPairsMultiple(t *testing.T) {
	assert := assert.New(t)

	res, err := GetScoreMemberPairs([]string{"1", "one", "2", "two"})
	assert.Nil(err)
	assert.Equal(store.NewScoreMember(int64(1), "one"), res[0])
	assert.Equal(store.NewScoreMember(int64(2), "two"), res[1])
}

func Test_GetScoreMemberPairs_InvalidInputLength(t *testing.T) {
	assert := assert.New(t)

	assert.Panics(func() { GetScoreMemberPairs([]string{"1"}) })
}

func Test_GetScoreMemberPairs_InvalidScore(t *testing.T) {
	assert := assert.New(t)

	_, err := GetScoreMemberPairs([]string{"NaN", "what"})
	assert.NotNil(err)
}
