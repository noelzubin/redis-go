package store

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_CreateStore_Success(t *testing.T) {
	InitStore()
}

func Test_Set_Key(t *testing.T) {
	assert := assert.New(t)
	s := InitStore()
	assert.NotPanics(func() {
		s.Set("foo", "bar", nil)
	})
}

func Test_Get_Key_Success(t *testing.T) {
	assert := assert.New(t)
	s := InitStore()
	s.Set("foo", "bar", nil)
	res := s.Get("foo")
	bar := "bar"
	assert.Equal(&bar, res)
}

func Test_Get_Key_Missing(t *testing.T) {
	assert := assert.New(t)
	s := InitStore()
	s.Set("foo", "bar", nil)
	res := s.Get("five")
	assert.Nil(res)
}

func Test_Get_Key_Expired(t *testing.T) {
	assert := assert.New(t)
	s := InitStore()
	now := time.Now().Add(-1 * time.Second)
	s.Set("foo", "bar", &now)
	res := s.Get("five")
	assert.Nil(res)
}

func Test_Del_Key(t *testing.T) {
	assert := assert.New(t)
	s := InitStore()
	s.Set("foo", "bar", nil)
	res := s.Get("foo")
	assert.NotNil(res)
	count := s.Del("foo")
	res = s.Get("foo")
	assert.Nil(res)
	assert.Equal(1, count)
}

func Test_Del_Multiple_Keys(t *testing.T) {
	assert := assert.New(t)
	s := InitStore()
	s.Set("foo", "bar", nil)
	s.Set("uno", "one", nil)
	s.Set("dos", "two", nil)

	count := s.Del("foo", "uno", "dos", "tres")
	assert.Equal(3, count)
}

func Test_Expire_Key(t *testing.T) {
	assert := assert.New(t)
	s := InitStore()
	s.Set("foo", "bar", nil)
	assert.NotNil(s.Get("foo"))
	s.Expire("foo", 1)
	assert.NotNil(s.Get("foo"))
	time.Sleep(2 * time.Second)
	assert.Nil(s.Get("foo"))
}

func Test_Keys(t *testing.T) {
	assert := assert.New(t)
	s := InitStore()
	s.Set("foo", "bar", nil)
	s.Set("uno", "bar", nil)
	s.Set("dos", "two", nil)

	keys := s.Keys("*")
	assert.ElementsMatch([]string{"foo", "uno", "dos"}, keys)
}

func Test_ZAdd(t *testing.T) {
	assert := assert.New(t)
	s := InitStore()
	assert.NotPanics(func() {
		scoreMembers := []ScoreMember{
			{3, "three"},
			{1, "one"},
			{2, "two"},
		}

		s.ZAdd("foo", scoreMembers)
	})
}

func Test_ZRange_All(t *testing.T) {
	assert := assert.New(t)
	s := InitStore()
	scoreMembers := []ScoreMember{
		{3, "three"},
		{1, "one"},
		{2, "two"},
	}

	s.ZAdd("foo", scoreMembers)

	assert.Equal([]string{"one", "two", "three"}, s.ZRange("foo", 1, -1, false))
}

func Test_ZRange_Multiple_Assign(t *testing.T) {
	assert := assert.New(t)
	s := InitStore()
	scoreMembers := []ScoreMember{
		{3, "three"},
		{1, "one"},
	}

	s.ZAdd("foo", scoreMembers)

	scoreMembers = []ScoreMember{
		{5, "five"},
		{4, "four"},
	}

	s.ZAdd("foo", scoreMembers)

	assert.Equal([]string{"one", "three", "four", "five"}, s.ZRange("foo", 1, -1, false))
}

func Test_ZRange_WithScores(t *testing.T) {
	assert := assert.New(t)
	s := InitStore()
	scoreMembers := []ScoreMember{
		{3, "three"},
		{1, "one"},
		{2, "two"},
	}

	s.ZAdd("foo", scoreMembers)

	assert.Equal([]string{"one", "1", "two", "2", "three", "3"}, s.ZRange("foo", 1, -1, true))
}

func Test_Value_IsExpired(t *testing.T) {
	assert := assert.New(t)
	v1 := Value{value: 1, expiry: nil}
	assert.Equal(false, v1.isExpired())

	afterASecond := time.Now().Add(1 * time.Second)
	v2 := Value{value: 1, expiry: &afterASecond}
	assert.Equal(false, v2.isExpired())

	beforeASecond := time.Now().Add(-1 * time.Second)
	v3 := Value{value: 1, expiry: &beforeASecond}
	assert.Equal(true, v3.isExpired())
}
