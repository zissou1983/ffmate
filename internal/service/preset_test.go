package service

import (
	"testing"

	"github.com/welovemedia/ffmate/internal/database/model"
	"github.com/welovemedia/ffmate/internal/dto"
	"github.com/welovemedia/ffmate/sev"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, *sev.Sev) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = db.AutoMigrate(&model.Preset{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	s := sev.New("test", "", "", 3000)
	s.SetDB(db)
	Init(s) // Initialize services

	return db, s
}

func TestPresetService(t *testing.T) {
	db, _ := setupTestDB(t)
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get underlying database: %v", err)
	}
	defer sqlDB.Close()

	t.Run("Create and find preset", func(t *testing.T) {
		newPreset := &dto.NewPreset{
			Name:        "Test Preset",
			Description: "Test Description",
			Command:     "ffmpeg -i ${INPUT_FILE} ${OUTPUT_FILE}",
		}

		preset, err := PresetService().NewPreset(newPreset)
		if err != nil {
			t.Fatalf("Failed to create preset: %v", err)
		}

		found, err := PresetService().FindByUuid(preset.Uuid)
		if err != nil {
			t.Fatalf("Failed to find preset: %v", err)
		}

		if found.Name != newPreset.Name {
			t.Errorf("Expected preset name %s, got %s", newPreset.Name, found.Name)
		}
	})

	t.Run("List presets", func(t *testing.T) {
		presets, total, err := PresetService().ListPresets(0, 10)
		if err != nil {
			t.Fatalf("Failed to list presets: %v", err)
		}

		if total < 1 {
			t.Error("Expected at least one preset")
		}

		if len(*presets) < 1 {
			t.Error("Expected at least one preset in list")
		}
	})

	t.Run("Delete preset", func(t *testing.T) {
		newPreset := &dto.NewPreset{
			Name:    "To Delete",
			Command: "test",
		}

		preset, err := PresetService().NewPreset(newPreset)
		if err != nil {
			t.Fatalf("Failed to create preset: %v", err)
		}

		err = PresetService().DeletePreset(preset.Uuid)
		if err != nil {
			t.Fatalf("Failed to delete preset: %v", err)
		}

		_, err = PresetService().FindByUuid(preset.Uuid)
		if err == nil {
			t.Error("Expected error when finding deleted preset")
		}
	})
}
