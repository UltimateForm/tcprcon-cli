package ansi

import "testing"

func TestFormatSingleFlag(t *testing.T) {
	result := Format("hello", Blue)
	expected := "\033[34mhello\033[0m"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestFormatMultipleFlags(t *testing.T) {
	result := Format("hello", Blue, Bold)
	expected := "\033[34;1mhello\033[0m"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}
