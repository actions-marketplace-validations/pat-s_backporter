// Package backport provides core backporting functionality.
package backport

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// CacheEntry represents a cached backport operation.
type CacheEntry struct {
	OriginalSHA  string    `json:"original_sha"`
	BackportSHA  string    `json:"backport_sha"`
	TargetBranch string    `json:"target_branch"`
	PRNumber     int       `json:"pr_number,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
	Message      string    `json:"message"`
}

// Cache manages the local cache of backported commits/PRs.
type Cache struct {
	path    string
	entries []CacheEntry
}

// NewCache creates a new cache instance.
func NewCache(path string) *Cache {
	if path == "" {
		home, err := os.UserHomeDir()
		if err == nil {
			path = filepath.Join(home, ".cache", "backporter", "history.json")
		}
	}

	cache := &Cache{path: path}
	_ = cache.load()

	return cache
}

// load loads the cache from disk.
func (c *Cache) load() error {
	if c.path == "" {
		return nil
	}

	data, err := os.ReadFile(c.path)
	if err != nil {
		if os.IsNotExist(err) {
			c.entries = []CacheEntry{}
			return nil
		}
		return err
	}

	return json.Unmarshal(data, &c.entries)
}

// save saves the cache to disk.
func (c *Cache) save() error {
	if c.path == "" {
		return nil
	}

	dir := filepath.Dir(c.path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	data, err := json.MarshalIndent(c.entries, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(c.path, data, 0o644)
}

// Add adds a new entry to the cache.
func (c *Cache) Add(entry CacheEntry) error {
	c.entries = append(c.entries, entry)
	return c.save()
}

// List returns all cache entries.
func (c *Cache) List() []CacheEntry {
	return c.entries
}

// FindByOriginalSHA finds entries by original SHA.
func (c *Cache) FindByOriginalSHA(sha string) []CacheEntry {
	var result []CacheEntry
	for _, entry := range c.entries {
		if entry.OriginalSHA == sha {
			result = append(result, entry)
		}
	}
	return result
}

// FindByPRNumber finds entries by PR number.
func (c *Cache) FindByPRNumber(number int) []CacheEntry {
	var result []CacheEntry
	for _, entry := range c.entries {
		if entry.PRNumber == number {
			result = append(result, entry)
		}
	}
	return result
}

// Clear clears all cache entries.
func (c *Cache) Clear() error {
	c.entries = []CacheEntry{}
	return c.save()
}
