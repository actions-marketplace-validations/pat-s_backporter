// Package forge provides abstraction over different git forges (GitHub, Forgejo, etc.).
package forge

import "time"

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
