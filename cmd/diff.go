package cmd

import (
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

			if k < 0 {
				log.Fatal("k should be non-negative")
			}
			s := adler.NewSum(buf, uint64(k))
			checksums = append(checksums, s.Current())
		}
		elapsed := time.Since(start)
		speed := float64(f1Size/(1024*1024)) / elapsed.Seconds()
		log.Info("Block Checksums: ", elapsed, speed, " MB/s")

		// fmt.Println(rollingChecksums)

	},
}
