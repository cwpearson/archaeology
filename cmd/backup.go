package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/cwpearson/archaeology/archaeology"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(backupCmd)
}

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatal("Expected at least one directory or file")
		}

		includes := args
		ignores := []string{}

		// Get destination from config
		dest := "~/.archaeology/backups"

		archaeology.Backup(includes, ignores, dest)
	},
}
