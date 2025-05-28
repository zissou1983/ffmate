package service

import (
	"testing"

	"github.com/welovemedia/ffmate/internal/database/model"
	"github.com/welovemedia/ffmate/internal/dto"
	"github.com/welovemedia/ffmate/internal/metrics"
	"github.com/welovemedia/ffmate/sev"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupWatchfolderTestDB(t *testing.T) (*gorm.DB, *sev.Sev) {
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

	// setup metrics
	metrics := &metrics.Metrics{}
	for name, gauge := range metrics.Gauges() {
		s.Metrics().RegisterGauge(name, gauge)
	}
	for name, gauge := range metrics.GaugesVec() {
		s.Metrics().RegisterGaugeVec(name, gauge)
	}

	return db, s
}

func TestWatchfolderService(t *testing.T) {
	db, _ := setupWatchfolderTestDB(t)
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get underlying database: %v", err)
	}
	defer sqlDB.Close()

	// First create a preset that we'll need for the watchfolder
	preset, err := PresetService().NewPreset(&dto.NewPreset{
		Name:    "Test Preset",
		Command: "test command",
	})
	if err != nil {
		t.Fatalf("Failed to create preset: %v", err)
	}

	t.Run("Create and find watchfolder", func(t *testing.T) {
		newWatchfolder := &dto.NewWatchfolder{
			Path:   "/test/watch",
			Preset: preset.Uuid,
		}

		wf, err := WatchfolderService().NewWatchfolder(newWatchfolder)
		if err != nil {
			t.Fatalf("Failed to create watchfolder: %v", err)
		}

		found, err := WatchfolderService().GetWatchfolderById(wf.Uuid)
		if err != nil {
			t.Fatalf("Failed to find watchfolder: %v", err)
		}

		if found.Path != newWatchfolder.Path {
			t.Errorf("Expected path %s, got %s", newWatchfolder.Path, found.Path)
		}
	})

	t.Run("List watchfolders", func(t *testing.T) {
		wfs, total, err := WatchfolderService().ListWatchfolders(0, 10)
		if err != nil {
			t.Fatalf("Failed to list watchfolders: %v", err)
		}

		if total < 1 {
			t.Error("Expected at least one watchfolder")
		}

		if len(*wfs) < 1 {
			t.Error("Expected at least one watchfolder in list")
		}
	})

	t.Run("Update watchfolder", func(t *testing.T) {
		wfs, _, _ := WatchfolderService().ListWatchfolders(0, 1)
		wf := (*wfs)[0]
		wf.GrowthChecks = 10

		updated, err := WatchfolderService().UpdateWatchfolderInternal(&wf)
		if err != nil {
			t.Fatalf("Failed to update watchfolder: %v", err)
		}

		if updated.GrowthChecks != 10 {
			t.Error("Expected watchfolder to have a growth checks of 10")
		}
	})

	t.Run("Delete watchfolder", func(t *testing.T) {
		wfs, _, _ := WatchfolderService().ListWatchfolders(0, 1)
		wf := (*wfs)[0]

		err := WatchfolderService().DeleteWatchfolder(wf.Uuid)
		if err != nil {
			t.Fatalf("Failed to delete watchfolder: %v", err)
		}

		_, err = WatchfolderService().GetWatchfolderById(wf.Uuid)
		if err == nil {
			t.Error("Expected no error when requesting deleted watchfolder")
		}
	})
}
