package mocks

import (
	"fmt"
	"io"
)

// MockReadWriteCloser is a custom type that implements io.ReadWriteCloser.
type MockReadWriteCloser struct {
	ReadVal []byte
	Done    bool
}

func NewMockReadWriteCloser(read string) *MockReadWriteCloser {
	return &MockReadWriteCloser{
		ReadVal: []byte(read),
		Done:    false,
	}
}

// Read calls the ReadFunc of the MockReadWriteCloser.
func (m *MockReadWriteCloser) Read(p []byte) (int, error) {
	if m.Done {
		return 0, io.EOF
	}
	copy(p, []byte(m.ReadVal))
	m.Done = true
	return len(p), nil
}

// Write calls the WriteFunc of the MockReadWriteCloser.
func (m *MockReadWriteCloser) Write(p []byte) (int, error) {
	return len(p), nil
}

// Close calls the CloseFunc of the MockReadWriteCloser.
func (m *MockReadWriteCloser) Close() error {
	fmt.Println("closed")
	return nil
}
