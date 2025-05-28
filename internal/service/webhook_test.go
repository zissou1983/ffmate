package service

import (
	"testing"

	"github.com/welovemedia/ffmate/internal/database/model"
	"github.com/welovemedia/ffmate/internal/dto"
	"github.com/welovemedia/ffmate/sev"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupWebhookTestDB(t *testing.T) (*gorm.DB, *sev.Sev) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = db.AutoMigrate(&model.Watchfolder{}, &model.Preset{}, &model.Webhook{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	s := sev.New("test", "", "", 3000)
	s.SetDB(db)
	Init(s)

	return db, s
}

func TestWebhookService(t *testing.T) {
	db, _ := setupWebhookTestDB(t)
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get underlying database: %v", err)
	}
	defer sqlDB.Close()

	t.Run("Create and list webhook", func(t *testing.T) {
		newWebhook := &dto.NewWebhook{
			Event: dto.TASK_CREATED,
			Url:   "http://localhost:8080/webhook",
		}

		_, err := WebhookService().NewWebhook(newWebhook)
		if err != nil {
			t.Fatalf("Failed to create webhook: %v", err)
		}

		webhooks, total, err := WebhookService().ListWebhooks(0, 10)
		if err != nil {
			t.Fatalf("Failed to list webhooks: %v", err)
		}

		if total < 1 {
			t.Error("Expected at least one webhook")
		}

		if (*webhooks)[0].Event != newWebhook.Event {
			t.Errorf("Expected event %s, got %s", newWebhook.Event, (*webhooks)[0].Event)
		}
	})

	t.Run("Delete webhook", func(t *testing.T) {
		newWebhook := &dto.NewWebhook{
			Event: dto.TASK_DELETED,
			Url:   "http://localhost:8080/webhook",
		}

		webhook, err := WebhookService().NewWebhook(newWebhook)
		if err != nil {
			t.Fatalf("Failed to create webhook: %v", err)
		}

		err = WebhookService().DeleteWebhook(webhook.Uuid)
		if err != nil {
			t.Fatalf("Failed to delete webhook: %v", err)
		}

		_, total, err := WebhookService().ListWebhooks(0, 10)
		if err != nil {
			t.Fatalf("Failed to list webhooks: %v", err)
		}

		if total != 1 { // One webhook from previous test
			t.Errorf("Expected 1 webhook, got %d", total)
		}
	})
}
