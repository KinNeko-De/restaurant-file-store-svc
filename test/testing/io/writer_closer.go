package io

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

func CreateWriterCloser(t *testing.T, writtenChunks [][]byte) *MockWriteCloser {
	fileWriter := &MockWriteCloser{}
	for _, chunk := range writtenChunks {
		fileWriter.EXPECT().Write(chunk).Return(len(chunk), nil).Times(1)
	}
	fileWriter.EXPECT().Close().Return(nil).Times(1)
	return fileWriter
}

func CreateWriterCloserRanIntoWriteError(t *testing.T, successfulWrittenChunks [][]byte, errorAfterSuccessfulWrites error) *MockWriteCloser {
	fileWriter := &MockWriteCloser{}
	for _, chunk := range successfulWrittenChunks {
		fileWriter.EXPECT().Write(chunk).Return(len(chunk), nil).Times(1)
	}
	fileWriter.EXPECT().Write(mock.Anything).Return(0, errorAfterSuccessfulWrites).Times(1)

	fileWriter.EXPECT().Close().Return(nil).Times(1)
	return fileWriter
}

func CreateWriterCloserRanIntoCloseError(t *testing.T, writtenChunks [][]byte, err error) *MockWriteCloser {
	fileWriter := &MockWriteCloser{}
	for _, chunk := range writtenChunks {
		fileWriter.EXPECT().Write(chunk).Return(len(chunk), nil).Times(1)
	}

	fileWriter.EXPECT().Close().Return(err).Times(1)
	return fileWriter
}
