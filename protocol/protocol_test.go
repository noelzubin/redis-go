package protocol

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeSimpleString(t *testing.T) {
	assert := assert.New(t)
	value, err := DecodeRESP(bufio.NewReader(bytes.NewBufferString("+foo\r\n")))

	assert.Nil(err)
	assert.Equal(value.typ, SimpleString)
	assert.Equal(value.String(), "foo")
}

func TestDecodeBulkString(t *testing.T) {
	assert := assert.New(t)
	value, err := DecodeRESP(bufio.NewReader(bytes.NewBufferString("$4\r\nabcd\r\n")))

	assert.Nil(err)
	assert.Equal(value.typ, BulkString)
	assert.Equal(value.String(), "abcd")
}

func TestDecodeBulkStringArray(t *testing.T) {
	assert := assert.New(t)

	value, err := DecodeRESP(bufio.NewReader(bytes.NewBufferString("*2\r\n$3\r\nGET\r\n$4\r\nthis\r\n")))

	assert.Nil(err)
	assert.Equal(value.typ, Array)
	assert.Equal(value.Array()[0].String(), "GET")
	assert.Equal(value.Array()[1].String(), "this")
}

func TestDecodeNil(t *testing.T) {
	assert := assert.New(t)

	value, err := DecodeRESP(bufio.NewReader(bytes.NewBufferString("_\r\n")))

	assert.Nil(err)
	assert.Equal(value.typ, Nil)
}

func TestDecodeInteger(t *testing.T) {
	assert := assert.New(t)

	value, err := DecodeRESP(bufio.NewReader(bytes.NewBufferString(":123\r\n")))

	assert.Nil(err)
	assert.Equal(value.typ, Integer)
	assert.Equal(value.Integer(), int64(123))
}

func TestEncodeString(t *testing.T) {
	assert := assert.New(t)

	hello := "hello"
	value := NewSimpleStringValue(&hello)

	assert.Equal(value.Encode(), []byte("+hello\r\n"))
}

func TestEncodeNil(t *testing.T) {
	assert := assert.New(t)

	value := NewNilValue()

	assert.Equal(value.Encode(), []byte("_\r\n"))
}

func TestEncodeInteger(t *testing.T) {
	assert := assert.New(t)

	value := NewSimpleIntValue(123)

	assert.Equal(value.Encode(), []byte(":123\r\n"))
}

func TestEncodeArrayOfStrings(t *testing.T) {
	assert := assert.New(t)

	value := NewArrayStringValue([]string{"GET", "THIS"})

	assert.Equal(value.Encode(), []byte("*2\r\n+GET\r\n+THIS\r\n"))
}
