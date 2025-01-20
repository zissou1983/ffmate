package dto

type NewPreset struct {
	Command string `json:"command"`

	Priority uint `json:"priority"`

	PreProcessing  *NewPrePostProcessing `json:"preProcessing"`
	PostProcessing *NewPrePostProcessing `json:"postProcessing"`

	Name        string `json:"name"`
	Description string `json:"description"`
}
