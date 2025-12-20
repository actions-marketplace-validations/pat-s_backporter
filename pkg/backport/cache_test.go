package backport

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCacheAddAndList(t *testing.T) {
	tmpDir := t.TempDir()
	cachePath := filepath.Join(tmpDir, "cache.json")

	cache := NewCache(cachePath)
	assert.Empty(t, cache.List())

	entry := CacheEntry{
		OriginalSHA:  "abc123def456",
		BackportSHA:  "789xyz000111",
		TargetBranch: "release-1.0",
		PRNumber:     42,
		Timestamp:    time.Now(),
		Message:      "Fix critical bug",
	}

	err := cache.Add(entry)
	require.NoError(t, err)

	entries := cache.List()
	assert.Len(t, entries, 1)
	assert.Equal(t, entry.OriginalSHA, entries[0].OriginalSHA)
	assert.Equal(t, entry.BackportSHA, entries[0].BackportSHA)
	assert.Equal(t, entry.TargetBranch, entries[0].TargetBranch)
	assert.Equal(t, entry.PRNumber, entries[0].PRNumber)
}

func TestCacheFindByOriginalSHA(t *testing.T) {
	tmpDir := t.TempDir()
	cachePath := filepath.Join(tmpDir, "cache.json")

	cache := NewCache(cachePath)

	entries := []CacheEntry{
		{
			OriginalSHA:  "sha1",
			BackportSHA:  "backport1",
			TargetBranch: "release-1.0",
			Timestamp:    time.Now(),
		},
		{
			OriginalSHA:  "sha2",
			BackportSHA:  "backport2",
			TargetBranch: "release-2.0",
			Timestamp:    time.Now(),
		},
		{
			OriginalSHA:  "sha1",
			BackportSHA:  "backport3",
			TargetBranch: "release-3.0",
			Timestamp:    time.Now(),
		},
	}

	for _, entry := range entries {
		err := cache.Add(entry)
		require.NoError(t, err)
	}

	// Find by SHA1 - should return 2 entries.
	found := cache.FindByOriginalSHA("sha1")
	assert.Len(t, found, 2)
	assert.Equal(t, "backport1", found[0].BackportSHA)
	assert.Equal(t, "backport3", found[1].BackportSHA)

	// Find by SHA2 - should return 1 entry.
	found = cache.FindByOriginalSHA("sha2")
	assert.Len(t, found, 1)
	assert.Equal(t, "backport2", found[0].BackportSHA)

	// Find non-existent SHA.
	found = cache.FindByOriginalSHA("nonexistent")
	assert.Empty(t, found)
}

func TestCacheFindByPRNumber(t *testing.T) {
	tmpDir := t.TempDir()
	cachePath := filepath.Join(tmpDir, "cache.json")

	cache := NewCache(cachePath)

	entries := []CacheEntry{
		{
			OriginalSHA:  "sha1",
			BackportSHA:  "backport1",
			TargetBranch: "release-1.0",
			PRNumber:     100,
			Timestamp:    time.Now(),
		},
		{
			OriginalSHA:  "sha2",
			BackportSHA:  "backport2",
			TargetBranch: "release-2.0",
			PRNumber:     200,
			Timestamp:    time.Now(),
		},
		{
			OriginalSHA:  "sha3",
			BackportSHA:  "backport3",
			TargetBranch: "release-3.0",
			PRNumber:     100,
			Timestamp:    time.Now(),
		},
	}

	for _, entry := range entries {
		err := cache.Add(entry)
		require.NoError(t, err)
	}

	// Find by PR 100 - should return 2 entries.
	found := cache.FindByPRNumber(100)
	assert.Len(t, found, 2)

	// Find by PR 200 - should return 1 entry.
	found = cache.FindByPRNumber(200)
	assert.Len(t, found, 1)

	// Find non-existent PR.
	found = cache.FindByPRNumber(999)
	assert.Empty(t, found)
}

func TestCacheClear(t *testing.T) {
	tmpDir := t.TempDir()
	cachePath := filepath.Join(tmpDir, "cache.json")

	cache := NewCache(cachePath)

	// Add some entries.
	for i := 0; i < 3; i++ {
		entry := CacheEntry{
			OriginalSHA:  "sha",
			BackportSHA:  "backport",
			TargetBranch: "release",
			Timestamp:    time.Now(),
		}
		err := cache.Add(entry)
		require.NoError(t, err)
	}

	assert.Len(t, cache.List(), 3)

	// Clear cache.
	err := cache.Clear()
	require.NoError(t, err)

	assert.Empty(t, cache.List())
}

func TestCachePersistence(t *testing.T) {
	tmpDir := t.TempDir()
	cachePath := filepath.Join(tmpDir, "cache.json")

	// Create cache and add entry.
	cache1 := NewCache(cachePath)
	entry := CacheEntry{
		OriginalSHA:  "persistent-sha",
		BackportSHA:  "persistent-backport",
		TargetBranch: "release-1.0",
		PRNumber:     42,
		Timestamp:    time.Now(),
		Message:      "Test persistence",
	}
	err := cache1.Add(entry)
	require.NoError(t, err)

	// Create new cache instance from same path.
	cache2 := NewCache(cachePath)

	// Should have the same entry.
	entries := cache2.List()
	assert.Len(t, entries, 1)
	assert.Equal(t, "persistent-sha", entries[0].OriginalSHA)
	assert.Equal(t, "persistent-backport", entries[0].BackportSHA)
	assert.Equal(t, 42, entries[0].PRNumber)
}

func TestCacheEmptyPath(t *testing.T) {
	// Use a temp directory to avoid loading default cache.
	tmpDir := t.TempDir()
	cachePath := filepath.Join(tmpDir, "nonexistent", "cache.json")

	cache := NewCache(cachePath)
	assert.Empty(t, cache.List())

	entry := CacheEntry{
		OriginalSHA:  "sha",
		BackportSHA:  "backport",
		TargetBranch: "release",
		Timestamp:    time.Now(),
	}

	err := cache.Add(entry)
	assert.NoError(t, err)

	// Should have the entry in memory.
	entries := cache.List()
	assert.Len(t, entries, 1)
}
