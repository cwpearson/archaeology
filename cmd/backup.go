package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/cwpearson/archaeology/archaeology"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var include []string
var ignore []string

func init() {
	backupCmd.PersistentFlags().StringSliceVarP(&include, "include", "i", []string{}, "Paths to include")
	backupCmd.PersistentFlags().StringSliceVarP(&ignore, "ignore", "x", []string{}, "Paths to ignore")
	viper.BindPFlag("include", backupCmd.PersistentFlags().Lookup("include"))
	viper.BindPFlag("ignore", backupCmd.PersistentFlags().Lookup("ignore"))

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

		ignore = viper.GetStringSlice("ignore")

		for _, i := range args {
			include = append(include, i)
		}

		// Get destination from config
		dest := "~/.archaeology/backups"

		err := archaeology.Backup(include, ignore, dest)
		if err != nil {
			log.Fatal(err)
		}
	},
}
