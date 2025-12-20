package backport

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLooksLikeSHA(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"abc1234", true},
		{"0123456789abcdef", true},
		{"ABC1234", true},
		{"abcdefg", false}, // 'g' is not hex.
		{"abc12", false},   // Too short.
		{"", false},
		{"abc123!", false},
		{"1234567", true},
		{"deadbeef1234567890", true},
		{"DEADBEEF1234567890", true},
		{"0000000", true},
		{"fffffff", true},
		{"FFFFFFF", true},
		{"123456", false}, // 6 chars, too short.
		{"xyz1234", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := looksLikeSHA(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
