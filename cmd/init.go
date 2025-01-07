package cmd

import (
	"fmt"
	"os"

	"github.com/sanbornm/go-selfupdate/selfupdate"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/welovemedia/ffmate/pkg/config"
)

var updater *selfupdate.Updater

var rootCmd = &cobra.Command{
	Use:   "ffmate",
	Short: "ffmate is a wrapper for ffmpeg",
	Long:  "ffmate is a wrapper for ffmpeg that adds a queue system on top of it",
}

func init() {
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "enable debugging")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
}

func Execute(args []string) {
	// unmarshal viper into config.Config
	config.Init()

	// parse cobra flags
	rootCmd.ParseFlags(args)

	if config.Config().Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
