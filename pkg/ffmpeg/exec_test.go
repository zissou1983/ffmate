package ffmpeg

import (
	"fmt"
	"testing"
	"time"

	"gopkg.in/go-playground/assert.v1"
)

func TestEvaluateWildcards(t *testing.T) {
	// Input and expected values
	inputFile := "input.mp4"
	outputFile := "output.mp4"

	// Capturing the current time once to avoid inconsistencies
	currentTime := time.Now()

	// Command template with placeholders
	cmd := "-i ${INPUT_FILE} -test ${DATE_YEAR} ${DATE_SHORTYEAR} ${DATE_MONTH} ${DATE_DAY} ${TIME_HOUR} ${TIME_MINUTE} ${TIME_SECOND} ${OUTPUT_FILE}"

	// Creating the execution request
	request := &ExecutionRequest{Task: nil, Logger: nil, Command: cmd, InputFile: inputFile, OutputFile: outputFile}

	// Evaluate the wildcards in the command
	evaluateWildcards(request)

	// Expected command after wildcard replacement
	expectedCommand := fmt.Sprintf(
		"-i %s -test %s %s %s %s %s %s %s %s",
		inputFile,
		currentTime.Format("2006"), // Year (4 digits)
		currentTime.Format("06"),   // Year (2 digits)
		currentTime.Format("01"),   // Month
		currentTime.Format("02"),   // Day
		currentTime.Format("15"),   // Hour
		currentTime.Format("04"),   // Minute
		currentTime.Format("05"),   // Second
		outputFile,
	)

	// Asserting that the evaluated command matches the expected command
	assert.Equal(t, expectedCommand, request.Command)
}
