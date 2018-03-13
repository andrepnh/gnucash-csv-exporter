package app

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func decompress(t *testing.T, file string) *bufio.Reader {
	gzipFile, err := os.Open("../xml-samples/" + file)
	if assert.NoError(t, err, "Could not open file") {
		return Decompress(gzipFile)
	}
	panic(err)
}

func TestShouldDecompress(t *testing.T) {
	reader := decompress(t, "simple-compressed.gnucash")
	line, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	assert.True(t, strings.HasPrefix(strings.TrimSpace(line), "<?xml"))
}
