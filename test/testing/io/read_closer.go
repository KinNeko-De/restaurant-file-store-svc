package io

import (
	"io"
	"testing"

	"github.com/stretchr/testify/mock"
)

func CreateReadCloser(t *testing.T, readFile []byte) *MockReadCloser {
	fileReader := &MockReadCloser{}
	setupReadOfFile(fileReader, readFile)
	fileReader.EXPECT().Close().Return(nil).Times(1)
	return fileReader
}

func setupReadOfFile(fileReader *MockReadCloser, readFile []byte) {
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
}

func CreateReadCloserRanIntoReadError(t *testing.T, errorAfterSuccessfulRead error) *MockReadCloser {
	fileReader := &MockReadCloser{}
	fileReader.EXPECT().Read(mock.Anything).Return(0, errorAfterSuccessfulRead).Times(1)
	return fileReader
}

func CreateReadCloserRanIntoCloseError(t *testing.T, readFile []byte, closeErr error) *MockReadCloser {
	fileReader := &MockReadCloser{}
	setupReadOfFile(fileReader, readFile)
	fileReader.EXPECT().Close().Return(closeErr).Times(1)
	return fileReader
}
