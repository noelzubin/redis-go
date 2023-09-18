package store

import (
	"fmt"
	"testing"
	"time"

	"noelzubin/redis-go/set/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var expireSet *mocks.IStringSet

func setup() {
	fmt.Println("SETTING UP TEST")
	expireSet = &mocks.IStringSet{}
	expireSet.On("Add", mock.AnythingOfType("string")).Return()
	expireSet.On("Remove", mock.AnythingOfType("string")).Return()
}

func Test_CreateStore_Success(t *testing.T) {
	setup()
	mockStringSet := &mocks.IStringSet{}
	InitStore(mockStringSet)
}

func Test_Set_Key(t *testing.T) {
	setup()
	assert := assert.New(t)
	s := InitStore(expireSet)
	assert.NotPanics(func() {
		s.Set("foo", "bar", nil)
	})
	expireSet.AssertNumberOfCalls(t, "Add", 0)
}

func Test_Set_Key_WithExpiry(t *testing.T) {
	setup()
	assert := assert.New(t)
	s := InitStore(expireSet)
	assert.NotPanics(func() {
		setup()
		now := time.Now()
		s.Set("foo", "bar", &now)
	})
	expireSet.AssertNumberOfCalls(t, "Add", 0)
}

func Test_Get_Key_Success(t *testing.T) {
	setup()
	assert := assert.New(t)
	s := InitStore(expireSet)
	s.Set("foo", "bar", nil)
	res := s.Get("foo")
	bar := "bar"
	assert.Equal(&bar, res)
	expireSet.AssertNumberOfCalls(t, "Add", 0)
	expireSet.AssertNumberOfCalls(t, "Remove", 0)
}

func Test_Get_Key_Missing(t *testing.T) {
	setup()
	assert := assert.New(t)
	s := InitStore(expireSet)
	s.Set("foo", "bar", nil)
	res := s.Get("five")
	assert.Nil(res)
	expireSet.AssertNumberOfCalls(t, "Add", 0)
	expireSet.AssertNumberOfCalls(t, "Remove", 0)
}

func Test_Get_Key_Expired(t *testing.T) {
	setup()
	assert := assert.New(t)
	s := InitStore(expireSet)
	now := time.Now().Add(-1 * time.Second)
	s.Set("foo", "bar", &now)
	res := s.Get("foo")
	assert.Nil(res)
	expireSet.AssertNumberOfCalls(t, "Add", 1)
	expireSet.AssertNumberOfCalls(t, "Remove", 1)
}

func Test_Del_Key(t *testing.T) {
	setup()
	assert := assert.New(t)
	s := InitStore(expireSet)
	s.Set("foo", "bar", nil)
	res := s.Get("foo")
	assert.NotNil(res)
	count := s.Del("foo")
	res = s.Get("foo")
	assert.Nil(res)
	assert.Equal(1, count)
	expireSet.AssertNumberOfCalls(t, "Remove", 1)
}

func Test_Del_Multiple_Keys(t *testing.T) {
	setup()
	assert := assert.New(t)
	s := InitStore(expireSet)
	s.Set("foo", "bar", nil)
	s.Set("uno", "one", nil)
	s.Set("dos", "two", nil)

	count := s.Del("foo", "uno", "dos", "tres")
	assert.Equal(3, count)
	expireSet.AssertNumberOfCalls(t, "Remove", 3)
}

func Test_Expire_Key(t *testing.T) {
	setup()
	assert := assert.New(t)
	s := InitStore(expireSet)
	s.Set("foo", "bar", nil)
	assert.NotNil(s.Get("foo"))
	s.Expire("foo", 1)
	assert.NotNil(s.Get("foo"))
	time.Sleep(2 * time.Second)
	assert.Nil(s.Get("foo"))

	expireSet.AssertNumberOfCalls(t, "Add", 1)
	expireSet.AssertNumberOfCalls(t, "Remove", 1)
}

func Test_Keys(t *testing.T) {
	setup()
	assert := assert.New(t)
	s := InitStore(expireSet)
	s.Set("foo", "bar", nil)
	s.Set("uno", "bar", nil)
	s.Set("dos", "two", nil)

	keys := s.Keys("*")
	assert.ElementsMatch([]string{"foo", "uno", "dos"}, keys)
}

func Test_ZAdd(t *testing.T) {
	setup()
	assert := assert.New(t)
	s := InitStore(expireSet)
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
	setup()
	assert := assert.New(t)
	s := InitStore(expireSet)
	scoreMembers := []ScoreMember{
		{3, "three"},
		{1, "one"},
		{2, "two"},
	}

	s.ZAdd("foo", scoreMembers)

	assert.Equal([]string{"one", "two", "three"}, s.ZRange("foo", 1, -1, false))
}

func Test_ZRange_Multiple_Assign(t *testing.T) {
	setup()
	assert := assert.New(t)
	s := InitStore(expireSet)
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
	setup()
	assert := assert.New(t)
	s := InitStore(expireSet)
	scoreMembers := []ScoreMember{
		{3, "three"},
		{1, "one"},
		{2, "two"},
	}

	s.ZAdd("foo", scoreMembers)

	assert.Equal([]string{"one", "1", "two", "2", "three", "3"}, s.ZRange("foo", 1, -1, true))
}

func Test_Value_IsExpired(t *testing.T) {
	setup()
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

func Test_CleanUp(t *testing.T) {
	setup()
	assert := assert.New(t)
	s := InitStore(expireSet)

	s.Set("one", "two", nil)
	s.Set("foo", "bar", nil)

	resp := make([]string, 0)
	expireSet.On("RandomN", mock.AnythingOfType("int")).Return(resp)

	s.CleanUp()
	expireSet.AssertNumberOfCalls(t, "RandomN", 1)

	bar := "bar"
	assert.Equal(s.Get("foo"), &bar)
}

func Test_CleanUp_Expired_Value(t *testing.T) {
	setup()
	s := InitStore(expireSet)

	now := time.Now().Add(-1 * time.Second)
	s.Set("one", "two", nil)
	s.Set("foo", "bar", &now)
	expireSet.AssertNumberOfCalls(t, "Add", 1)

	// return 1 value with expired key
	resp := []string{"foo"}
	expireSet.On("RandomN", mock.AnythingOfType("int")).Return(resp)

	s.CleanUp()
	expireSet.AssertNumberOfCalls(t, "Remove", 1)

	// No more extra Remove calls
	s.Get("foo")
	expireSet.AssertNumberOfCalls(t, "Remove", 1)
}

func Test_CleanUp_All_Expired_Values(t *testing.T) {
	setup()
	s := InitStore(expireSet)

	now := time.Now().Add(-1 * time.Second)
	s.Set("one", "two", &now)
	s.Set("foo", "bar", &now)
	expireSet.AssertNumberOfCalls(t, "Add", 2)

	// return 1 value with expired key
	resp := []string{"foo", "one"}
	expireSet.On("RandomN", mock.AnythingOfType("int")).Return(resp)

	s.CleanUp()
	expireSet.AssertNumberOfCalls(t, "Remove", 2)

	s.Get("foo")
	s.Get("one")
	expireSet.AssertNumberOfCalls(t, "Remove", 2)
}

func Test_CleanUp_Call_Multiple_Times_If_More_Than_5(t *testing.T) {
	setup()
	s := InitStore(expireSet)

	now := time.Now().Add(-1 * time.Second)
	later := time.Now().Add(20 * time.Second)
	s.Set("one", "one", &now)
	s.Set("two", "two", &now)
	s.Set("three", "two", &now)
	s.Set("four", "two", &now)
	s.Set("five", "two", &now)
	s.Set("six", "two", &now)
	s.Set("seven", "two", &later)
	s.Set("eight", "two", &later)

	// 6 expired keys
	resp := []string{"one", "two", "three", "four", "five", "six"}
	expireSet.On("RandomN", mock.AnythingOfType("int")).Return(resp).Once()
	resp2 := []string{"seven", "eight"}
	expireSet.On("RandomN", mock.AnythingOfType("int")).Return(resp2).Once()

	s.CleanUp()
	expireSet.AssertNumberOfCalls(t, "Remove", 6)
	expireSet.AssertNumberOfCalls(t, "RandomN", 2)
}

func Test_CleanUp_Call_Multiple_Times_If_Less_Than_5(t *testing.T) {
	setup()
	s := InitStore(expireSet)

	now := time.Now().Add(-1 * time.Second)
	later := time.Now().Add(20 * time.Second)
	s.Set("one", "one", &now)
	s.Set("two", "two", &now)
	s.Set("three", "two", &now)
	s.Set("four", "two", &now)
	s.Set("five", "two", &now)
	s.Set("six", "two", &now)
	s.Set("seven", "two", &later)
	s.Set("eight", "two", &later)

	// only 4 have exired
	resp := []string{"one", "two", "three", "seven", "eight", "six"}
	expireSet.On("RandomN", mock.AnythingOfType("int")).Return(resp).Once()
	resp2 := []string{"four", "five"}
	expireSet.On("RandomN", mock.AnythingOfType("int")).Return(resp2).Once()

	s.CleanUp()
	expireSet.AssertNumberOfCalls(t, "Remove", 4)
	expireSet.AssertNumberOfCalls(t, "RandomN", 1)
}
