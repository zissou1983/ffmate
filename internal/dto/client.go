package dto

type Client struct {
	Version string `json:"version"`
	Os      string `json:"os"`
	Arch    string `json:"arch"`
}
