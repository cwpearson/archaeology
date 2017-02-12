package cmd

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "archaeology",
	Short: "Short",
	Long: `multi
    line`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	// Do Stuff Here
	// },
}
