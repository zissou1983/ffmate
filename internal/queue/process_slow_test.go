//go:build slow

package queue

import (
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/welovemedia/ffmate/internal/config"
	"github.com/welovemedia/ffmate/internal/database/model"
	"github.com/welovemedia/ffmate/internal/database/repository"
	"github.com/welovemedia/ffmate/internal/dto"
	"github.com/welovemedia/ffmate/internal/metrics"
	"github.com/welovemedia/ffmate/internal/service"
	"github.com/welovemedia/ffmate/sev"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, *sev.Sev) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = db.AutoMigrate(&model.Task{}, &model.Webhook{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	s := sev.New("test", "", "", 3000)
	s.SetDB(db)

	// Set config values before initializing services
	viper.Set("maxConcurrentTasks", uint(2))
	config.Init()

	metrics := &metrics.Metrics{}
	for name, gauge := range metrics.Gauges() {
		s.Metrics().RegisterGauge(name, gauge)
	}
	for name, gauge := range metrics.GaugesVec() {
		s.Metrics().RegisterGaugeVec(name, gauge)
	}

	service.Init(s)

	return db, s
}

func TestQueue(t *testing.T) {
	db, s := setupTestDB(t)
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get underlying database: %v", err)
	}
	defer sqlDB.Close()

	queue := &Queue{
		Sev:                s,
		TaskRepository:     &repository.Task{DB: db},
		MaxConcurrentTasks: config.Config().MaxConcurrentTasks,
	}
	queue.Init()

	t.Run("Process task", func(t *testing.T) {
		task := &model.Task{
			InputFile:  &dto.RawResolved{Raw: "/test/input.mp4"},
			OutputFile: &dto.RawResolved{Raw: "/test/output.mp4"},
			Command:    &dto.RawResolved{Raw: "echo test"},
			Status:     dto.QUEUED,
		}

		db.Create(task)

		// Wait for task to be processed
		deadline := time.Now().Add(5 * time.Second)
		for time.Now().Before(deadline) {
			db.First(task, task.ID)
			if task.Status == dto.DONE_ERROR {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}

		if task.Status != dto.DONE_ERROR {
			t.Errorf("Expected task status %s, got %s", dto.DONE_ERROR, task.Status)
		}
		if task.Error == "" {
			t.Error("Expected error message to be set")
		}
	})
}
