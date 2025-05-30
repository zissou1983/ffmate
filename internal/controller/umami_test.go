package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/welovemedia/ffmate/internal/dto"
	"github.com/welovemedia/ffmate/internal/metrics"
	"github.com/welovemedia/ffmate/sev"
)

type UmamiPayload struct {
	Url      string `json:"url"`
	Screen   string `json:"screen"`
	Language string `json:"language"`
}

func TestUmamiController(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		payload        dto.Umami
		expectedStatus int
	}{
		{
			name: "Valid analytics data",
			payload: dto.Umami{
				Payload: dto.UmamiPayload{
					Url:    "/test",
					Screen: "1920x1080",
				},
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "Empty payload",
			payload:        dto.Umami{},
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := sev.New("test", "", "", 3000)

			// setup metrics
			metrics := &metrics.Metrics{}
			for name, gauge := range metrics.Gauges() {
				s.Metrics().RegisterGauge(name, gauge)
			}
			for name, gauge := range metrics.GaugesVec() {
				s.Metrics().RegisterGaugeVec(name, gauge)
			}

			controller := &UmamiController{
				Prefix: "",
			}
			controller.Setup(s)

			// Prepare request
			payload, _ := json.Marshal(tt.payload)
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/metrics/umami", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")

			// Execute request
			s.Gin().ServeHTTP(w, req)

			// Verify response
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
