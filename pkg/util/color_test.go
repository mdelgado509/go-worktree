package util

import (
	"testing"
)

// TestColorize tests the Colorize function
func TestColorize(t *testing.T) {
	testCases := []struct {
		text     string
		color    string
		expected string
	}{
		{"test", ColorRed, "\033[31mtest\033[0m"},
		{"hello", ColorBlue, "\033[34mhello\033[0m"},
		{"", ColorGreen, "\033[32m\033[0m"},
	}

	for _, tc := range testCases {
		result := Colorize(tc.text, tc.color)
		if result != tc.expected {
			t.Errorf("Expected %q, got %q", tc.expected, result)
		}
	}
}

// TestBold tests the Bold function
func TestBold(t *testing.T) {
	testCases := []struct {
		text     string
		expected string
	}{
		{"test", "\033[1mtest\033[0m"},
		{"hello", "\033[1mhello\033[0m"},
		{"", "\033[1m\033[0m"},
	}

	for _, tc := range testCases {
		result := Bold(tc.text)
		if result != tc.expected {
			t.Errorf("Expected %q, got %q", tc.expected, result)
		}
	}
}
