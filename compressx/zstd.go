package compressx

import (
	"io"

	"github.com/klauspost/compress/zstd"
)

func Zstd(in io.Reader, out io.Writer) error {
	enc, err := zstd.NewWriter(out)
	if err != nil {
		return err
	}
	_, err = io.Copy(enc, in)
	if err != nil {
		_ = enc.Close()
		return err
	}
	return enc.Close()
}

func UnZstd(in io.Reader, out io.Writer) error {
	d, err := zstd.NewReader(in)
	if err != nil {
		return err
	}
	defer d.Close()

	// Copy content...
	_, err = io.Copy(out, d)
	return err
}
