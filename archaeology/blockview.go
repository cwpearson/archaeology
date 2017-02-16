package archaeology

import (
	"bytes"
	"io"
	"os"
)

type BlockView struct {
	blockSize int64
	file      *os.File
	length    int64
}

func NewBlockView(f *os.File, blockSize int64) (*BlockView, error) {
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	numBlocks := (fi.Size() + 4095) / 4096
	return &BlockView{blockSize, f, numBlocks}, nil

}

func (b *BlockView) Length() int64 {
	return b.length
}

func (bv *BlockView) GetBlock(i int64) (*Block, error) {
	// Seek to the right file offset
	_, err := bv.file.Seek(i*bv.blockSize, io.SeekStart)
	if err != nil {
		return nil, err
	}

	// Read in the block
	block := &Block{}
	block.data = make([]byte, bv.blockSize)
	n, err := bv.file.Read(block.data)
	if err != nil && err != io.EOF {
		return nil, err
	}
	block.length = n

	return block, nil
}

func (bv *BlockView) Get(i int) (*Block, error) {
	return bv.GetBlock(int64(i))
}

type Block struct {
	data   []byte
	length int
}

func (b *Block) Equals(rhs *Block) bool {
	if b.length != rhs.length {
		return false
	}
	return bytes.Equal(b.data, rhs.data)
}
