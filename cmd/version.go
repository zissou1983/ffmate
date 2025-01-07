package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/welovemedia/ffmate/pkg/config"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print ffmate version",
	Run:   version,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func version(cmd *cobra.Command, args []string) {
	fmt.Printf("version: %s\n", config.Config().AppVersion)
}
