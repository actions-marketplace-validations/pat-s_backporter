package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	assert.Equal(t, "", cfg.ForgeType)
	assert.Empty(t, cfg.TargetBranches)
	assert.Equal(t, "", cfg.CommitMessage)
	assert.Equal(t, "", cfg.AuthorName)
	assert.Equal(t, "", cfg.AuthorEmail)
	assert.Equal(t, "main", cfg.DefaultBranch)
	assert.Equal(t, "origin", cfg.Remote)
	assert.True(t, cfg.Cache.Enabled)
	assert.Equal(t, "", cfg.Cache.Path)
}

func TestConfigMerge(t *testing.T) {
	tests := []struct {
		name     string
		base     *Config
		other    *Config
		expected *Config
	}{
		{
			name: "merge empty config",
			base: DefaultConfig(),
			other: &Config{
				ForgeType: "github",
			},
			expected: &Config{
				ForgeType:      "github",
				TargetBranches: []string{},
				CommitMessage:  "",
				AuthorName:     "",
				AuthorEmail:    "",
				DefaultBranch:  "main",
				Remote:         "origin",
				Cache: CacheConfig{
					Enabled: false,
					Path:    "",
				},
			},
		},
		{
			name: "merge full config",
			base: DefaultConfig(),
			other: &Config{
				ForgeType:      "forgejo",
				TargetBranches: []string{"release-.*"},
				CommitMessage:  "custom message",
				AuthorName:     "Test Author",
				AuthorEmail:    "test@example.com",
				DefaultBranch:  "develop",
				Remote:         "upstream",
				Cache: CacheConfig{
					Enabled: true,
					Path:    "/tmp/cache.json",
				},
			},
			expected: &Config{
				ForgeType:      "forgejo",
				TargetBranches: []string{"release-.*"},
				CommitMessage:  "custom message",
				AuthorName:     "Test Author",
				AuthorEmail:    "test@example.com",
				DefaultBranch:  "develop",
				Remote:         "upstream",
				Cache: CacheConfig{
					Enabled: true,
					Path:    "/tmp/cache.json",
				},
			},
		},
		{
			name:     "merge nil config",
			base:     DefaultConfig(),
			other:    nil,
			expected: DefaultConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.base.Merge(tt.other)
			assert.Equal(t, tt.expected.ForgeType, tt.base.ForgeType)
			assert.Equal(t, tt.expected.TargetBranches, tt.base.TargetBranches)
			assert.Equal(t, tt.expected.CommitMessage, tt.base.CommitMessage)
			assert.Equal(t, tt.expected.AuthorName, tt.base.AuthorName)
			assert.Equal(t, tt.expected.AuthorEmail, tt.base.AuthorEmail)
			assert.Equal(t, tt.expected.DefaultBranch, tt.base.DefaultBranch)
			assert.Equal(t, tt.expected.Remote, tt.base.Remote)
			assert.Equal(t, tt.expected.Cache.Path, tt.base.Cache.Path)
		})
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name      string
		config    *Config
		wantError bool
	}{
		{
			name:      "valid empty forge type",
			config:    DefaultConfig(),
			wantError: false,
		},
		{
			name: "valid github forge type",
			config: &Config{
				ForgeType: "github",
			},
			wantError: false,
		},
		{
			name: "valid forgejo forge type",
			config: &Config{
				ForgeType: "forgejo",
			},
			wantError: false,
		},
		{
			name: "invalid forge type",
			config: &Config{
				ForgeType: "gitlab",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLoadFromFile(t *testing.T) {
	// Create a temporary config file.
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `
forge_type: github
target_branches:
  - release-1.x
  - release-2.x
commit_message: "backport: {{.OriginalMessage}}"
author_name: "Backporter Bot"
author_email: "bot@example.com"
default_branch: main
remote: origin
cache:
  enabled: true
  path: /tmp/backporter-cache.json
`

	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err)

	cfg, err := LoadFromFile(configPath)
	require.NoError(t, err)

	assert.Equal(t, "github", cfg.ForgeType)
	assert.Equal(t, []string{"release-1.x", "release-2.x"}, cfg.TargetBranches)
	assert.Equal(t, "backport: {{.OriginalMessage}}", cfg.CommitMessage)
	assert.Equal(t, "Backporter Bot", cfg.AuthorName)
	assert.Equal(t, "bot@example.com", cfg.AuthorEmail)
	assert.Equal(t, "main", cfg.DefaultBranch)
	assert.Equal(t, "origin", cfg.Remote)
	assert.True(t, cfg.Cache.Enabled)
	assert.Equal(t, "/tmp/backporter-cache.json", cfg.Cache.Path)
}

func TestLoadFromFileNotFound(t *testing.T) {
	_, err := LoadFromFile("/nonexistent/path/config.yaml")
	assert.Error(t, err)
}

func TestSaveToFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "subdir", "config.yaml")

	cfg := &Config{
		ForgeType:      "github",
		TargetBranches: []string{"release-.*"},
		CommitMessage:  "test message",
		AuthorName:     "Test",
		AuthorEmail:    "test@test.com",
		DefaultBranch:  "main",
		Remote:         "origin",
		Cache: CacheConfig{
			Enabled: true,
			Path:    "/tmp/cache.json",
		},
	}

	err := cfg.SaveToFile(configPath)
	require.NoError(t, err)

	// Verify file exists.
	_, err = os.Stat(configPath)
	assert.NoError(t, err)

	// Load it back.
	loaded, err := LoadFromFile(configPath)
	require.NoError(t, err)

	assert.Equal(t, cfg.ForgeType, loaded.ForgeType)
	assert.Equal(t, cfg.TargetBranches, loaded.TargetBranches)
	assert.Equal(t, cfg.CommitMessage, loaded.CommitMessage)
	assert.Equal(t, cfg.AuthorName, loaded.AuthorName)
	assert.Equal(t, cfg.AuthorEmail, loaded.AuthorEmail)
}

func TestGlobalConfigPath(t *testing.T) {
	path := GlobalConfigPath()
	// Should contain .config/backporter.
	assert.Contains(t, path, ".config")
	assert.Contains(t, path, "backporter")
	assert.Contains(t, path, "config.yaml")
}

func TestRepoConfigPath(t *testing.T) {
	path := RepoConfigPath()
	assert.Equal(t, ".backporter.yaml", path)
}
