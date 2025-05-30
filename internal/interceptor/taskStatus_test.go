package interceptor

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/welovemedia/ffmate/internal/dto"
)

func TestTaskStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		queryStatus    string
		expectedStatus string
	}{
		{
			name:           "Empty status",
			queryStatus:    "",
			expectedStatus: "",
		},
		{
			name:           "Status QUEUED",
			queryStatus:    "QUEUED",
			expectedStatus: string(dto.QUEUED),
		},
		{
			name:           "Status RUNNING",
			queryStatus:    "RUNNING",
			expectedStatus: string(dto.RUNNING),
		},
		{
			name:           "Status DONE_SUCCESSFUL",
			queryStatus:    "DONE_SUCCESSFUL",
			expectedStatus: string(dto.DONE_SUCCESSFUL),
		},
		{
			name:           "Status DONE_ERROR",
			queryStatus:    "DONE_ERROR",
			expectedStatus: string(dto.DONE_ERROR),
		},
		{
			name:           "Status DONE_CANCELED",
			queryStatus:    "DONE_CANCELED",
			expectedStatus: string(dto.DONE_CANCELED),
		},
		{
			name:           "Invalid status",
			queryStatus:    "INVALID_STATUS",
			expectedStatus: "",
		},
		{
			name:           "Lowercase status",
			queryStatus:    "queued",
			expectedStatus: string(dto.QUEUED),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Create test request with query parameter
			req := httptest.NewRequest("GET", "/", nil)
			q := req.URL.Query()
			if tt.queryStatus != "" {
				q.Add("status", tt.queryStatus)
			}
			req.URL.RawQuery = q.Encode()
			c.Request = req

			// Call the interceptor
			TaskStatus(c)

			// Check the result
			status, exists := c.Get("status")

			// Only check existence for valid status values
			if tt.expectedStatus != "" && !exists {
				t.Error("Expected status to be set for valid status value")
				return
			}

			if exists && status.(string) != tt.expectedStatus {
				t.Errorf("Expected status %s, got %s", tt.expectedStatus, status.(string))
			}
		})
	}
}
