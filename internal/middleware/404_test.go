package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestE404(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		initialStatus  int
		expectedStatus int
		shouldAbort    bool
	}{
		{
			name:           "404 status",
			initialStatus:  404,
			expectedStatus: 404,
			shouldAbort:    true,
		},
		{
			name:           "200 status",
			initialStatus:  200,
			expectedStatus: 200,
			shouldAbort:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, engine := gin.CreateTestContext(w)

			// Create a test endpoint that sets the status
			engine.GET("/test", func(c *gin.Context) {
				c.Status(tt.initialStatus)
			})

			// Perform the request
			req := httptest.NewRequest("GET", "/test", nil)
			c.Request = req
			engine.HandleContext(c)

			// Now call our middleware
			E404(c, nil)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.shouldAbort != c.IsAborted() {
				t.Errorf("Expected aborted to be %v but was %v", tt.shouldAbort, c.IsAborted())
			}
		})
	}
}
