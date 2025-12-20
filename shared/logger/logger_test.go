package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsCI(t *testing.T) {
	tests := []struct {
		name     string
		envKey   string
		envValue string
		expected bool
	}{
		{
			name:     "CI env set",
			envKey:   "CI",
			envValue: "true",
			expected: true,
		},
		{
			name:     "GITHUB_ACTIONS env set",
			envKey:   "GITHUB_ACTIONS",
			envValue: "true",
			expected: true,
		},
		{
			name:     "no CI env",
			envKey:   "",
			envValue: "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear relevant env vars using t.Setenv which auto-restores.
			t.Setenv("CI", "")
			t.Setenv("GITHUB_ACTIONS", "")

			if tt.envKey != "" {
				t.Setenv(tt.envKey, tt.envValue)
			}

			result := IsCI()
			assert.Equal(t, tt.expected, result)
		})
	}
}
