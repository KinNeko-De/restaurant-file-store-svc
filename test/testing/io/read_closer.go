package io

import (
	"io"
	"testing"

	"github.com/stretchr/testify/mock"
)

func CreateReadCloser(t *testing.T, readFile []byte) *MockReadCloser {
	fileReader := &MockReadCloser{}
	readIndex := 0
	fileReader.EXPECT().Read(mock.Anything).RunAndReturn(func(data []byte) (int, error) {
		if readIndex >= len(readFile) {
			return 0, io.EOF
		}

		remainingBytes := len(readFile) - readIndex
		bytesToRead := min(len(data), remainingBytes)
		endbytes := readIndex + bytesToRead
		copy(data, readFile[readIndex:endbytes])
		readIndex = endbytes
		return bytesToRead, nil
	})
	fileReader.EXPECT().Close().Return(nil).Times(1)
	return fileReader
}
