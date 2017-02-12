package cmd

import "github.com/spf13/cobra"

func init() {
	RootCmd.AddCommand(restoreCmd)
}

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}
