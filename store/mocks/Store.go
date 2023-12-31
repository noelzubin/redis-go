// Code generated by mockery v2.33.3. DO NOT EDIT.

package mocks

import (
	store "noelzubin/redis-go/store"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// Store is an autogenerated mock type for the Store type
type Store struct {
	mock.Mock
}

// CleanUp provides a mock function with given fields:
func (_m *Store) CleanUp() {
	_m.Called()
}

// Del provides a mock function with given fields: keys
func (_m *Store) Del(keys ...string) int {
	_va := make([]interface{}, len(keys))
	for _i := range keys {
		_va[_i] = keys[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 int
	if rf, ok := ret.Get(0).(func(...string) int); ok {
		r0 = rf(keys...)
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// Expire provides a mock function with given fields: k, seconds
func (_m *Store) Expire(k string, seconds int) int {
	ret := _m.Called(k, seconds)

	var r0 int
	if rf, ok := ret.Get(0).(func(string, int) int); ok {
		r0 = rf(k, seconds)
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// Get provides a mock function with given fields: k
func (_m *Store) Get(k string) *string {
	ret := _m.Called(k)

	var r0 *string
	if rf, ok := ret.Get(0).(func(string) *string); ok {
		r0 = rf(k)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*string)
		}
	}

	return r0
}

// Keys provides a mock function with given fields: k
func (_m *Store) Keys(k string) []string {
	ret := _m.Called(k)

	var r0 []string
	if rf, ok := ret.Get(0).(func(string) []string); ok {
		r0 = rf(k)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// Ping provides a mock function with given fields:
func (_m *Store) Ping() *string {
	ret := _m.Called()

	var r0 *string
	if rf, ok := ret.Get(0).(func() *string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*string)
		}
	}

	return r0
}

// Set provides a mock function with given fields: k, v, e
func (_m *Store) Set(k string, v string, e *time.Time) {
	_m.Called(k, v, e)
}

// ZAdd provides a mock function with given fields: k, s
func (_m *Store) ZAdd(k string, s []store.ScoreMember) int {
	ret := _m.Called(k, s)

	var r0 int
	if rf, ok := ret.Get(0).(func(string, []store.ScoreMember) int); ok {
		r0 = rf(k, s)
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// ZRange provides a mock function with given fields: k, start, stop, withScores
func (_m *Store) ZRange(k string, start int, stop int, withScores bool) []string {
	ret := _m.Called(k, start, stop, withScores)

	var r0 []string
	if rf, ok := ret.Get(0).(func(string, int, int, bool) []string); ok {
		r0 = rf(k, start, stop, withScores)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// NewStore creates a new instance of Store. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *Store {
	mock := &Store{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
