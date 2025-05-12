package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/welovemedia/ffmate/sev"
)

func E404(c *gin.Context, _ *sev.Sev) {
	if c.Writer.Status() == 404 {
		c.AbortWithStatus(404)
		return
	}

	c.Next()
}
