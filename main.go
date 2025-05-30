package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ffmate",
	Short: "FFmate - A modern automation layer for FFmpeg",
	Long: `FFmate is a modern and powerful automation layer built on top of FFmpeg.
It provides REST API, Web UI, Webhooks, and more for video/audio transcoding.`,
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the FFmate server",
	Long:  `Start the FFmate server with REST API and Web UI`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting FFmate server...")
		fmt.Println("API will be available at http://localhost:3000")
		// TODO: Implement server startup logic
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
