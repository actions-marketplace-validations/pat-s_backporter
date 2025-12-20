package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// CherryPickResult represents the result of a cherry-pick operation.
type CherryPickResult struct {
	Success     bool
	HasConflict bool
	Message     string
}

// CherryPick performs a git cherry-pick operation.
// Note: go-git doesn't support cherry-pick natively, so we use git command.
func CherryPick(sha string) (*CherryPickResult, error) {
	cmd := exec.Command("git", "cherry-pick", sha)
	output, err := cmd.CombinedOutput()
	if err != nil {
		outputStr := string(output)

		// Check if it's a conflict.
		if strings.Contains(outputStr, "CONFLICT") || strings.Contains(outputStr, "after resolving the conflicts") {
			return &CherryPickResult{
				Success:     false,
				HasConflict: true,
				Message:     outputStr,
			}, nil
		}

		return nil, fmt.Errorf("cherry-pick failed: %s - %w", outputStr, err)
	}

	return &CherryPickResult{
		Success:     true,
		HasConflict: false,
		Message:     string(output),
	}, nil
}

// AbortCherryPick aborts an in-progress cherry-pick.
func AbortCherryPick() error {
	cmd := exec.Command("git", "cherry-pick", "--abort")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to abort cherry-pick: %w", err)
	}
	return nil
}

// ContinueCherryPick continues a cherry-pick after conflicts are resolved.
func ContinueCherryPick() error {
	cmd := exec.Command("git", "cherry-pick", "--continue")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to continue cherry-pick: %w", err)
	}
	return nil
}

// CheckoutBranch switches to the specified branch.
// Note: We don't use "--" separator here because it would treat the branch as a file path.
// Branch existence is validated by the caller using go-git before calling this function.
func CheckoutBranch(branch string) error {
	cmd := exec.Command("git", "checkout", branch)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to checkout %s: %s - %w", branch, string(output), err)
	}
	return nil
}

// CreateBranch creates a new branch from the current HEAD.
func CreateBranch(name string) error {
	cmd := exec.Command("git", "branch", "--", name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create branch %s: %s - %w", name, string(output), err)
	}
	return nil
}

// CreateBranchFrom creates a new branch from a specific ref.
func CreateBranchFrom(name, ref string) error {
	cmd := exec.Command("git", "branch", "--", name, ref)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create branch %s from %s: %s - %w", name, ref, string(output), err)
	}
	return nil
}

// DeleteBranch deletes a branch.
func DeleteBranch(name string) error {
	cmd := exec.Command("git", "branch", "-D", "--", name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete branch %s: %s - %w", name, string(output), err)
	}
	return nil
}

// AmendCommitMessage amends the last commit message.
func AmendCommitMessage(message string) error {
	cmd := exec.Command("git", "commit", "--amend", "-m", message)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to amend commit: %s - %w", string(output), err)
	}
	return nil
}

// GetCurrentCommitSHA returns the SHA of the current HEAD.
func GetCurrentCommitSHA() (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current commit SHA: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// Fetch fetches from the specified remote.
func Fetch(remote string) error {
	cmd := exec.Command("git", "fetch", remote)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to fetch from %s: %s - %w", remote, string(output), err)
	}
	return nil
}
