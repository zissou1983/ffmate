package repository

import (
	"github.com/google/uuid"
	"github.com/welovemedia/ffmate/internal/database/model"
	"github.com/welovemedia/ffmate/internal/dto"
	"gorm.io/gorm"
)

type Preset struct {
	DB *gorm.DB
}

func (t *Preset) Setup() {
	t.DB.AutoMigrate(&model.Preset{})
}

func (m *Preset) List(page int, perPage int) (*[]model.Preset, int64, error) {
	total, _ := m.Count()
	var presets = &[]model.Preset{}
	m.DB.Order("created_at DESC").Limit(perPage).Offset(perPage * page).Find(&presets)
	return presets, total, m.DB.Error
}

func (m *Preset) Update(w *model.Preset) error {
	m.DB.Save(w)
	return m.DB.Error
}

func (m *Preset) Delete(w *model.Preset) error {
	m.DB.Delete(w)
	return m.DB.Error
}

func (m *Preset) Create(newPreset *dto.NewPreset) (*model.Preset, error) {
	preset := &model.Preset{
		Uuid:           uuid.NewString(),
		Command:        newPreset.Command,
		Name:           newPreset.Name,
		Description:    newPreset.Description,
		Priority:       newPreset.Priority,
		OutputFile:     newPreset.OutputFile,
		PreProcessing:  newPreset.PreProcessing,
		PostProcessing: newPreset.PostProcessing,
	}
	db := m.DB.Create(preset)
	return preset, db.Error
}

func (m *Preset) First(uuid string) (*model.Preset, error) {
	var preset *model.Preset
	db := m.DB.Where("uuid", uuid).First(&preset)
	return preset, db.Error
}

func (m *Preset) FindByUuid(uuid string) (*model.Preset, error) {
	var preset *model.Preset
	db := m.DB.Where("uuid", uuid).Find(&preset)
	return preset, db.Error
}

func (m *Preset) Count() (int64, error) {
	var count int64
	db := m.DB.Model(&model.Preset{}).Count(&count)
	return count, db.Error
}

func (m *Preset) CountDeleted() (int64, error) {
	var count int64
	db := m.DB.Unscoped().Model(&model.Preset{}).Where("deleted_at IS NOT NULL").Count(&count)
	return count, db.Error
}
