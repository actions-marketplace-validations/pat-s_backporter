// Package version provides version information for the application.
package version

import "fmt"

// Version is set at build time via ldflags.
var Version = "dev"

// BuildDate is set at build time via ldflags.
var BuildDate = "unknown"

// GitURL is the repository URL.
const GitURL = "https://codefloe.com/pat-s/backporter"

// String returns the version string with build date.
func String() string {
	return fmt.Sprintf("%s (%s)", Version, BuildDate)
}

// Full returns the full version string with git URL.
func Full() string {
	return fmt.Sprintf("backporter %s (%s)", Version, GitURL)
}

// SignatureMessage returns the message to append to backport commits.
func SignatureMessage(originalSHA string) string {
	return fmt.Sprintf("Backported from %s using backporter %s (%s)", originalSHA, Version, GitURL)
}
