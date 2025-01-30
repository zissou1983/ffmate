package ffmpeg

import (
	"context"
	"fmt"
	"math"

	"github.com/sirupsen/logrus"
	"github.com/welovemedia/ffmate/internal/database/model"
)

// FFmpegProgress holds parsed progress data
type FFmpegProgress struct {
	Frame   int
	FPS     float64
	Bitrate string
	Time    float64
	Speed   string
}

// EstimateRemainingTime calculates the estimated remaining time based on the current progress and speed.
func (p *FFmpegProgress) EstimateRemainingTime(duration float64) (float64, error) {
	speed, err := parseSpeed(p.Speed)
	if err != nil {
		return 0, err
	}
	remainingTime := (duration - p.Time) / speed
	return math.Round(remainingTime), nil
}

// parseSpeed parses the speed string and returns the speed as a float64.
func parseSpeed(speedStr string) (float64, error) {
	var speed float64
	_, err := fmt.Sscanf(speedStr, "%fx", &speed)
	if err != nil {
		return 0, err
	}
	return speed, nil
}

type ExecutionRequest struct {
	Task *model.Task

	Command string

	Logger *logrus.Logger

	UpdateFunc func(progress float64, remaining float64)

	Ctx context.Context
}
