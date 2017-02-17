package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/cwpearson/archaeology/archaeology/adler"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(diffCmd)
}

func compareBlocks(blockSize int, f1 *os.File, f1BlockOffsets []int64, f2 *os.File, f2BlockOffset int64) ([]int64, error) {
	f2Buf := make([]byte, blockSize)
	_, err := f2.Seek(f2BlockOffset, io.SeekStart)
	if err != nil {
		return nil, err
	}
	n, err := f2.Read(f2Buf)
	if err != nil && err != io.EOF {
		return nil, err
	}
	f2Buf = f2Buf[:n]
	if err == io.EOF {
	}

	matches := []int64{}
	f1Buf := make([]byte, blockSize)
	for _, f1BlockOffset := range f1BlockOffsets {
		f1Buf = f1Buf[:cap(f1Buf)]
		_, err := f1.Seek(f1BlockOffset, io.SeekStart)
		if err != nil {
			return nil, err
		}
		n, err = f1.Read(f1Buf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		f1Buf = f1Buf[:n]
		if err == io.EOF {
		}

		if bytes.Equal(f1Buf, f2Buf) {
			matches = append(matches, f1BlockOffset)
		}
	}
	return matches, nil
}

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Compute block-level diff between files",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			log.Fatal("Expected two files as arguments.")
		}

		f1, err := os.Open(args[0])
		if err != nil {
			log.Fatal(err)
		}
		defer f1.Close()

		f2, err := os.Open(args[1])
		if err != nil {
			log.Fatal(err)
		}
		defer f2.Close()

		fi, err := f1.Stat()
		if err != nil {
			log.Fatal(err)
		}
		f1Size := fi.Size()

		start := time.Now()
		checksums := map[uint32][]int64{} // map of file1 block checksums to offsets in file1

		buf := make([]byte, blockSize)

		for k := int64(0); k < f1Size; k += blockSize {

			buf = buf[:cap(buf)] // resize buffer to its original capacity
			_, err := f1.Seek(k, io.SeekStart)
			if err != nil {
				log.Fatal(err)
			}

			n, err := f1.Read(buf)
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}
			buf = buf[:n]

			s := adler.NewSum(buf, uint64(k))
			checksum := s.Current()
			checksums[checksum] = append(checksums[checksum], k)

			if err == io.EOF {
				log.Warn("Reached end of file1")
				break
			}
		}
		elapsed := time.Since(start)
		speed := float64(f1Size/(1024*1024)) / elapsed.Seconds()
		log.Info("Block Checksums: ", elapsed, speed, " MB/s")

		fi, err = f2.Stat()
		if err != nil {
			log.Fatal(err)
		}
		f2Size := fi.Size()

		start = time.Now()
		// Create an initial checksum of f2
		buf = make([]byte, blockSize)
		n, err := f2.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		buf = buf[:n]

		s := adler.NewSum(buf, uint64(0))

		endPrevMatch := int64(-1)
		for k := int64(0); k < f2Size; k++ {
			l := k + blockSize

			checksum := s.Current()
			if offsets, ok := checksums[checksum]; ok {
				// fmt.Printf("file2 block checksum (%d-%d) matches file1 block offsets %v\n", k, l, offsets)

				// Actually compare the blocks
				matches, err := compareBlocks(blockSize, f1, offsets, f2, k)
				if err != nil {
					log.Fatal(err)
				}
				if len(matches) > 0 {
					if k > endPrevMatch {
						fmt.Printf("New data in file2 at %d\n", k)
					}
					fmt.Printf("file2 block (%d-%d) matches file1 offsets %v\n", k, l, matches)
					endPrevMatch = l
				}

				endPrevMatch = l
			} else {

			}

			f2.Seek(l, io.SeekStart)
			oneByte := make([]byte, 1)
			n, err := f2.Read(oneByte)
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}
			oneByte = oneByte[:n]

			if len(oneByte) == 1 {
				s.Roll(oneByte[0])
			} else {
				log.Warn("Read zero bytes")
				break
			}
			if err == io.EOF {
				log.Warn("Read end of file2")
				break

			}

		}
		elapsed = time.Since(start)
		speed = float64(f2Size/(1024*1024)) / elapsed.Seconds()
		log.Info("File2 Rolling Checksums: ", elapsed, speed, " MB/s")

	},
}
