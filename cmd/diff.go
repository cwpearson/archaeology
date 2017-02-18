package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime/pprof"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/cwpearson/archaeology/archaeology"
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

		cpuprofile, err := RootCmd.PersistentFlags().GetString("cpuprofile")
		if err != nil {
			log.Fatal(err)
		}
		if cpuprofile != "" {
			log.Infof("Doing CPU profiling, output=%s\n", cpuprofile)
			f, err := os.Create(cpuprofile)
			if err != nil {
				log.Fatal(err)
			}
			pprof.StartCPUProfile(f)
			defer pprof.StopCPUProfile()
		}

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

		file2Recipe := []*archaeology.Instruction{}

		// blockRegionStart := int64(-1)
		blockRegionEnd := int64(0)
		// newRegionStart := int64(0)

		for k := int64(0); k < f2Size; k++ {
			l := k + blockSize
			checksum := s.Current()
			if f1Offsets, ok := checksums[checksum]; ok {
				// fmt.Printf("file2 block checksum (%d-%d) matches file1 block offsets %v\n", k, l, offsets)

				// Actually compare the blocks
				matches, err := compareBlocks(blockSize, f1, f1Offsets, f2, k)
				if err != nil {
					log.Fatal(err)
				}
				if len(matches) > 0 {
					if k == blockRegionEnd { // k is at the end of the previous block
						fmt.Printf("Found a matching block at %d (matches %d)\n", k, matches[0])
						// blockRegionStart = k
						blockRegionEnd = l
						file2Recipe = append(file2Recipe, archaeology.NewBlockRef(matches[0])) // Add the matching block
					} else if k > blockRegionEnd { // k is past the end of a block Region (in a newRegion)
						fmt.Printf("Found a match at %d to end a newRegion (started %d)\n", k, blockRegionEnd)
						// blockRegionStart = k
						blockRegionEnd = l
						file2Recipe = append(file2Recipe, archaeology.NewNewData([]byte{})) // Add the new data
					}
				} else { // k-l doesn't match anything (there is new data in there)
					// if k < blockRegionEnd && newRegionStart < blockRegionEnd {
					// 	fmt.Printf("Found the beginning of a potential new region at %d\n", l)
					// 	newRegionStart = l
					// }
				}
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

		for _, inst := range file2Recipe {
			if inst.Ty == archaeology.BlockRef {
				fmt.Print("r")
			} else if inst.Ty == archaeology.NewData {
				fmt.Print("+")
			}
		}
		fmt.Println()

	},
}
