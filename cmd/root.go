package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/Sirupsen/logrus"
)

var RootCmd = &cobra.Command{
	Use:   "archaeology",
	Short: "Short",
	Long: `multi
    line`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	// Do Stuff Here
	// },
}

func init() {

	RootCmd.PersistentFlags().Bool("viper", true, "Use Viper for configuration")
	viper.SetConfigName("config")             // name of config file (without extension)
	viper.AddConfigPath("/etc/archaeology/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.archaeology") // call multiple times to add many search paths
	viper.AddConfigPath(".")                  // optionally look for config in the working directory

	viper.SetConfigType("json")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatal(err)
	}
}
