package main

import (
	_ "embed"
	"os"

	"github.com/welovemedia/ffmate/cmd"
	"github.com/welovemedia/ffmate/docs"
)

//go:embed .version
var version string

//	@title			ffmate API
//	@version		{will be injected}
//	@description	A wrapper around ffmpeg

//	@contact.name	We love media
//	@contact.email	sev@welovemedia.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost
// @BasePath	/api/v1
func main() {
	docs.SwaggerInfo.Schemes = []string{"http"}

	cmd.Execute(os.Args, version)
}
