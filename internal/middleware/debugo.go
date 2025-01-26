package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/welovemedia/ffmate/sev"
	"github.com/yosev/debugo"
)

var debug = debugo.New("gin")

func Debugo(c *gin.Context, s *sev.Sev) {
	if strings.HasPrefix(c.Request.URL.Path, "/api/") {
		s.Metrics().GaugeVec("rest.api").WithLabelValues(c.Request.Method, c.FullPath()).Inc()
	}
	debug.Debugf("%s \"%s\"", c.Request.Method, c.Request.URL)
	c.Next()
}
