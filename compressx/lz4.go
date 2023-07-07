package compressx

import (
	"errors"

	"github.com/pierrec/lz4/v4"
)

func LZ4(data []byte) ([]byte, error) {
	buf := make([]byte, lz4.CompressBlockBound(len(data)))
	var c lz4.Compressor
	n, err := c.CompressBlock(data, buf)
	if err != nil {
		return nil, err
	}
	if n >= len(data) {
		return nil, errors.New("not compressible")
	}
	return buf[:n], nil
}

func UnLZ4(data []byte) ([]byte, error) {
	// Allocate a very large buffer for decompression.
	out := make([]byte, 10*len(data))
	n, err := lz4.UncompressBlock(data, out)
	if err != nil {
		return nil, err
	}

	return out[:n], nil
}
