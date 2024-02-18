package io

import "testing"

func CreateWriterCloserMock(t *testing.T, writtenChunks [][]byte) *MockWriteCloser {
	fileWriter := &MockWriteCloser{}
	for _, chunk := range writtenChunks {
		fileWriter.EXPECT().Write(chunk).Return(len(chunk), nil).Times(1)
	}
	fileWriter.EXPECT().Close().Return(nil).Times(1)
	return fileWriter
}
