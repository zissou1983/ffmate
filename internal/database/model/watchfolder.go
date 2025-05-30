package model

import (
	"github.com/welovemedia/ffmate/internal/dto"
	"gorm.io/gorm"
)

type Watchfolder struct {
	ID uint `gorm:"primarykey"`

	CreatedAt int64          `gorm:"autoCreateTime:milli"`
	UpdatedAt int64          `gorm:"autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Uuid string

	Name        string
	Description string

	Path         string
	Interval     int
	GrowthChecks int

	Filter *dto.WatchfolderFilter

	Preset string

	Suspended bool

	Error     string
	LastCheck int64
}

func (m *Watchfolder) ToDto() *dto.Watchfolder {
	return &dto.Watchfolder{
		Uuid: m.Uuid,

		Name:        m.Name,
		Description: m.Description,

		Path:         m.Path,
		Interval:     m.Interval,
		GrowthChecks: m.GrowthChecks,

		Preset: m.Preset,

		Filter: m.Filter,

		Suspended: m.Suspended,

		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,

		Error:     m.Error,
		LastCheck: m.LastCheck,
	}
}

func (Watchfolder) TableName() string {
	return "watchfolder"
}
