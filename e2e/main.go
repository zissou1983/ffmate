package main

import (
	"embed"
	"fmt"
	"os"
	"time"

	"github.com/welovemedia/ffmate/cmd"
	"github.com/welovemedia/ffmate/e2e/tests/presets"
	"github.com/welovemedia/ffmate/e2e/tests/tasks"
)

//go:embed ui/*
var frontend embed.FS

func main() {
	go cmd.Execute(os.Args, frontend)
	time.Sleep(1 * time.Second)

	// Run test suites
	if err := presets.RunTests(); err != nil {
		fmt.Printf("Preset tests failed: %v\n", err)
		os.Exit(1)
	}

	if err := tasks.RunTests(); err != nil {
		fmt.Printf("Task tests failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("All tests passed successfully!")
	os.Exit(0)
}
