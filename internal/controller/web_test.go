package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/welovemedia/ffmate/sev"
)

func TestWebController(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "Root redirect to UI",
			method:         "GET",
			path:           "/",
			expectedStatus: http.StatusMovedPermanently,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := sev.New("test", "", "", 3000)
			controller := &WebController{
				Prefix: "",
			}
			controller.Setup(s)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.path, nil)
			s.Gin().ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusMovedPermanently {
				if location := w.Header().Get("Location"); location != "/ui" {
					t.Errorf("Expected redirect to /ui, got %s", location)
				}
			}
		})
	}
}
