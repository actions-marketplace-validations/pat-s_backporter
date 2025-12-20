package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRemoteURL(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		wantOwner string
		wantRepo  string
		wantError bool
	}{
		{
			name:      "HTTPS URL with .git",
			url:       "https://github.com/owner/repo.git",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantError: false,
		},
		{
			name:      "HTTPS URL without .git",
			url:       "https://github.com/owner/repo",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantError: false,
		},
		{
			name:      "SSH URL",
			url:       "git@github.com:owner/repo.git",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantError: false,
		},
		{
			name:      "SSH URL without .git",
			url:       "git@github.com:owner/repo",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantError: false,
		},
		{
			name:      "Forgejo HTTPS URL",
			url:       "https://codeberg.org/myorg/myrepo.git",
			wantOwner: "myorg",
			wantRepo:  "myrepo",
			wantError: false,
		},
		{
			name:      "Forgejo SSH URL",
			url:       "git@codeberg.org:myorg/myrepo.git",
			wantOwner: "myorg",
			wantRepo:  "myrepo",
			wantError: false,
		},
		{
			name:      "Invalid SSH URL missing colon",
			url:       "git@github.com/owner/repo.git",
			wantOwner: "",
			wantRepo:  "",
			wantError: true,
		},
		{
			name:      "Invalid HTTPS URL too few parts",
			url:       "https://github.com/repo",
			wantOwner: "",
			wantRepo:  "",
			wantError: true,
		},
		{
			name:      "HTTP URL",
			url:       "http://github.com/owner/repo.git",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			owner, repo, err := ParseRemoteURL(tt.url)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantOwner, owner)
				assert.Equal(t, tt.wantRepo, repo)
			}
		})
	}
}

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
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := looksLikeSHA(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Helper function from interactive.go that we're testing.
func looksLikeSHA(s string) bool {
	if len(s) < 7 {
		return false
	}

	for _, c := range s {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') && (c < 'A' || c > 'F') {
			return false
		}
	}

	return true
}
