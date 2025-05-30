package interceptor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestPageLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		queryParams    map[string]string
		expectedPage   int
		expectedPer    int
		expectedStatus int
	}{
		{
			name:           "Default values",
			queryParams:    map[string]string{},
			expectedPage:   0,
			expectedPer:    50,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid values",
			queryParams:    map[string]string{"page": "2", "perPage": "20"},
			expectedPage:   2,
			expectedPer:    20,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid perPage",
			queryParams:    map[string]string{"perPage": "150"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid page",
			queryParams:    map[string]string{"page": "-1"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Non-numeric values",
			queryParams:    map[string]string{"page": "abc", "perPage": "def"},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Create test request with query parameters
			req := httptest.NewRequest("GET", "/", nil)
			q := req.URL.Query()
			for k, v := range tt.queryParams {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()
			c.Request = req

			// Call the interceptor
			PageLimit(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
				return
			}

			if tt.expectedStatus == http.StatusOK {
				page, exists := c.Get("page")
				if !exists {
					t.Error("Expected page to be set")
				} else if page.(int) != tt.expectedPage {
					t.Errorf("Expected page %d, got %d", tt.expectedPage, page.(int))
				}

				perPage, exists := c.Get("perPage")
				if !exists {
					t.Error("Expected perPage to be set")
				} else if perPage.(int) != tt.expectedPer {
					t.Errorf("Expected perPage %d, got %d", tt.expectedPer, perPage.(int))
				}
			} else {
				var response map[string]string
				err := json.NewDecoder(w.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if _, exists := response["error"]; !exists {
					t.Error("Expected error message in response")
				}
			}
		})
	}
}
