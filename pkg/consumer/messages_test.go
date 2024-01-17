package consumer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTelemetryMessage_isMessage(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Test TelemetryMessage isMessage()",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotPanics(t, TelemetryMessage{}.isMessage)
		})
	}
}
