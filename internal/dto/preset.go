package dto

import "time"

type Preset struct {
	Uuid string `json:"uuid"`

	Command     string `json:"command"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`

	OutputFile string `json:"outputFile"`

	Priority uint `json:"priority"`

	PreProcessing  *NewPrePostProcessing `json:"preProcessing,omitempty"`
	PostProcessing *NewPrePostProcessing `json:"postProcessing,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
