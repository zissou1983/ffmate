package dto

type Watchfolder struct {
	Uuid string `json:"uuid"`

	Name        string `json:"name"`
	Description string `json:"description"`

	Path        string `json:"path"`
	Interval    int    `json:"interval"`
	GrowthCheck int    `json:"growthCheck"`

	Suspended bool `json:"suspended"`

	CreatedAt int64 `json:"createdAt"`
	UpdatedAt int64 `json:"updatedAt"`
}
