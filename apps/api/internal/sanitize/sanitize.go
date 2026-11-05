package sanitize

import (
	"bytes"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"golang.org/x/text/unicode/norm"
)

var (
	mdParser = goldmark.New()
	ugcPolicy = bluemonday.UGCPolicy()
)

// NormalizeText trims spaces and normalizes unicode to NFC form.
func NormalizeText(input string) string {
	trimmed := strings.TrimSpace(input)
	return norm.NFC.String(trimmed)
}

// MarkdownToHTML parses markdown safely (raw HTML disabled) and sanitizes output with bluemonday UGCPolicy.
func MarkdownToHTML(markdownInput string) (string, error) {
	cleanInput := NormalizeText(markdownInput)
	if cleanInput == "" {
		return "", nil
	}

	var buf bytes.Buffer
	// Goldmark by default does NOT parse raw HTML (html.WithUnsafe() is omitted)
	if err := mdParser.Convert([]byte(cleanInput), &buf); err != nil {
		return "", err
	}

	// Sanitize output through bluemonday UGC Policy to eliminate any residual XSS vectors
	sanitizedBytes := ugcPolicy.SanitizeBytes(buf.Bytes())
	return string(sanitizedBytes), nil
}

// ValidateHTTPSURL checks that a URL uses https:// and does not contain javascript: or data:
func ValidateHTTPSURL(url string) bool {
	trimmed := strings.TrimSpace(url)
	if trimmed == "" {
		return true // empty URL is allowed if optional
	}
	lower := strings.ToLower(trimmed)
	if strings.HasPrefix(lower, "javascript:") || strings.HasPrefix(lower, "data:") {
		return false
	}
	return strings.HasPrefix(lower, "https://")
}
