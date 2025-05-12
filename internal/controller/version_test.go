package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/welovemedia/ffmate/internal/config"
	"github.com/welovemedia/ffmate/internal/dto"
	"github.com/welovemedia/ffmate/sev"
)

func TestVersionController(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	versionString := "1.0.0-test"
	s := sev.New("test", versionString, "", 3000)

	// Set test version in config
	viper.Set("appVersion", versionString)
	config.Init()

	controller := &VersionController{
		Prefix: "",
	}
	controller.Setup(s)

	// Test getVersion endpoint
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/version", nil)
	s.Gin().ServeHTTP(w, req)

	// Assert response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response dto.Version
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatal("Failed to unmarshal response:", err)
	}

	if response.Version != versionString {
		t.Errorf("Expected version %s, got %s", versionString, response.Version)
	}
}
