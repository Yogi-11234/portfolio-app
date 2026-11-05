package sanitize

import (
	"strings"
	"testing"
)

func TestMarkdownSanitizer(t *testing.T) {
	// 1. Standard markdown
	input := "# Hello World\n\nThis is **bold** text."
	html, err := MarkdownToHTML(input)
	if err != nil {
		t.Fatalf("MarkdownToHTML failed: %v", err)
	}
	if !strings.Contains(html, "<h1>Hello World</h1>") && !strings.Contains(html, "<h1") {
		t.Errorf("Expected h1 heading, got: %s", html)
	}
	if !strings.Contains(html, "<strong>bold</strong>") {
		t.Errorf("Expected strong tag, got: %s", html)
	}

	// 2. XSS injection attempt in markdown
	xssInput := "Check this: <script>alert('xss')</script> and <img src=x onerror=alert(1)>"
	cleanHtml, err := MarkdownToHTML(xssInput)
	if err != nil {
		t.Fatalf("MarkdownToHTML failed on XSS: %v", err)
	}
	if strings.Contains(cleanHtml, "<script") || strings.Contains(cleanHtml, "onerror=") {
		t.Fatalf("XSS payload leaked! Result: %s", cleanHtml)
	}
}

func TestValidateHTTPSURL(t *testing.T) {
	if !ValidateHTTPSURL("https://github.com/myproject") {
		t.Error("Expected valid https URL")
	}
	if ValidateHTTPSURL("http://insecure.com") {
		t.Error("Expected non-https URL to be rejected")
	}
	if ValidateHTTPSURL("javascript:alert(1)") {
		t.Error("Expected javascript URL to be rejected")
	}
	if ValidateHTTPSURL("data:text/html;base64,PHNjcmlwdD4=") {
		t.Error("Expected data URL to be rejected")
	}
}
