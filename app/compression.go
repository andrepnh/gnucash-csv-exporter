package app

import (
	"bufio"
	"compress/gzip"
	"io"
)

// Decompress - decompress with gzip
func Decompress(file io.Reader) *bufio.Reader {
	reader, err := gzip.NewReader(bufio.NewReader(file))
	if err != nil {
		panic(err)
	}
	return bufio.NewReader(reader)
}
