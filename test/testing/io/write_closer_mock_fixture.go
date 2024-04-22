package io

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

func (mockWriteCloser *MockWriteCloser) SetupSuccessfulWrite(t *testing.T, writtenChunks [][]byte) {
	mockWriteCloser.SetupWrite(t, writtenChunks)
	mockWriteCloser.SetupClose(t)
}

func (mockWriteCloser *MockWriteCloser) SetupWrite(t *testing.T, writtenChunks [][]byte) {
	for _, chunk := range writtenChunks {
		mockWriteCloser.EXPECT().Write(chunk).Return(len(chunk), nil).Times(1)
	}
}

func (mockWriteCloser *MockWriteCloser) SetupWriteError(t *testing.T, writeError error) {
	mockWriteCloser.EXPECT().Write(mock.Anything).Return(0, writeError).Times(1)
}

func (mockWriteCloser *MockWriteCloser) SetupClose(t *testing.T) {
	mockWriteCloser.EXPECT().Close().Return(nil).Times(1)
}

func (mockWriteCloser *MockWriteCloser) SetupCloseError(t *testing.T, writeError error) {
	mockWriteCloser.EXPECT().Close().Return(writeError).Times(1)
}
