package model

import (
	"time"

	"github.com/welovemedia/ffmate/internal/dto"
	"gorm.io/gorm"
)

type Preset struct {
	ID uint `gorm:"primarykey"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Uuid string

	Command string
	Name    string

	OutputFile string

	Priority uint

	PreProcessing  *dto.NewPrePostProcessing `gorm:"type:json"`
	PostProcessing *dto.NewPrePostProcessing `gorm:"type:json"`

	Description string
}

func (m *Preset) ToDto() *dto.Preset {
	return &dto.Preset{
		Uuid: m.Uuid,

		Command:     m.Command,
		Name:        m.Name,
		Description: m.Description,

		OutputFile: m.OutputFile,

		Priority: m.Priority,

		PreProcessing:  m.PreProcessing,
		PostProcessing: m.PostProcessing,

		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (Preset) TableName() string {
	return "presets"
}
