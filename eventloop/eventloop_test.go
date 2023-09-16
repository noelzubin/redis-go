package eventloop

import (
	"noelzubin/redis-go/mocks"
	"noelzubin/redis-go/store"
	"testing"
	"time"
)

var st = &mocks.Store{}
var el Eventloop = *InitEventloop(st)

func Test_Set(t *testing.T) {
	go el.RunLoop()
	var tm *time.Time = nil
	st.On("Set", "foo", "bar", tm).Return()
	el.HandleConnection(
		mocks.NewMockReadWriteCloser("*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"),
	)
	st.AssertNumberOfCalls(t, "Set", 1)
}

func Test_Get(t *testing.T) {
	// assert := assert.New(t)
	go el.RunLoop()
	bar := "bar"
	st.On("Get", "foo").Return(&bar)
	el.HandleConnection(
		mocks.NewMockReadWriteCloser("*2\r\n$3\r\nGET\r\n$3\r\nfoo\r\n"),
	)
	st.AssertNumberOfCalls(t, "Get", 1)
}

func Test_Expire(t *testing.T) {
	go el.RunLoop()
	var tm *time.Time = nil
	st.On("Set", "foo", "bar", tm).Return()
	el.HandleConnection(
		mocks.NewMockReadWriteCloser("*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"),
	)
	st.AssertNumberOfCalls(t, "Set", 1)
	st.AssertCalled(t, "Set", "foo", "bar", tm)
}

func Test_Ping(t *testing.T) {
	go el.RunLoop()
	PONG := "PONG"
	st.On("Ping").Return(&PONG)
	el.HandleConnection(
		mocks.NewMockReadWriteCloser("*1\r\n$4\r\nPING"),
	)
	st.AssertNumberOfCalls(t, "Ping", 1)
}

func Test_Del(t *testing.T) {
	go el.RunLoop()
	st.On("Del", "foo").Return(1)
	el.HandleConnection(
		mocks.NewMockReadWriteCloser("*2\r\n$3\r\nDEL\r\n$3\r\nfoo\r\n"),
	)
	st.AssertNumberOfCalls(t, "Del", 1)
}

func Test_Zadd(t *testing.T) {
	go el.RunLoop()
	members := []store.ScoreMember{store.NewScoreMember(4, "bar")}
	st.On("ZAdd", "foo", members).Return(1)
	el.HandleConnection(
		mocks.NewMockReadWriteCloser("*4\r\n$4\r\nZADD\r\n$3\r\nfoo\r\n$1\r\n4\r\n$3\r\nbar\r\n"),
	)
	st.AssertNumberOfCalls(t, "ZAdd", 1)
}

func Test_Zrange(t *testing.T) {
	go el.RunLoop()
	st.On("ZRange", "foo", 1, -1, false).Return([]string{})
	el.HandleConnection(
		mocks.NewMockReadWriteCloser("*4\r\n$6\r\nZRANGE\r\n$3\r\nfoo\r\n$1\r\n1\r\n$2\r\n-1\r\n"),
	)
	st.AssertNumberOfCalls(t, "ZRange", 1)
}

func Test_Zrange_WithScores(t *testing.T) {
	go el.RunLoop()
	st.On("ZRange", "foo", 1, -1, true).Return([]string{})
	el.HandleConnection(
		mocks.NewMockReadWriteCloser("*5\r\n$6\r\nZRANGE\r\n$3\r\nfoo\r\n$1\r\n1\r\n$2\r\n-1\r\n$10\r\nWITHSCORES\r\n"),
	)
	st.AssertNumberOfCalls(t, "ZRange", 1)
}
