package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/welovemedia/ffmate/internal/database/model"
	"github.com/welovemedia/ffmate/internal/dto"
	"github.com/welovemedia/ffmate/internal/metrics"
	"github.com/welovemedia/ffmate/internal/service"
	"github.com/welovemedia/ffmate/sev"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupWatchfolderTestDB(t *testing.T) (*gorm.DB, *sev.Sev) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = db.AutoMigrate(&model.Watchfolder{}, &model.Preset{}, &model.Webhook{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	s := sev.New("test", "", "", 3000)
	s.SetDB(db)
	service.Init(s)

	metrics := &metrics.Metrics{}
	for name, gauge := range metrics.Gauges() {
		s.Metrics().RegisterGauge(name, gauge)
	}
	for name, gauge := range metrics.GaugesVec() {
		s.Metrics().RegisterGaugeVec(name, gauge)
	}

	return db, s
}

func TestWatchfolderController(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, s := setupWatchfolderTestDB(t)
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get underlying database: %v", err)
	}
	defer sqlDB.Close()

	// Setup webhook server and expected calls
	webhookURL, webhookCalls := setupWebhookServer(t, 4) // Expect: preset_created, watchfolder_created, watchfolder_deleted

	// Setup webhooks for watchfolder events
	service.WebhookService().NewWebhook(&dto.NewWebhook{
		Event: dto.PRESET_CREATED,
		Url:   webhookURL,
	})
	service.WebhookService().NewWebhook(&dto.NewWebhook{
		Event: dto.WATCHFOLDER_CREATED,
		Url:   webhookURL,
	})
	service.WebhookService().NewWebhook(&dto.NewWebhook{
		Event: dto.WATCHFOLDER_DELETED,
		Url:   webhookURL,
	})

	// Create a preset first that we'll need for the watchfolder and wait for its webhook
	preset, err := service.PresetService().NewPreset(&dto.NewPreset{
		Name:    "Test Preset",
		Command: "test command",
	})
	if err != nil {
		t.Fatalf("Failed to create preset: %v", err)
	}
	waitForWebhook(t, webhookCalls, dto.PRESET_CREATED)

	controller := &WatchfolderController{
		Prefix: "",
	}
	controller.Setup(s)

	t.Run("Add watchfolder", func(t *testing.T) {
		newWatchfolder := dto.NewWatchfolder{
			Path:   "/test/watch",
			Preset: preset.Uuid,
		}

		body, _ := json.Marshal(newWatchfolder)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/v1/watchfolders", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		s.Gin().ServeHTTP(w, req)
		waitForWebhook(t, webhookCalls, dto.WATCHFOLDER_CREATED)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response dto.Watchfolder
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response.Path != newWatchfolder.Path {
			t.Errorf("Expected path %s, got %s", newWatchfolder.Path, response.Path)
		}
	})

	t.Run("List watchfolders", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/v1/watchfolders?page=0&perPage=10", nil)
		s.Gin().ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response []dto.Watchfolder
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if len(response) < 1 {
			t.Error("Expected at least one watchfolder")
		}
	})

	// Get first watchfolder for subsequent tests
	var firstWatchfolder dto.Watchfolder
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/v1/watchfolders?page=0&perPage=1", nil)
	s.Gin().ServeHTTP(w, req)
	var watchfolders []dto.Watchfolder
	json.Unmarshal(w.Body.Bytes(), &watchfolders)
	if len(watchfolders) > 0 {
		firstWatchfolder = watchfolders[0]

		t.Run("Get single watchfolder", func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/v1/watchfolders/"+firstWatchfolder.Uuid, nil)
			s.Gin().ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
			}
		})

		t.Run("Delete watchfolder", func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/v1/watchfolders/"+firstWatchfolder.Uuid, nil)
			s.Gin().ServeHTTP(w, req)
			waitForWebhook(t, webhookCalls, dto.WATCHFOLDER_DELETED)

			if w.Code != http.StatusNoContent {
				t.Errorf("Expected status %d, got %d", http.StatusNoContent, w.Code)
			}
		})
	}
}
