package dto

type NewWatchfolder struct {
	Name        string `json:"name"`
	Description string `json:"description"`

	Path         string `json:"path"`
	Interval     int    `json:"interval"`
	GrowthChecks int    `json:"growthChecks"`

	Filter *WatchfolderFilter `json:"filter"`

	Suspended bool `json:"suspended"`

	Preset string `json:"preset"`
}
