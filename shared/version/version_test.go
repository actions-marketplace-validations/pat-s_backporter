package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	originalVersion := Version
	originalBuildDate := BuildDate
	defer func() {
		Version = originalVersion
		BuildDate = originalBuildDate
	}()

	Version = "1.2.3"
	BuildDate = "2025-01-15"
	assert.Equal(t, "1.2.3 (2025-01-15)", String())
}

func TestFull(t *testing.T) {
	original := Version
	defer func() { Version = original }()

	Version = "1.2.3"
	full := Full()

	assert.Contains(t, full, "1.2.3")
	assert.Contains(t, full, GitURL)
	assert.Contains(t, full, "backporter")
}

func TestSignatureMessage(t *testing.T) {
	original := Version
	defer func() { Version = original }()

	Version = "2.0.0"
	sha := "abc123def456"
	msg := SignatureMessage(sha)

	assert.Contains(t, msg, sha)
	assert.Contains(t, msg, "2.0.0")
	assert.Contains(t, msg, GitURL)
	assert.Contains(t, msg, "Backported from")
}
