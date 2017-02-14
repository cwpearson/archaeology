package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(patchCmd)
}

var patchCmd = &cobra.Command{
	Use:   "patch",
	Short: "Apply an archaeology edit script to a file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		log.Fatal("Unimplemented.")
	},
}
