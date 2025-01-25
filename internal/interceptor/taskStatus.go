package interceptor

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/welovemedia/ffmate/internal/dto"
)

func TaskStatus(c *gin.Context) {
	status := c.Query("status")
	if status != "" {
		switch strings.ToUpper(status) {
		case "QUEUED":
			c.Set("status", string(dto.QUEUED))
		case "RUNNING":
			c.Set("status", string(dto.RUNNING))
		case "DONE_SUCCESSFUL":
			c.Set("status", string(dto.DONE_SUCCESSFUL))
		case "DONE_ERROR":
			c.Set("status", string(dto.DONE_ERROR))
		case "DONE_CANCELED":
			c.Set("status", string(dto.DONE_CANCELED))
		}
	} else {
		c.Set("status", "")
	}
	c.Next()
}
