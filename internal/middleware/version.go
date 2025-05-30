package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/welovemedia/ffmate/internal/config"
	"github.com/welovemedia/ffmate/sev"
)

func Version(c *gin.Context, s *sev.Sev) {
	c.Writer.Header().Set("X-Server", config.Config().AppName+"/v"+config.Config().AppVersion)
	c.Next()
}
