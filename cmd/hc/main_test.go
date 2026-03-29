package main

import (
	"testing"
)

func TestMatchByPatternExact(t *testing.T) {
	results, err := matchByPattern(statusCodes, "404")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Code != 404 {
		t.Errorf("expected code 404, got %d", results[0].Code)
	}
	if results[0].Message != "Not Found" {
		t.Errorf("expected 'Not Found', got %q", results[0].Message)
	}
}

func TestMatchByPatternWildcard(t *testing.T) {
	results, err := matchByPattern(statusCodes, "1xx")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 4 {
		t.Fatalf("expected 4 results for 1xx, got %d", len(results))
	}
	for _, sc := range results {
		if sc.Code < 100 || sc.Code >= 200 {
			t.Errorf("unexpected code %d in 1xx results", sc.Code)
		}
	}
}

func TestMatchByPatternRegex(t *testing.T) {
	results, err := matchByPattern(statusCodes, "30[12]")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	codes := map[int]bool{results[0].Code: true, results[1].Code: true}
	if !codes[301] || !codes[302] {
		t.Errorf("expected 301 and 302, got %v", codes)
	}
}

func TestMatchByPatternNoMatch(t *testing.T) {
	results, err := matchByPattern(statusCodes, "999")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestMatchByPatternInvalidRegex(t *testing.T) {
	_, err := matchByPattern(statusCodes, "[invalid")
	if err == nil {
		t.Error("expected error for invalid regex")
	}
}

func TestMatchByTextCaseInsensitive(t *testing.T) {
	results, err := matchByText(statusCodes, "teapot")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Code != 418 {
		t.Errorf("expected code 418, got %d", results[0].Code)
	}
}

func TestMatchByTextMultipleResults(t *testing.T) {
	results, err := matchByText(statusCodes, "timeout")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) < 3 {
		t.Errorf("expected at least 3 results for 'timeout', got %d", len(results))
	}
}

func TestMatchByTextNoMatch(t *testing.T) {
	results, err := matchByText(statusCodes, "zzzzzznothing")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestMatchByTextInvalidRegex(t *testing.T) {
	_, err := matchByText(statusCodes, "[invalid")
	if err == nil {
		t.Error("expected error for invalid regex")
	}
}

func TestMatchByTextSearchesMessageAndExplain(t *testing.T) {
	// "WebDAV" appears in Explain, not Message
	results, err := matchByText(statusCodes, "WebDAV")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 {
		t.Error("expected results for 'WebDAV' (appears in explanations)")
	}
}

func TestStatusCodesNotEmpty(t *testing.T) {
	if len(statusCodes) == 0 {
		t.Fatal("statusCodes should not be empty")
	}
}

func TestNoDuplicateCodes(t *testing.T) {
	seen := make(map[int]bool)
	for _, sc := range statusCodes {
		if seen[sc.Code] {
			t.Errorf("duplicate code: %d", sc.Code)
		}
		seen[sc.Code] = true
	}
}
