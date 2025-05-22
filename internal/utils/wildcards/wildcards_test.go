package wildcards

import (
	"runtime"
	"testing"
)

func TestReplace(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		inputFile   string
		outputFile  string
		source      string
		escapePaths bool
		want        string
	}{
		{
			name:        "File paths with spaces",
			input:       "${INPUT_FILE} to ${OUTPUT_FILE}",
			inputFile:   "/path/to/input file.mp4",
			outputFile:  "/path/to/output file.mp4",
			source:      "test",
			escapePaths: true,
			want:        "/path/to/input\\ file.mp4 to /path/to/output\\ file.mp4",
		},
		{
			name:       "File components",
			input:      "${INPUT_FILE_BASE} ${INPUT_FILE_EXTENSION} ${INPUT_FILE_BASENAME} ${INPUT_FILE_DIR}",
			inputFile:  "/path/to/input.mp4",
			outputFile: "/path/to/output.mp4",
			source:     "test",
			want:       "input.mp4 .mp4 input /path/to",
		},
		{
			name:       "System info",
			input:      "OS: ${OS_NAME} ${OS_ARCH}",
			inputFile:  "test.mp4",
			outputFile: "out.mp4",
			source:     "test",
			want:       "OS: " + runtime.GOOS + " " + runtime.GOARCH,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Replace(tt.input, tt.inputFile, tt.outputFile, tt.source, tt.escapePaths)
			if got != tt.want {
				t.Errorf("Replace() = %v, want %v", got, tt.want)
			}
		})
	}
}
