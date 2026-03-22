package fullterm

import (
	"testing"
)

func TestVisibleContentShorterThanWindow(t *testing.T) {
	content := []string{"line1\n", "line2\n", "line3\n"}
	result := visibleContent(content, 10, 0)
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
	result := visibleContent(content, len(content)-1, 0)
	if len(result) != len(content) {
		t.Fatalf("expected %d rows, got %d", len(content), len(result))
	}
}

func TestVisibleContentOverflowsWindow(t *testing.T) {
	content := []string{"line1\n", "line2\n", "line3\n", "line4\n", "line5\n"}
	height := 2
	result := visibleContent(content, height, 0)
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
	result := visibleContent([]string{}, 10, 0)
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d rows", len(result))
	}
}

func TestVisibleContentZeroHeight(t *testing.T) {
	content := []string{"line1\n", "line2\n", "line3\n"}
	result := visibleContent(content, 0, 0)
	// startRow = max(3 - 1, 0) = 2, so only last row
	if len(result) != 1 {
		t.Fatalf("expected 1 row, got %d", len(result))
	}
	if result[0] != "line3\n" {
		t.Errorf("expected 'line3\\n', got %q", result[0])
	}
}

func TestVisibleContentScrolled(t *testing.T) {
	content := []string{"line1\n", "line2\n", "line3\n", "line4\n", "line5\n"}
	height := 2
	// scrollOffset=2: endRow = 5-2 = 3, startRow = max(3-3, 0) = 0, so rows 0,1,2
	result := visibleContent(content, height, 2)
	expectedLen := height + 1
	if len(result) != expectedLen {
		t.Fatalf("expected %d rows, got %d", expectedLen, len(result))
	}
	if result[0] != "line1\n" {
		t.Errorf("expected first visible row to be 'line1\\n', got %q", result[0])
	}
	if result[len(result)-1] != "line3\n" {
		t.Errorf("expected last visible row to be 'line3\\n', got %q", result[len(result)-1])
	}
}

func TestVisibleContentScrollOffsetPastTop(t *testing.T) {
	content := []string{"line1\n", "line2\n", "line3\n"}
	// scrollOffset larger than content — should return empty, not panic
	result := visibleContent(content, 2, 100)
	if len(result) != 0 {
		t.Fatalf("expected empty result when scrolled past top, got %d rows", len(result))
	}
}

func TestVisibleContentScrolledPartial(t *testing.T) {
	content := []string{"line1\n", "line2\n", "line3\n", "line4\n", "line5\n"}
	height := 10
	// scrollOffset=1: endRow = 4, startRow = max(4-11, 0) = 0, all 4 rows visible
	result := visibleContent(content, height, 1)
	if len(result) != 4 {
		t.Fatalf("expected 4 rows, got %d", len(result))
	}
	if result[len(result)-1] != "line4\n" {
		t.Errorf("expected last row to be 'line4\\n', got %q", result[len(result)-1])
	}
}
