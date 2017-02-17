package cmd

import (
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
		checksums := []uint32{}
		firstOffsets := map[uint32]int64{}
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
			checksums = append(checksums, checksum)
			if _, ok := firstOffsets[checksum]; !ok {
				firstOffsets[checksum] = k
			}
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

		// Create an initial checksum of f2
		buf = make([]byte, blockSize)
		n, err := f2.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		buf = buf[:n]

		s := adler.NewSum(buf, uint64(0))

		for k := int64(0); k < f2Size; k++ {
			l := k + blockSize

			checksum := s.Current()
			offset, ok := firstOffsets[checksum]
			if ok {
				fmt.Printf("file2 offset %d checksum matches file1 offset %d\n", k, offset)
			} else {
				// fmt.Printf("New checksum at offset %d\n", k)
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

	},
}
