package protocol

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Type represents a Value type
type Type byte

const (
	SimpleString Type = '+'
	BulkString   Type = '$'
	Array        Type = '*'
	Integer      Type = ':'
	Nil          Type = '_'
	Error        Type = '-'
)

// Value represents the data of a valid RESP type.
type Value struct {
	typ    Type
	bytes  []byte
	array  []Value
	intVal int64
}

// String converts Value to a string.
//
// If Value cannot be converted, an empty string is returned.
func (v Value) String() string {
	if v.typ == BulkString || v.typ == SimpleString {
		return string(v.bytes)
	}

	return ""
}

// Integer converts Value to an int64
//
// If Value cannot be converted 1 is returned
func (v Value) Integer() int64 {
	if v.typ == Integer {
		val, err := strconv.Atoi(string(v.bytes))
		if err != nil {
			return 0
		}
		return int64(val)
	}

	return 1
}

// Array converts Value to an array.
//
// If Value cannot be converted, an empty array is returned.
func (v Value) Array() []Value {
	if v.typ == Array {
		return v.array
	}
	return []Value{}
}

// Output converts Value to a string for output.
func (v Value) Output() string {
	s := make([]string, 0)

	if v.typ == Array {
		for _, v := range v.array {
			s = append(s, v.String())
		}
	} else if v.typ == SimpleString || v.typ == BulkString {
		s = append(s, v.String())
	} else if v.typ == Integer {
		s = append(s, strconv.Itoa(int(v.Integer())))
	} else {
		s = append(s, "(nil)")
	}

	return strings.Join(s, "\n")
}

// DecodeRESP parses a RESP message and returns a Value
func DecodeRESP(byteStream *bufio.Reader) (Value, error) {
	dataTypeByte, err := byteStream.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch string(dataTypeByte) {
	case "+":
		return decodeSimpleString(byteStream)
	case "$":
		return decodeBulkString(byteStream)
	case "*":
		return decodeArray(byteStream)
	case ":":
		return decodeInteger(byteStream)
	case "-":
		return decodeError(byteStream)
	case "_":
		return decodeNil(byteStream)
	}
	return Value{}, fmt.Errorf("invalid RESP data type byte: %s", string(dataTypeByte))
}

func decodeSimpleString(byteStream *bufio.Reader) (Value, error) {
	readBytes, err := readUntilCRLF(byteStream)
	if err != nil {
		return Value{}, err
	}

	return Value{
		typ:   SimpleString,
		bytes: readBytes,
	}, nil
}

func decodeInteger(byteStream *bufio.Reader) (Value, error) {
	readBytes, err := readUntilCRLF(byteStream)
	if err != nil {
		return Value{}, err
	}

	return Value{
		typ:   Integer,
		bytes: readBytes,
	}, nil
}

func decodeError(byteStream *bufio.Reader) (Value, error) {
	readBytes, err := readUntilCRLF(byteStream)
	if err != nil {
		return Value{}, err
	}

	return Value{
		typ:   SimpleString,
		bytes: readBytes,
	}, nil
}

func decodeNil(byteStream *bufio.Reader) (Value, error) {
	readBytes, err := readUntilCRLF(byteStream)
	if err != nil {
		return Value{}, err
	}

	return Value{
		typ:   Nil,
		bytes: readBytes,
	}, nil
}

func decodeBulkString(byteStream *bufio.Reader) (Value, error) {
	readBytesForCount, err := readUntilCRLF(byteStream)
	if err != nil {
		return Value{}, fmt.Errorf("failed to read bulk string length: %s", err)
	}

	count, err := strconv.Atoi(string(readBytesForCount))
	if err != nil {
		return Value{}, fmt.Errorf("failed to parse bulk string length: %s", err)
	}

	readBytes := make([]byte, count+2)

	if _, err := io.ReadFull(byteStream, readBytes); err != nil {
		return Value{}, fmt.Errorf("failed to read bulk string contents: %s", err)
	}

	return Value{
		typ:   BulkString,
		bytes: readBytes[:count],
	}, nil
}

func decodeArray(byteStream *bufio.Reader) (Value, error) {
	readBytesForCount, err := readUntilCRLF(byteStream)
	if err != nil {
		return Value{}, fmt.Errorf("failed to read bulk string length: %s", err)
	}

	count, err := strconv.Atoi(string(readBytesForCount))
	if err != nil {
		return Value{}, fmt.Errorf("failed to parse bulk string length: %s", err)
	}

	array := []Value{}

	for i := 1; i <= count; i++ {
		value, err := DecodeRESP(byteStream)
		if err != nil {
			return Value{}, err
		}

		array = append(array, value)
	}

	return Value{
		typ:   Array,
		array: array,
	}, nil

}

func readUntilCRLF(byteStream *bufio.Reader) ([]byte, error) {
	readBytes := []byte{}

	for {
		b, err := byteStream.ReadBytes('\n')
		if err != nil {
			return nil, err
		}

		readBytes = append(readBytes, b...)
		if len(readBytes) >= 2 && readBytes[len(readBytes)-2] == '\r' {
			break
		}
	}

	return readBytes[:len(readBytes)-2], nil
}

// Encode encodes a value into a RESP Message
func (v *Value) Encode() []byte {
	buf := new(bytes.Buffer)
	switch v.typ {
	case BulkString:
	case SimpleString:
		buf.Write([]byte("+" + string(v.bytes) + "\r\n"))
	case Array:
		buf.Write([]byte("*"))
		buf.Write([]byte(strconv.Itoa(len(v.array))))
		buf.Write([]byte("\r\n"))
		for _, v := range v.array {
			buf.Write(v.Encode())
		}
	case Integer:
		buf.Write([]byte(":" + strconv.Itoa(int(v.intVal)) + "\r\n"))
	case Error:
		buf.Write([]byte("-" + string(v.bytes) + "\r\n"))
	case Nil:
		buf.Write([]byte("_\r\n"))
	}

	return buf.Bytes()
}

// NewErrorValue creates a new Error Value
func NewErrorValue(err string) Value {
	return Value{
		typ:   '-',
		bytes: []byte(err),
	}
}

// NewSimpleStringValue creates a new String Value
func NewSimpleStringValue(str *string) Value {
	return Value{
		typ:   '+',
		bytes: []byte(*str),
	}
}

// NewSimpleIntValue creates a new Integer Value
func NewSimpleIntValue(val int64) Value {
	return Value{
		typ:    ':',
		intVal: val,
	}
}

// NewNilValue creates a new Nil Value
func NewNilValue() Value {
	return Value{
		typ: '_',
	}
}

// New ArrayStringValue creates a new Array Value
func NewArrayStringValue(arr []string) Value {
	vals := make([]Value, 0)

	for _, a := range arr {
		vals = append(vals, NewSimpleStringValue(&a))
	}

	return Value{
		typ:   '*',
		array: vals,
	}
}
