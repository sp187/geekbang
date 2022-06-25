package web

import (
	"io"
)

// FileStream 文件流，抽象化
type FileStream interface {
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Closer
	Size() int64
}
