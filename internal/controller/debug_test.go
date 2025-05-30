package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/welovemedia/ffmate/sev"
)

func TestDebugController(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "Set debug namespace",
			method:         "PATCH",
			path:           "/v1/debug/namespace/test",
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "Disable debug",
			method:         "DELETE",
			path:           "/v1/debug/namespace",
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := sev.New("test", "", "", 3000)
			controller := &DebugController{
				Prefix: "",
			}
			controller.Setup(s)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.path, nil)

			s.Gin().ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
