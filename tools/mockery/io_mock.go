package io_mock

import "io"

type WriteCloser interface {
	io.WriteCloser
}

type ReadCloser interface {
	io.ReadCloser
}
