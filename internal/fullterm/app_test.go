package fullterm

import (
	"testing"
)

func TestVisibleContentShorterThanWindow(t *testing.T) {
	content := []string{"line1\n", "line2\n", "line3\n"}
	result := visibleContent(content, 10)
	if len(result) != len(content) {
		t.Fatalf("expected %d rows, got %d", len(content), len(result))
	}
	for i, row := range result {
		if row != content[i] {
			t.Errorf("row %d: expected %q, got %q", i, content[i], row)
		}
	}
}

func TestVisibleContentExactlyFitsWindow(t *testing.T) {
	content := []string{"line1\n", "line2\n", "line3\n"}
	// height+1 == len(content), so startRow == 0
	result := visibleContent(content, len(content)-1)
	if len(result) != len(content) {
		t.Fatalf("expected %d rows, got %d", len(content), len(result))
	}
}

func TestVisibleContentOverflowsWindow(t *testing.T) {
	content := []string{"line1\n", "line2\n", "line3\n", "line4\n", "line5\n"}
	height := 2
	result := visibleContent(content, height)
	// startRow = max(5 - 3, 0) = 2, so rows 2,3,4
	expectedLen := height + 1
	if len(result) != expectedLen {
		t.Fatalf("expected %d rows, got %d", expectedLen, len(result))
	}
	if result[0] != "line3\n" {
		t.Errorf("expected first visible row to be 'line3\\n', got %q", result[0])
	}
	if result[len(result)-1] != "line5\n" {
		t.Errorf("expected last visible row to be 'line5\\n', got %q", result[len(result)-1])
	}
}

func TestVisibleContentEmpty(t *testing.T) {
	result := visibleContent([]string{}, 10)
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d rows", len(result))
	}
}

func TestVisibleContentZeroHeight(t *testing.T) {
	content := []string{"line1\n", "line2\n", "line3\n"}
	result := visibleContent(content, 0)
	// startRow = max(3 - 1, 0) = 2, so only last row
	if len(result) != 1 {
		t.Fatalf("expected 1 row, got %d", len(result))
	}
	if result[0] != "line3\n" {
		t.Errorf("expected 'line3\\n', got %q", result[0])
	}
}
