package repository

import (
	"github.com/google/uuid"
	"github.com/welovemedia/ffmate/internal/database/model"
	"github.com/welovemedia/ffmate/internal/dto"
	"gorm.io/gorm"
)

type Watchfolder struct {
	DB *gorm.DB
}

func (t *Watchfolder) Setup() {
	t.DB.AutoMigrate(&model.Watchfolder{})
}

func (m *Watchfolder) List(page int, perPage int) (*[]model.Watchfolder, int64, error) {
	total, _ := m.Count()
	var watchfolder = &[]model.Watchfolder{}
	if page >= 0 && perPage >= 0 {
		m.DB.Order("created_at DESC").Limit(perPage).Offset(page * perPage).Find(&watchfolder)
	} else {
		m.DB.Order("created_at DESC").Find(&watchfolder)
	}
	return watchfolder, total, m.DB.Error
}

func (m *Watchfolder) First(uuid string) (*model.Watchfolder, error) {
	var watchfolder = &model.Watchfolder{}
	err := m.DB.Where("uuid = ?", uuid).First(&watchfolder).Error
	if err != nil {
		return nil, err
	}
	return watchfolder, nil
}

func (m *Watchfolder) Count() (int64, error) {
	var count int64
	db := m.DB.Model(&model.Watchfolder{}).Count(&count)
	return count, db.Error
}

func (m *Watchfolder) Delete(w *model.Watchfolder) error {
	m.DB.Delete(w)
	return m.DB.Error
}

func (m *Watchfolder) Update(w *model.Watchfolder) (*model.Watchfolder, error) {
	m.DB.Save(w)
	return w, m.DB.Error
}

func (m *Watchfolder) Create(newWatchfolder *dto.NewWatchfolder) (*model.Watchfolder, error) {
	watchfolder := &model.Watchfolder{
		Uuid:         uuid.NewString(),
		Name:         newWatchfolder.Name,
		Description:  newWatchfolder.Description,
		Preset:       newWatchfolder.Preset,
		Path:         newWatchfolder.Path,
		Interval:     newWatchfolder.Interval,
		Filter:       newWatchfolder.Filter,
		GrowthChecks: newWatchfolder.GrowthChecks,
		Suspended:    newWatchfolder.Suspended,
	}
	db := m.DB.Create(watchfolder)
	return watchfolder, db.Error
}
