package cmd

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"math"
	"os"
	"runtime/pprof"
	"time"

	log "github.com/Sirupsen/logrus"
	arch "github.com/cwpearson/archaeology/archaeology"
	"github.com/cwpearson/archaeology/archaeology/adler"
	bv "github.com/cwpearson/archaeology/blockview"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(diffCmd)
}

const (
	blockSize = 4096
)

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

		f1, err := bv.New(args[0])
		if err != nil {
			log.Fatal(err)
		}
		defer f1.Close()

		f2, err := bv.New(args[1])
		if err != nil {
			log.Fatal(err)
		}
		defer f2.Close()

		f1Size, err := f1.Size()
		if err != nil {
			log.Fatal(err)
		}

		start := time.Now()
		f1BlockAdlerSums := map[uint32][]int64{} // map of file1 block checksums to offsets in file1
		f1BlockShaSums := map[[sha1.Size]byte][]int64{}

		err = nil
		for err == nil {
			// Read blocks
			k := f1.Offset()
			buf, err := f1.NextBlock(blockSize)
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}

			// Create the adler and sha1 sums
			adlerSum := adler.Sum(buf)
			f1BlockAdlerSums[adlerSum] = append(f1BlockAdlerSums[adlerSum], k)
			shaSum := sha1.Sum(buf)
			f1BlockShaSums[shaSum] = append(f1BlockShaSums[shaSum], k)
		}
		elapsed := time.Since(start)
		speed := float64(f1Size) / (1024 * 1024) / elapsed.Seconds()
		log.Info("Block Checksums: ", elapsed, speed, " MB/s")
		log.Info(len(f1BlockAdlerSums), " unique adler sums")
		log.Info(len(f1BlockShaSums), " unique sha1 sums")
		if len(f1BlockShaSums) != int(math.Ceil(float64(f1Size)/blockSize)) {
			log.Fatal("Sha1 collision in file1 blocks!")
		}

		start = time.Now()
		// Get file 2 size
		f2Size, err := f2.Size()
		if err != nil {
			log.Fatal(err)
		}

		file2Recipe := []*arch.Instruction{}
		rollingStart := int64(-1)

		err = nil
		for err == nil {
			k := f2.Offset()
			fmt.Printf("k = %d\n", k)

			if rollingStart >= 0 { // one byte at a time
				l := k + blockSize
				if l >= f2Size {
					break
				}
				// fmt.Printf("Rolling in byte %d...\n", l)
				data, err := f2.NextByte()
				if err == bv.ErrNoByte { // done!
					log.Info("Reached end of file 2")
					break
				} else if err != nil && err != io.EOF {
					log.Fatal(err)
				}
				k := f2.Offset()
				// Should be guaranteed to read 1 byte by here

				adlerSum := s.Roll(data)

				// Check if any blocks match the current one - that means the new region has ended
				if _, ok := f1BlockAdlerSums[adlerSum]; ok {
					// fmt.Printf("file2 @ %d matches somewhere in file1 (weak)\n", k)

					// Use the first matching block from file1
					f2BlockShaSum := sha1.Sum(buf)
					if f1BlockMatches, ok := f1BlockShaSums[f2BlockShaSum]; ok {
						f1BlockRef := f1BlockMatches[0]
						// fmt.Printf("file2 @ %d matches file1 @ %d (strong)\n", k, f1BlockRef)
						file2Recipe = append(file2Recipe, archaeology.NewBlockRef(f1BlockRef))
						k += blockSize
						rollingStart = -1 // end of new region

						continue
					} else {
						continue
					}
				}

				if err == io.EOF {
					// fmt.Println("Reached end of file2 while rolling")

					// Read the new data region at the end of the file
					newData := make([]byte, l-rollingStart)
					_, err = f2.Seek(rollingStart, io.SeekStart)
					if err != nil {
						log.Fatal(err)
					}
					n, err := f2.Read(newData)
					if err != nil {
						log.Fatal(err)
					}
					newData = newData[:n]
					file2Recipe = append(file2Recipe, archaeology.NewNewData(newData))
					break
				}

			} else { // one block at a time
				// Seek and read file
				_, err = f2.Seek(k, io.SeekStart)
				if err != nil {
					log.Fatal(err)
				}
				n, err := f2.Read(buf)
				if err != nil {
					log.Fatal(err)
				}
				buf = buf[:n]

				s = adler.NewSum(buf)
				adlerSum := s.Current()
				if k == 0 {
					log.Info(adlerSum)
				}
				if _, ok := f1BlockAdlerSums[adlerSum]; ok {
					// fmt.Printf("file2 block checksum (%d-%d) matches file1 block offsets %v\n", k, l, offsets)

					// Use the first matching block from file1
					f2BlockShaSum := sha1.Sum(buf)
					if f1BlockMatches, ok := f1BlockShaSums[f2BlockShaSum]; ok {
						f1BlockRef := f1BlockMatches[0]
						// fmt.Printf("file2[%d...] == file1[%d...]\n", k, f1BlockRef)
						file2Recipe = append(file2Recipe, archaeology.NewBlockRef(f1BlockRef))
						k += blockSize
						continue
					} else {
						// fmt.Printf("file2[%d...] eventually had no match \n", k)
						rollingStart = k
						continue
					}
				} else {
					// fmt.Printf("file2[%d...] had no weak match\n", k)
					rollingStart = k
					continue
				}
			}

		}

		elapsed = time.Since(start)
		speed = float64(f2Size/(1024*1024)) / elapsed.Seconds()
		log.Info("File2 Rolling Checksums: ", elapsed, speed, " MB/s")

		// for _, inst := range file2Recipe {
		// 	if inst.Ty == archaeology.BlockRef {
		// 		fmt.Print("r")
		// 	} else if inst.Ty == archaeology.NewData {
		// 		fmt.Print("+")
		// 	}
		// }
		// fmt.Println()

	},
}
