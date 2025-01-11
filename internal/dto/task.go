package dto

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type TaskStatus string

const (
	QUEUED          TaskStatus = "QUEUED"
	RUNNING         TaskStatus = "RUNNING"
	DONE_SUCCESSFUL TaskStatus = "DONE_SUCCESSFUL"
	DONE_ERROR      TaskStatus = "DONE_ERROR"
	DONE_CANCELED   TaskStatus = "DONE_CANCELED"
)

type PostProcessing struct {
	ScriptPath  string `json:"scriptPath"`
	SidecarPath string `json:"sidecarPath"`
}

type Task struct {
	Uuid  string `json:"uuid"`
	Batch string `json:"batch,omitempty"`

	Name string `json:"name,omitempty"`

	Command    string `json:"command"`
	InputFile  string `json:"inputFile"`
	OutputFile string `json:"outputFile"`

	Status   TaskStatus `json:"status"`
	Progress float64    `json:"progress"`

	Priority uint `json:"priority"`

	PostProcessing *PostProcessing `json:"postProcessing,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (p PostProcessing) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *PostProcessing) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, p)
}
