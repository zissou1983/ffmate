package interceptor

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func PageLimit(c *gin.Context) {
	perPage := c.Query("perPage")
	page := c.Query("page")
	if perPage != "" {
		perPage, err := strconv.Atoi(perPage)
		if err != nil || perPage < 1 || perPage > 100 {
			c.JSON(400, gin.H{"error": "Invalid value for query string perPage"})
			c.Abort()
			return
		}
		c.Set("perPage", perPage)
	} else {
		c.Set("perPage", 50)
	}
	if page != "" {
		page, err := strconv.Atoi(page)
		if err != nil || page < 0 || page > 100 {
			c.JSON(400, gin.H{"error": "Invalid value for query string page"})
			c.Abort()
			return
		}
		c.Set("page", page)
	} else {
		c.Set("page", 0)
	}
	c.Next()
}
