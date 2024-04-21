package io

import (
	"io"
	"testing"

	"github.com/stretchr/testify/mock"
)

func (fileReader *MockReadCloser) SetupSuccessfulRead(t *testing.T, readFile []byte) {
	fileReader.SetupRead(t, readFile)
	fileReader.SetupClose(t)
}

func (fileReader *MockReadCloser) SetupRead(t *testing.T, readFile []byte) {
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

func (fileReader *MockReadCloser) SetupReadError(t *testing.T, readError error) {
	fileReader.EXPECT().Read(mock.Anything).Return(0, readError).Times(1)
}

func (fileReader *MockReadCloser) SetupClose(t *testing.T) {
	fileReader.EXPECT().Close().Return(nil).Times(1)
}

func (fileReader *MockReadCloser) SetupCloseError(t *testing.T, closeErr error) {
	fileReader.EXPECT().Close().Return(closeErr).Times(1)
}
