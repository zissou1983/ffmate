package main

import (
	"embed"
	_ "embed"
	"os"

	"github.com/spf13/viper"
	"github.com/welovemedia/ffmate/cmd"
	"github.com/welovemedia/ffmate/docs"
)

//go:embed .version
var version string

//go:embed all:ui-build/*
var frontend embed.FS

// @title ffmate API
// @version
// @description	A wrapper around ffmpeg

// @contact.name We love media
// @contact.email sev@welovemedia.io

// @license.name AGPL-3.0
// @license.url https://opensource.org/license/agpl-v3

// @host localhost
// @BasePath /api/v1
func main() {
	viper.Set("appName", "ffmate")
	viper.Set("appVersion", version)

	docs.SwaggerInfo.Schemes = []string{"http"}

	cmd.Execute(os.Args, frontend)
}
