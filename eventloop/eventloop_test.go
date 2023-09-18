package eventloop

import (
	rwMock "noelzubin/redis-go/eventloop/mocks"
	"noelzubin/redis-go/store"
	storeMock "noelzubin/redis-go/store/mocks"
	"testing"
	"time"
)

var st storeMock.Store
var el Eventloop

func setup() {
	st = storeMock.Store{}
	el = *InitEventloop(&st)
	go el.RunLoop()
}

func Test_Set(t *testing.T) {
	setup()
	go el.RunLoop()
	var tm *time.Time = nil
	st.On("Set", "foo", "bar", tm).Return()
	el.HandleConnection(
		rwMock.NewMockReadWriteCloser("*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"),
	)
	st.AssertNumberOfCalls(t, "Set", 1)
}

func Test_Get(t *testing.T) {
	setup()
	bar := "bar"
	st.On("Get", "foo").Return(&bar)
	el.HandleConnection(
		rwMock.NewMockReadWriteCloser("*2\r\n$3\r\nGET\r\n$3\r\nfoo\r\n"),
	)
	st.AssertNumberOfCalls(t, "Get", 1)
}

func Test_Expire(t *testing.T) {
	setup()
	var tm *time.Time = nil
	st.On("Set", "foo", "bar", tm).Return()
	el.HandleConnection(
		rwMock.NewMockReadWriteCloser("*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"),
	)
	st.AssertNumberOfCalls(t, "Set", 1)
}

func Test_Ping(t *testing.T) {
	setup()
	PONG := "PONG"
	st.On("Ping").Return(&PONG)
	el.HandleConnection(
		rwMock.NewMockReadWriteCloser("*1\r\n$4\r\nPING"),
	)
	st.AssertNumberOfCalls(t, "Ping", 1)
}

func Test_Del(t *testing.T) {
	setup()
	st.On("Del", "foo").Return(1)
	el.HandleConnection(
		rwMock.NewMockReadWriteCloser("*2\r\n$3\r\nDEL\r\n$3\r\nfoo\r\n"),
	)
	st.AssertNumberOfCalls(t, "Del", 1)
}

func Test_Zadd(t *testing.T) {
	setup()
	members := []store.ScoreMember{store.NewScoreMember(4, "bar")}
	st.On("ZAdd", "foo", members).Return(1)
	el.HandleConnection(
		rwMock.NewMockReadWriteCloser("*4\r\n$4\r\nZADD\r\n$3\r\nfoo\r\n$1\r\n4\r\n$3\r\nbar\r\n"),
	)
	st.AssertNumberOfCalls(t, "ZAdd", 1)
}

func Test_Zrange(t *testing.T) {
	go el.RunLoop()
	st.On("ZRange", "foo", 1, -1, false).Return([]string{})
	el.HandleConnection(
		rwMock.NewMockReadWriteCloser("*4\r\n$6\r\nZRANGE\r\n$3\r\nfoo\r\n$1\r\n1\r\n$2\r\n-1\r\n"),
	)
	st.AssertNumberOfCalls(t, "ZRange", 1)
}

func Test_Zrange_WithScores(t *testing.T) {
	setup()
	st.On("ZRange", "foo", 1, -1, true).Return([]string{})
	el.HandleConnection(
		rwMock.NewMockReadWriteCloser("*5\r\n$6\r\nZRANGE\r\n$3\r\nfoo\r\n$1\r\n1\r\n$2\r\n-1\r\n$10\r\nWITHSCORES\r\n"),
	)
	st.AssertNumberOfCalls(t, "ZRange", 1)
}

func Test_CleanUp(t *testing.T) {
	setup()
	st.On("CleanUp").Return()
	go el.StartCleanUpTimer()
	<-time.NewTimer(500 * time.Millisecond).C
	st.AssertNumberOfCalls(t, "CleanUp", 5)
}
