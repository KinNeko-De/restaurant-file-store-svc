package io

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

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

func CreateWriterCloser2(t *testing.T, writtenChunks [][]byte) *MockWriteCloser {
	fileWriter := NewMockWriteCloser(t)
	for _, chunk := range writtenChunks {
		fileWriter.EXPECT().Write(chunk).Return(len(chunk), nil).Times(1)
	}
	fileWriter.EXPECT().Close().Return(nil).Times(1)
	return fileWriter
}

func CreateWriterCloserRanIntoWriteError2(t *testing.T, successfulWrittenChunks [][]byte, errorAfterSuccessfulWrites error) *MockWriteCloser {
	fileWriter := NewMockWriteCloser(t)
	for _, chunk := range successfulWrittenChunks {
		fileWriter.EXPECT().Write(chunk).Return(len(chunk), nil).Times(1)
	}
	fileWriter.EXPECT().Write(mock.Anything).Return(0, errorAfterSuccessfulWrites).Times(1)

	fileWriter.EXPECT().Close().Return(nil).Times(1)
	return fileWriter
}

func CreateWriterCloserRanIntoCloseError2(t *testing.T, writtenChunks [][]byte, err error) *MockWriteCloser {
	fileWriter := NewMockWriteCloser(t)
	for _, chunk := range writtenChunks {
		fileWriter.EXPECT().Write(chunk).Return(len(chunk), nil).Times(1)
	}

	fileWriter.EXPECT().Close().Return(err).Times(1)
	return fileWriter
}
