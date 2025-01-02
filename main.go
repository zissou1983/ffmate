package main

import (
	_ "embed"
	"os"

	"github.com/welovemedia/ffmate/cmd"
)

//go:embed .version
var version string

func main() {
	cmd.Execute(os.Args, version)
}
