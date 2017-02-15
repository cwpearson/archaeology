package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/cwpearson/archaeology/archaeology"
	"github.com/cwpearson/archaeology/archaeology/levenshtein"
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
		bv1, err := archaeology.NewBlockView(f1, blockSize)
		if err != nil {
			log.Fatal(err)
		}

		f2, err := os.Open(args[1])
		if err != nil {
			log.Fatal(err)
		}
		bv2, err := archaeology.NewBlockView(f2, blockSize)
		if err != nil {
			log.Fatal(err)
		}

		script := levenshtein.EditScriptForStrings(bv1, bv2, levenshtein.DefaultLevenshtein)

		fmt.Println(script)

	},
}
