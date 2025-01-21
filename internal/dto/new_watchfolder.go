package dto

type NewWatchfolder struct {
	Name        string `json:"name"`
	Description string `json:"description"`

	Path        string `json:"path"`
	Interval    int    `json:"interval"`
	GrowthCheck int    `json:"growthCheck"`
}
