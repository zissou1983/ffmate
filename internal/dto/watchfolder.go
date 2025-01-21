package dto

type Watchfolder struct {
	Uuid string `json:"uuid"`

	Name        string `json:"name"`
	Description string `json:"description"`

	Path         string `json:"path"`
	Interval     int    `json:"interval"`
	GrowthChecks int    `json:"growthChecks"`

	Suspended bool `json:"suspended"`

	Preset string `json:"preset"`

	CreatedAt int64 `json:"createdAt"`
	UpdatedAt int64 `json:"updatedAt"`
}
