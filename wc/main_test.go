package main

import (
	"os"
	"strings"
	"testing"
)

func TestCountBytes(t *testing.T) {
	tests := []struct {
		name     string
		content  []byte
		expected int
	}{
		{"empty", []byte(""), 0},
		{"single char", []byte("a"), 1},
		{"hello world", []byte("hello world"), 11},
		{"with newlines", []byte("hello\nworld"), 11},
		{"unicode", []byte("café"), 5}, // é is 2 bytes in UTF-8
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := countBytes(tt.content)
			if result != tt.expected {
				t.Errorf("countBytes() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCountLines(t *testing.T) {
	tests := []struct {
		name     string
		content  []byte
		expected int
	}{
		{"empty", []byte(""), 0},
		{"single line", []byte("hello world"), 0},
		{"one line with newline", []byte("hello world\n"), 1},
		{"two lines", []byte("hello\nworld"), 1},
		{"three lines", []byte("hello\nworld\ntest"), 2},
		{"empty lines", []byte("\n\n\n"), 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := countLines(tt.content)
			if result != tt.expected {
				t.Errorf("countLines() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCountWords(t *testing.T) {
	tests := []struct {
		name     string
		content  []byte
		expected int
	}{
		{"empty", []byte(""), 0},
		{"single word", []byte("hello"), 1},
		{"two words", []byte("hello world"), 2},
		{"multiple words", []byte("hello world test"), 3},
		{"with newlines", []byte("hello\nworld"), 2},
		{"multiple spaces", []byte("hello   world"), 2},
		{"tabs and spaces", []byte("hello\tworld test"), 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := countWords(tt.content)
			if result != tt.expected {
				t.Errorf("countWords() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCountCharacters(t *testing.T) {
	tests := []struct {
		name     string
		content  []byte
		expected int
	}{
		{"empty", []byte(""), 0},
		{"single char", []byte("a"), 1},
		{"hello world", []byte("hello world"), 11},
		{"with newlines", []byte("hello\nworld"), 11},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := countCharacters(tt.content)
			if result != tt.expected {
				t.Errorf("countCharacters() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIntegrationWithFile(t *testing.T) {
	// Create a temporary test file
	testContent := "hello world\nthis is a test\nwith multiple lines"
	tmpfile, err := os.CreateTemp("", "wc_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(testContent)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	// Test each flag
	content, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Test bytes
	bytes := countBytes(content)
	if bytes != len(testContent) {
		t.Errorf("Expected %d bytes, got %d", len(testContent), bytes)
	}

	// Test lines
	lines := countLines(content)
	expectedLines := strings.Count(testContent, "\n")
	if lines != expectedLines {
		t.Errorf("Expected %d lines, got %d", expectedLines, lines)
	}

	// Test words
	words := countWords(content)
	expectedWords := len(strings.Fields(testContent))
	if words != expectedWords {
		t.Errorf("Expected %d words, got %d", expectedWords, words)
	}

	// Test characters
	chars := countCharacters(content)
	expectedChars := len([]rune(testContent))
	if chars != expectedChars {
		t.Errorf("Expected %d characters, got %d", expectedChars, chars)
	}
}

func TestEdgeCases(t *testing.T) {
	// Test with very large content
	largeContent := strings.Repeat("hello world\n", 1000)
	content := []byte(largeContent)

	lines := countLines(content)
	if lines != 1000 {
		t.Errorf("Expected 1000 lines, got %d", lines)
	}

	words := countWords(content)
	if words != 2000 { // 2 words per line * 1000 lines
		t.Errorf("Expected 2000 words, got %d", words)
	}

	// Test with only whitespace
	whitespaceContent := []byte("   \t\n   ")
	words = countWords(whitespaceContent)
	if words != 0 {
		t.Errorf("Expected 0 words for whitespace-only content, got %d", words)
	}
}

func TestLorumFile(t *testing.T) {
	// Check if lorum.txt exists
	if _, err := os.Stat("lorum.txt"); os.IsNotExist(err) {
		t.Skip("lorum.txt file not found, skipping test")
	}

	content, err := os.ReadFile("lorum.txt")
	if err != nil {
		t.Fatalf("Failed to read lorum.txt: %v", err)
	}

	// Test based on the expected values from README
	bytes := countBytes(content)
	if bytes != 445 {
		t.Errorf("Expected 445 bytes, got %d", bytes)
	}

	lines := countLines(content)
	if lines != 4 {
		t.Errorf("Expected 4 lines, got %d", lines)
	}

	words := countWords(content)
	if words != 69 {
		t.Errorf("Expected 69 words, got %d", words)
	}

	chars := countCharacters(content)
	if chars != 445 {
		t.Errorf("Expected 445 characters, got %d", chars)
	}
}
