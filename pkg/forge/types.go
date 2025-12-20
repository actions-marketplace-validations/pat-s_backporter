// Package forge provides abstraction over different git forges (GitHub, Forgejo, etc.).
package forge

import (
	"strings"
	"time"
)

// PRInfo contains information about a pull request.
type PRInfo struct {
	Number      int
	Title       string
	Body        string
	State       string
	MergeCommit string
	HeadSHA     string
	BaseBranch  string
	HeadBranch  string
	Merged      bool
	Squashed    bool
	Author      string
	MergedAt    time.Time
	Labels      []string
}

// HasBackportLabel checks if the PR has any label containing "backport".
func (p *PRInfo) HasBackportLabel() bool {
	for _, label := range p.Labels {
		if strings.Contains(strings.ToLower(label), "backport") {
			return true
		}
	}
	return false
}

// CommitInfo contains information about a commit.
type CommitInfo struct {
	SHA       string
	Message   string
	Author    string
	Email     string
	Timestamp time.Time
	Parents   []string
}

// IsSquashMerge checks if the PR was squash merged (single parent in merge commit).
func (p *PRInfo) IsSquashMerge() bool {
	return p.Squashed
}
