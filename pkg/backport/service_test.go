package backport

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"codefloe.com/pat-s/backporter/pkg/config"
)

func TestNewService(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		Cache: config.CacheConfig{
			Enabled: true,
			Path:    filepath.Join(tmpDir, "cache.json"),
		},
	}

	service := NewService(nil, nil, cfg, "owner", "repo")

	assert.NotNil(t, service)
	assert.Equal(t, "owner", service.owner)
	assert.Equal(t, "repo", service.repoN)
	assert.NotNil(t, service.cache)
}

func TestNewServiceWithCacheDisabled(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		Cache: config.CacheConfig{
			Enabled: false,
			Path:    filepath.Join(tmpDir, "cache.json"),
		},
	}

	service := NewService(nil, nil, cfg, "owner", "repo")

	assert.NotNil(t, service)
	// Cache is still created but disabled (won't persist).
	assert.NotNil(t, service.cache)
}

func TestListBackportsEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		Cache: config.CacheConfig{
			Enabled: false,
			Path:    filepath.Join(tmpDir, "cache.json"),
		},
	}

	service := NewService(nil, nil, cfg, "owner", "repo")
	entries := service.ListBackports()

	assert.Empty(t, entries)
}

func TestClearCache(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		Cache: config.CacheConfig{
			Enabled: false,
			Path:    filepath.Join(tmpDir, "cache.json"),
		},
	}

	service := NewService(nil, nil, cfg, "owner", "repo")
	err := service.ClearCache()

	assert.NoError(t, err)
}
