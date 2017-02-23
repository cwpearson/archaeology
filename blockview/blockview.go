package archaeology

import (
	"errors"
	"io"
	"os"
)

type BlockView struct {
	file   *os.File
	offset int64
}

var ErrNoByte = errors.New("No byte")

func New(path string) (*BlockView, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return &BlockView{f, 0}, nil
}

func (bv *BlockView) Close() {
	bv.file.Close()
}

func (bv *BlockView) Offset() int64 {
	return bv.offset
}

func (bv *BlockView) Size() (int64, error) {
	fi, err := bv.file.Stat()
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}

func (bv *BlockView) Read(buf []byte) (int, error) {
	_, err := bv.file.Seek(bv.offset, io.SeekStart)
	if err != nil {
		return 0, err
	}

	return bv.file.Read(buf)
}

func (bv *BlockView) NextBlock(i int64) ([]byte, error) {
	// Seek to the right file offset
	_, err := bv.file.Seek(bv.offset, io.SeekStart)
	if err != nil {
		return nil, err
	}

	// Read in the block
	buf := make([]byte, i)
	n, err := bv.file.Read(buf)

	bv.offset += int64(n)

	return buf, err
}

func (bv *BlockView) NextByte() (byte, error) {
	buf, err := bv.NextBlock(1)

	if len(buf) == 0 {
		return 0, ErrNoByte
	}

	return buf[0], err
}
