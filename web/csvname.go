package web

import (
	"regexp"
	"strings"
)

var (
	// nonAlphaNum matches anything that is not a letter, digit, hyphen, or underscore.
	nonAlphaNum = regexp.MustCompile(`[^a-z0-9\-_]+`)
	// multiHyphen collapses consecutive hyphens into one.
	multiHyphen = regexp.MustCompile(`-{2,}`)
)

const maxCSVNameLen = 100

// SanitizeJobName converts a job name into a safe, lowercase, hyphen-separated
// basename suitable for use as a filesystem filename (without extension).
//
// Examples:
//
//	"My Search Query"   → "my-search-query"
//	"Coffee Shops, NYC" → "coffee-shops-nyc"
//	"  leading/trailing " → "leading-trailing"
func SanitizeJobName(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	// Replace whitespace and slashes with hyphens before stripping.
	s = strings.NewReplacer(" ", "-", "/", "-", "\\", "-").Replace(s)
	// Remove everything that is not alphanumeric, hyphen, or underscore.
	s = nonAlphaNum.ReplaceAllString(s, "")
	// Collapse consecutive hyphens.
	s = multiHyphen.ReplaceAllString(s, "-")
	// Trim leading/trailing hyphens.
	s = strings.Trim(s, "-")

	if s == "" {
		s = "untitled"
	}

	if len(s) > maxCSVNameLen {
		s = s[:maxCSVNameLen]
		s = strings.TrimRight(s, "-")
	}

	return s
}

// CSVFileName returns the full CSV filename (with .csv extension) for a job name.
func CSVFileName(jobName string) string {
	return SanitizeJobName(jobName) + ".csv"
}
