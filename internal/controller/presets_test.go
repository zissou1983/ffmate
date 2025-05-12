package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/welovemedia/ffmate/internal/database/model"
	"github.com/welovemedia/ffmate/internal/dto"
	"github.com/welovemedia/ffmate/internal/metrics"
	"github.com/welovemedia/ffmate/internal/service"
	"github.com/welovemedia/ffmate/sev"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type webhookCall struct {
	event dto.WebhookEvent
	data  interface{}
}

func setupTestDB(t *testing.T) (*gorm.DB, *sev.Sev) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Auto migrate models
	err = db.AutoMigrate(&model.Preset{}, &model.Webhook{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	s := sev.New("test", "", "", 3000)
	s.SetDB(db)

	// Initialize services
	service.Init(s)

	// Setup metrics
	metrics := &metrics.Metrics{}
	for name, gauge := range metrics.Gauges() {
		s.Metrics().RegisterGauge(name, gauge)
	}
	for name, gauge := range metrics.GaugesVec() {
		s.Metrics().RegisterGaugeVec(name, gauge)
	}

	return db, s
}

func TestPresetController(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, s := setupTestDB(t)

	// Get underlying database for cleanup
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get underlying database: %v", err)
	}
	defer sqlDB.Close()

	// Setup webhook server and expected calls
	webhookURL, webhookCalls := setupWebhookServer(t, 3) // Expect 3 calls: created, deleted

	// Setup webhooks for preset events
	service.WebhookService().NewWebhook(&dto.NewWebhook{
		Event: dto.PRESET_CREATED,
		Url:   webhookURL,
	})
	service.WebhookService().NewWebhook(&dto.NewWebhook{
		Event: dto.PRESET_DELETED,
		Url:   webhookURL,
	})

	controller := &PresetController{
		Prefix: "",
	}
	controller.Setup(s)

	t.Run("Add preset with webhook", func(t *testing.T) {
		newPreset := dto.NewPreset{
			Name:        "Test Preset",
			Description: "Test Description",
			Command:     "ffmpeg -i ${INPUT_FILE} ${OUTPUT_FILE}",
		}

		body, _ := json.Marshal(newPreset)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/v1/presets", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		s.Gin().ServeHTTP(w, req)

		// Wait for webhook call
		select {
		case call := <-webhookCalls:
			if call.event != dto.PRESET_CREATED {
				t.Errorf("Expected webhook event %s, got %s", dto.PRESET_CREATED, call.event)
			}
		case <-time.After(time.Second):
			t.Error("Webhook call timeout")
		}

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response dto.Preset
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response.Name != newPreset.Name {
			t.Errorf("Expected name %s, got %s", newPreset.Name, response.Name)
		}
	})

	t.Run("List presets", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/v1/presets?page=0&perPage=10", nil)

		s.Gin().ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response []dto.Preset
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if len(response) < 1 {
			t.Error("Expected at least one preset in response")
		}
	})

	t.Run("Delete preset with webhook", func(t *testing.T) {
		// First create a preset to delete
		preset, err := service.PresetService().NewPreset(&dto.NewPreset{
			Name:    "To Delete",
			Command: "test",
		})
		if err != nil {
			t.Fatalf("Failed to create preset: %v", err)
		}

		// Wait for creation webhook first
		waitForWebhook(t, webhookCalls, dto.PRESET_CREATED)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", "/v1/presets/"+preset.Uuid, nil)

		s.Gin().ServeHTTP(w, req)

		// Then wait for deletion webhook
		waitForWebhook(t, webhookCalls, dto.PRESET_DELETED)

		if w.Code != http.StatusNoContent {
			t.Errorf("Expected status %d, got %d", http.StatusNoContent, w.Code)
		}
	})
}
