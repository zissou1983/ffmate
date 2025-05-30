package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/welovemedia/ffmate/internal/config"
	"github.com/welovemedia/ffmate/sev"
)

func TestVersion(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		appName        string
		appVersion     string
		expectedHeader string
	}{
		{
			name:           "Normal version header",
			appName:        "TestApp",
			appVersion:     "1.0.0",
			expectedHeader: "TestApp/v1.0.0",
		},
		{
			name:           "Empty version",
			appName:        "TestApp",
			appVersion:     "",
			expectedHeader: "TestApp/v",
		},
		{
			name:           "Empty app name",
			appName:        "",
			appVersion:     "1.0.0",
			expectedHeader: "/v1.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Set config values
			viper.Set("appName", tt.appName)
			viper.Set("appVersion", tt.appVersion)
			config.Init()

			s := &sev.Sev{} // Mock Sev instance

			// Call middleware
			Version(c, s)

			// Verify header
			if header := w.Header().Get("X-Server"); header != tt.expectedHeader {
				t.Errorf("Expected header %s, got %s", tt.expectedHeader, header)
			}
		})
	}
}
