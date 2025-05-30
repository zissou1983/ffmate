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

func setupTaskTestDB(t *testing.T) (*gorm.DB, *sev.Sev) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = db.AutoMigrate(&model.Task{}, &model.Preset{}, &model.Webhook{})
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

func TestTaskController(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, s := setupTaskTestDB(t)
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get underlying database: %v", err)
	}
	defer sqlDB.Close()

	// Setup webhook server and expected calls
	webhookURL, webhookCalls := setupWebhookServer(t, 5) // Expect: task_created (single), batch_created, task_created (batch), task_deleted

	// Setup webhooks for task events
	service.WebhookService().NewWebhook(&dto.NewWebhook{
		Event: dto.TASK_CREATED,
		Url:   webhookURL,
	})
	service.WebhookService().NewWebhook(&dto.NewWebhook{
		Event: dto.TASK_DELETED,
		Url:   webhookURL,
	})
	service.WebhookService().NewWebhook(&dto.NewWebhook{
		Event: dto.BATCH_CREATED,
		Url:   webhookURL,
	})

	controller := &TaskController{
		Prefix: "",
	}
	controller.Setup(s)

	t.Run("Add single task", func(t *testing.T) {
		newTask := dto.NewTask{
			InputFile:  "/test/input.mp4",
			OutputFile: "/test/output.mp4",
			Command:    "ffmpeg -i ${INPUT_FILE} ${OUTPUT_FILE}",
		}

		body, _ := json.Marshal(newTask)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/v1/tasks", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		s.Gin().ServeHTTP(w, req)
		waitForWebhook(t, webhookCalls, dto.TASK_CREATED)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("Add batch tasks", func(t *testing.T) {
		newTasks := []dto.NewTask{
			{
				InputFile:  "/test/input1.mp4",
				OutputFile: "/test/output1.mp4",
				Command:    "test command 1",
			},
			{
				InputFile:  "/test/input2.mp4",
				OutputFile: "/test/output2.mp4",
				Command:    "test command 2",
			},
		}

		body, _ := json.Marshal(newTasks)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/v1/tasks/batch", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		s.Gin().ServeHTTP(w, req)

		// Wait for both batch created and task created events
		waitForWebhook(t, webhookCalls, dto.TASK_CREATED) // First task
		waitForWebhook(t, webhookCalls, dto.TASK_CREATED) // Second task
		waitForWebhook(t, webhookCalls, dto.BATCH_CREATED)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response []dto.Task
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if len(response) != 2 {
			t.Errorf("Expected 2 tasks in response, got %d", len(response))
		}
	})

	t.Run("List tasks", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/v1/tasks?page=0&perPage=10", nil)

		s.Gin().ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	// Get first task for subsequent tests
	var firstTask dto.Task
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/v1/tasks?page=0&perPage=1", nil)
	s.Gin().ServeHTTP(w, req)
	var tasks []dto.Task
	json.Unmarshal(w.Body.Bytes(), &tasks)
	if len(tasks) > 0 {
		firstTask = tasks[0]

		t.Run("Get single task", func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/v1/tasks/"+firstTask.Uuid, nil)
			s.Gin().ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
			}
		})

		t.Run("Cancel task", func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("PATCH", "/v1/tasks/"+firstTask.Uuid+"/cancel", nil)
			s.Gin().ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
			}
		})

		t.Run("Restart task", func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("PATCH", "/v1/tasks/"+firstTask.Uuid+"/restart", nil)
			s.Gin().ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
			}
		})

		t.Run("Delete task", func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/v1/tasks/"+firstTask.Uuid, nil)
			s.Gin().ServeHTTP(w, req)
			waitForWebhook(t, webhookCalls, dto.TASK_DELETED)

			if w.Code != http.StatusNoContent {
				t.Errorf("Expected status %d, got %d", http.StatusNoContent, w.Code)
			}
		})
	}
}
