package utils

import (
	"bytes"
	"testing"
)

func TestLogBroadcaster_Write(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		wantN    int
		wantErr  bool
	}{
		{
			name:     "Write normal log message",
			input:    []byte("test log message"),
			wantN:    16,
			wantErr:  false,
		},
		{
			name:     "Write empty message",
			input:    []byte{},
			wantN:    0,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var received []byte
			lb := &LogBroadcaster{
				Callback: func(p []byte) {
					received = p
				},
			}

			gotN, err := lb.Write(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("LogBroadcaster.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if gotN != tt.wantN {
				t.Errorf("LogBroadcaster.Write() = %v, want %v", gotN, tt.wantN)
			}

			if !bytes.Equal(received, tt.input) {
				t.Errorf("LogBroadcaster callback received %v, want %v", received, tt.input)
			}
		})
	}
}
